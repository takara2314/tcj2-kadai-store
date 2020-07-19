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

// ttUpdateSchedule はTimeTreeに指定されたスケジュールを更新する関数
func ttUpdateSchedule(hwID string) {
	var POSTjson []byte
	// POSTに必要なJSONの情報を格納する構造体を初期化
	var POSTjsonData SetSchedule
	var scheData SetScheduleData
	var scheAttribute SetScheduleAttributes
	var scheRelate SetScheduleRelationships
	var scheLabel SetScheduleLabel
	var scheLabelData SetScheduleLabelData

	scheAttribute = SetScheduleAttributes{
		Category:      "schedule",
		Title:         hwStatus[hwID][5].(string),
		AllDay:        false,
		StartAt:       hwStatus[hwID][4].(time.Time),
		StartTimezone: "Asia/Tokyo",
		EndAt:         hwStatus[hwID][4].(time.Time),
		EndTimezone:   "Asia/Tokyo",
		Description:   hwStatus[hwID][0].(string) + "の課題です。",
	}

	scheLabelData = SetScheduleLabelData{
		ID:   "KADAI_LABEL_ID",
		Type: "label",
	}

	scheLabel = SetScheduleLabel{
		Data: scheLabelData,
	}
	scheRelate = SetScheduleRelationships{
		Label: scheLabel,
	}
	scheData = SetScheduleData{
		Attributes:    scheAttribute,
		Relationships: scheRelate,
	}
	POSTjsonData = SetSchedule{
		Data: scheData,
	}

	// J2AのカレンダーにPOST
	POSTjsonData.Data.Relationships.Label.Data.ID = os.Getenv("J2A_KADAI_LABEL_ID")
	// 構造体をJSONに変換
	POSTjson, _ = json.Marshal(POSTjsonData)
	ttUpdateSchedulePOST(os.Getenv("J2A_CALENDAR_ID"), hwStatus[hwID][6].(string), POSTjson)

	// J2BのカレンダーにPOST
	POSTjsonData.Data.Relationships.Label.Data.ID = os.Getenv("J2B_KADAI_LABEL_ID")
	// 構造体をJSONに変換
	POSTjson, _ = json.Marshal(POSTjsonData)
	ttUpdateSchedulePOST(os.Getenv("J2B_CALENDAR_ID"), hwStatus[hwID][7].(string), POSTjson)
}

// ttUpdateSchedulePOST はTimeTreeに予定変更リクエスト(PUT)を送る関数
func ttUpdateSchedulePOST(calendarID string, eventID string, POSTjson []byte) {
	// TimeTree API
	var baseURL string = "https://timetreeapis.com/"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	// calendars にアクセス
	reqURL.Path = path.Join(reqURL.Path, "calendars")
	reqURL.Path = path.Join(reqURL.Path, calendarID)
	reqURL.Path = path.Join(reqURL.Path, "events")
	reqURL.Path = path.Join(reqURL.Path, eventID)

	// リクエスト詳細を定義
	req, _ := http.NewRequest("PUT", reqURL.String(), bytes.NewBuffer(POSTjson))
	// アプリ情報やトークン情報を明記
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.timetree.v1+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TIMETREE_API_TOKEN"))

	// PUTリクエストを送信
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("エラーが発生しました: %v\n", err)
	}
}
