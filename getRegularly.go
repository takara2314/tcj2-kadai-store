package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// getRegularly は定期的に TCJ2 Kadai Store API から課題一覧を取得
func getRegularly(getTime []int) {
	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら実行
		if containsInt(getTime, nowMinute) {
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
			reqURLvar.Add("timezone", "JST")
			reqURL.RawQuery = reqURLvar.Encode()

			// リクエスト詳細を定義
			req, _ := http.NewRequest("GET", reqURL.String(), nil)
			// トークン情報を明記
			req.Header.Add("Authorization", "Bearer "+os.Getenv("KADAI_API_TOKEN"))

			// APIを叩いてレスポンスを受け取る
			response, _ := http.DefaultClient.Do(req)

			body, _ := ioutil.ReadAll(response.Body)
			// fmt.Println(string(body))

			// レスポンスされたJSONを構造体化
			var hwStatusStruct GetHomeworks
			err = json.Unmarshal(body, &hwStatusStruct)
			if err != nil {
				panic(err)
			}

			// 構造体化したものをmapに変換
			for _, hwInfo := range hwStatusStruct.Homeworks {
				hwSubject := hwInfo.Subject
				hwOmitted := hwInfo.Omitted
				hwName := hwInfo.Name
				hwID := hwInfo.ID
				hwDue := hwInfo.Due

				// 課題情報を抽出
				hwStatus[hwID] = []interface{}{hwSubject, hwOmitted, hwName, hwDue}
				// 課題リストを抽出
				hwList = append(hwList, hwInfo.ID)
			}

			// 前の課題情報と比べて変更点がないかをチェック
			newHW, updateHW, deleteHW := checkChanges()
			fmt.Println("新規追加ID:", newHW)
			fmt.Println("内容変更ID:", updateHW)
			fmt.Println("削除ID:", deleteHW)
			time.Sleep(1 * time.Minute)
		}
	}
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
