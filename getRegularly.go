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
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// getRegularly は定期的に課題一覧を取得したり、TimeTreeのスケジュールを変更する関数
func getRegularly(getTime []int) {
	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら実行
		if _, isExist := containsInt(getTime, nowMinute); isExist {

			// 新鮮ピッチピチな課題情報を入れるために、要素数を0にする
			hwStatus = make(map[string][]interface{}, 0)
			hwList = make([]string, 0)

			// TCJ2 Kadai Store API
			var baseURL string = "https://kadai-store-api.appspot.com/"
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

			// ステータスコードが200でなければ、処理はパス
			if res.StatusCode != 200 {
				time.Sleep(1 * time.Minute)
				continue
			}

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
				hwCourse := hwInfo.Course
				hwSubject := hwInfo.Subject
				hwSubjectID := hwInfo.SubjectID
				hwName := hwInfo.Name
				hwID := hwInfo.ID
				hwDue := hwInfo.Due
				// TimeTreeに関連する課題情報
				// スケジュール名は省略形教科名と課題名を元に決める
				hwTTscheName := scheNameGen(hwSubjectID, hwName)
				hwTTIDA := ""
				hwTTIDB := ""

				// 課題情報を抽出
				// (教科名、省略された教科名、課題名、課題ID、提出期限、TimeTreeスケジュール名、TimeTreeカレンダーID(A)、 TimeTreeカレンダーID(B))
				hwStatus[hwID] = []interface{}{
					hwCourse,
					hwSubject,
					hwSubjectID,
					hwName,
					hwID,
					hwDue,
					hwTTscheName,
					hwTTIDA,
					hwTTIDB,
				}
				// 課題リストを抽出
				hwList = append(hwList, hwID)
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
				hwStatus[hwID][7], hwStatus[hwID][8] = ttAddSchedule(
					hwStatus[hwID][3].(string),
					hwStatus[hwID][1].(string),
					hwStatus[hwID][5].(time.Time),
				)

				// 過去の課題情報に今の課題情報を書き加える
				hwStatusPast[hwID] = hwStatus[hwID]
				hwListPast = append(hwListPast, hwID)

				// 課題情報をFirebaseに保存
				dbSetKadai(ctx, client, hwID, hwStatus[hwID])
			}

			// 内容変更があったものをスケジュールに反映・データベースを更新
			for _, hwID := range updateHW {
				fmt.Println("課題の内容が変更されました:", hwID)
				// TimeTree関連のデータは新規作成時にしか取得できないので、過去のものを引き継ぐ
				hwStatus[hwID][6] = hwStatusPast[hwID][6]
				hwStatus[hwID][7] = hwStatusPast[hwID][7]
				hwStatus[hwID][8] = hwStatusPast[hwID][8]
				// そのIDを渡してカレンダー情報を変更してもらう
				ttUpdateSchedule(hwID)

				// 過去の課題情報に今の課題情報を書きかえる
				hwStatusPast[hwID] = hwStatus[hwID]

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
func scheNameGen(hwSubjectID string, hwName string) (gened string) {
	// 教科ID表の中に、そのIDがあるかどうか
	index, exist := containID(subjects.Subjects, hwSubjectID)
	if exist {
		return fmt.Sprintf("%s %s", subjects.Subjects[index], hwName)
	}

	// 含まれていない場合は、課題名をそのまま返す
	return hwName
}

// containsInt はint型のスライスから特定の整数の要素数番号と存在するかを返す関数
func containsInt(tSlice []int, tNum int) (int, bool) {
	for i, num := range tSlice {
		if tNum == num {
			return i, true
		}
	}
	return -1, false
}

// containID は、登録されている教科ID表の中にそのIDがあるかどうかを返す関数
func containID(tSlice [][]string, id string) (int, bool) {
	for i, subject := range tSlice {
		if id == subject[0] {
			return i, true
		}
	}
	return -1, false
}
