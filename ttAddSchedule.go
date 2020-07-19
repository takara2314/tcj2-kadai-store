package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// ttAddSchedule はTimeTreeに指定されたスケジュールを追加する関数
func ttAddSchedule(hwName string, hwDue time.Time) (hwTTIDA string, hwTTIDB string) {
	var err error
	var POSTjson []byte
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
		StartAt:       hwDue,
		StartTimezone: "Asia/Tokyo",
		EndAt:         hwDue,
		EndTimezone:   "Asia/Tokyo",
	}

	scheLabelData = AddScheduleLabelData{
		ID:   "KADAI_LABEL_ID",
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

	// J2AのカレンダーにPOST
	POSTjsonData.Data.Relationships.Label.Data.ID = os.Getenv("J2A_KADAI_LABEL_ID")
	// 構造体をJSONに変換
	POSTjson, _ = json.Marshal(POSTjsonData)
	hwTTIDA, err = ttAddSchedulePOST(os.Getenv("J2A_CALENDAR_ID"), POSTjson)
	if err != nil {
		panic(err)
	}

	// J2BのカレンダーにPOST
	POSTjsonData.Data.Relationships.Label.Data.ID = os.Getenv("J2B_KADAI_LABEL_ID")
	// 構造体をJSONに変換
	POSTjson, _ = json.Marshal(POSTjsonData)
	hwTTIDB, err = ttAddSchedulePOST(os.Getenv("J2B_CALENDAR_ID"), POSTjson)
	if err != nil {
		panic(err)
	}

	return
}

// ttAddSchedulePOST はTimeTreeに予定追加リクエスト(POST)を送る関数
func ttAddSchedulePOST(calendarID string, POSTjson []byte) (string, error) {
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

	// リクエスト詳細を定義
	req, _ := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(POSTjson))
	// アプリ情報やトークン情報を明記
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.timetree.v1+json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("TIMETREE_API_TOKEN"))

	// POSTリクエストを送信
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("エラーが発生しました: %s\n", err)
		// 予定IDは発行できなかったので、空白にしてエラーとして返す
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	// レスポンスされたものを格納する構造体
	var resData ResAddSchedule
	// JSONをmapに変換
	json.Unmarshal(body, &resData)
	// 予定IDを返す
	return resData.Data.ID, nil
}
