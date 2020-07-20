package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// getRegularly は定期的に課題一覧を取得したり、TimeTreeのスケジュールを変更する関数
func getRegularly(getTime []int) {
	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら実行
		if containsInt(getTime, nowMinute) {
			// 新鮮ピッチピチな課題情報を入れるために、要素数を0にする
			hwStatus = make(map[string][]interface{}, 0)
			hwList = make([]string, 0)

			// TCJ2 Kadai Store API
			var baseURL string = "http://tcj2-kadai-store-api.2314.tk/"
			reqURL, err := url.Parse(baseURL)
			if err != nil {
				panic(err)
			}

			// get にアクセス
			reqURL.Path = path.Join(reqURL.Path, "get")
			// dueパラメータを指定
			reqURLvar, _ := url.ParseQuery(reqURL.RawQuery)
			reqURLvar.Add("due", "future")
			reqURLvar.Add("timezone", "Asia/Tokyo")
			reqURL.RawQuery = reqURLvar.Encode()

			// リクエスト詳細を定義
			req, _ := http.NewRequest("GET", reqURL.String(), nil)
			// トークン情報を明記
			req.Header.Add("Authorization", "Bearer "+os.Getenv("KADAI_API_TOKEN"))

			// APIを叩いてレスポンスを受け取る
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalf("エラーが発生しました: %v\n", err)
			}
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)

			// レスポンスされたJSONを構造体化
			var hwStatusStruct GetHomeworks
			err = json.Unmarshal(body, &hwStatusStruct)
			if err != nil {
				panic(err)
			}

			// 構造体化したものをmapに変換
			for _, hwInfo := range hwStatusStruct.Homeworks {
				// 基本的な課題情報
				hwSubject := hwInfo.Subject
				hwOmitted := hwInfo.Omitted
				hwName := hwInfo.Name
				hwID := hwInfo.ID
				hwDue := hwInfo.Due
				// TimeTreeに関連する課題情報
				// スケジュール名は省略形教科名と課題名を元に決める
				hwTTscheName := scheNameGen(hwOmitted, hwName)
				hwTTIDA := ""
				hwTTIDB := ""

				// 課題情報を抽出
				// (教科名、省略された教科名、課題名、課題ID、提出期限、TimeTreeスケジュール名、TimeTreeカレンダーID(A)、 TimeTreeカレンダーID(B))
				hwStatus[hwID] = []interface{}{hwSubject, hwOmitted, hwName, hwID, hwDue, hwTTscheName, hwTTIDA, hwTTIDB}
				// 課題リストを抽出
				hwList = append(hwList, hwInfo.ID)
			}

			// Firebaseを初期化
			ctx := context.Background()
			sa := option.WithCredentialsFile("tcj2-kadai-store-ed48273c015c.json")
			app, err := firebase.NewApp(ctx, nil, sa)
			if err != nil {
				log.Fatalln(err)
			}

			client, err := app.Firestore(ctx)
			if err != nil {
				log.Fatalln(err)
			}
			defer client.Close()

			// 前の課題情報と比べて変更点がないかをチェック
			newHW, updateHW, _ := checkChanges()
			fmt.Println(time.Now().Format("2006年1月2日 15時4分5秒"), "現在")
			fmt.Println("新規追加ID:", newHW)
			fmt.Println("内容変更ID:", updateHW)
			// fmt.Println("削除ID:", deleteHW)

			// 新規追加されたものをスケジュールとデータベースに追加
			for _, hwID := range newHW {
				// スケジュール名と課題名と提出期限を渡してTimeTreeに予定として追加してもらい、
				// 予定IDを取得
				hwStatus[hwID][6], hwStatus[hwID][7] = ttAddSchedule(
					hwStatus[hwID][5].(string),
					hwStatus[hwID][0].(string),
					hwStatus[hwID][4].(time.Time),
				)

				// 課題情報をFirebaseに保存
				dbSetKadai(ctx, client, hwID, hwStatus[hwID])
			}

			// 内容変更があったものをスケジュールに反映・データベースを更新
			for _, hwID := range updateHW {
				fmt.Println("変更します！！！")
				// TimeTree関連のデータは新規作成時にしか取得できないので、過去のものを引き継ぐ
				hwStatus[hwID][5] = hwStatusPast[hwID][5]
				hwStatus[hwID][6] = hwStatusPast[hwID][6]
				hwStatus[hwID][7] = hwStatusPast[hwID][7]
				// そのIDを渡してカレンダー情報を変更してもらう
				ttUpdateSchedule(hwID)
				// 課題情報をFirebaseに上書き保存
				dbSetKadai(ctx, client, hwID, hwStatus[hwID])
			}

			// 削除された課題をスケジュールとデータベースからも削除
			// → 今後のアップデートで実装予定

			fmt.Println("task finished")
			time.Sleep(1 * time.Minute)
		}
	}
}

// scheNameGen はTimeTreeで表示するスケジュール名を生成する関数
func scheNameGen(hwOmitted string, hwName string) (gened string) {
	// 課題名の中に省略された教科名が含まれていたら、課題名をそのまま返す
	if strings.Contains(hwName, hwOmitted) {
		return hwName
	}
	// 含まれていない場合は、"省略教科名 課題名"という形にする
	return hwOmitted + " " + hwName
}

// containsInt はint型のスライスから特定の整数があればtrueを返す関数
func containsInt(tSlice []int, tNum int) bool {
	for _, num := range tSlice {
		if tNum == num {
			return true
		}
	}
	return false
}
