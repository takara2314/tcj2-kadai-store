package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// ttAddSchedule はTimeTreeに指定されたスケジュールを追加する関数
func ttAddSchedule(hwName string, dueTime time.Time) {
	// POSTに必要なJSONの情報を格納する構造体を初期化
	var POSTjsonData AddSchedule
	var scheData AddScheduleData
	var scheAttribute AddScheduleAttributes
	var scheRelate AddScheduleRelationships
	var scheLabel AddScheduleLabel
	var scheLabelData AddScheduleLabelData

	scheAttribute = AddScheduleAttributes{
		Category:      "schedule",
		Title:         hwName,
		AllDay:        false,
		StartAt:       dueTime,
		StartTimezone: "Asia/Tokyo",
		EndAt:         dueTime,
		EndTimezone:   "Asia/Tokyo",
	}

	scheLabelData = AddScheduleLabelData{
		ID:   os.Getenv("KADAI_LABEL_ID"),
		Type: "label",
	}

	scheLabel = AddScheduleLabel{
		Data: scheLabelData,
	}
	scheRelate = AddScheduleRelationships{
		Label: scheLabel,
	}
	scheData = AddScheduleData{
		Attributes:    scheAttribute,
		Relationships: scheRelate,
	}
	POSTjsonData = AddSchedule{
		Data: scheData,
	}

	// 構造体をJSONに変換
	POSTjson, _ := json.Marshal(POSTjsonData)

	// TimeTree API
	var baseURL string = "https://timetreeapis.com/"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	// calendars にアクセス
	reqURL.Path = path.Join(reqURL.Path, "calendars")
	reqURL.Path = path.Join(reqURL.Path, os.Getenv("J2A_CALENDAR_ID"))
	reqURL.Path = path.Join(reqURL.Path, "events")

	// リクエスト詳細を定義
	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(POSTjson))
	// アプリ情報やトークン情報を明記
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.timetree.v1+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TIMETREE_API_TOKEN"))

	// POSTリクエストを送信
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("エラーが発生しました: %s\n", err)
	}
}

// // ttUpdateSchedule はTimeTreeに指定されたスケジュールを更新する関数
// func ttUpdateSchedule() {
// 	hoge
// }

// // ttDeleteSchedule はTimeTreeに指定されたスケジュールを削除する関数
// func ttDeleteSchedule() {
// 	hoge
// }
