package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
)

const location = "Asia/Tokyo"

var (
	// 現在の課題情報
	hwStatus map[string][]interface{} = make(map[string][]interface{}, 0)
	// 現在の課題リスト (ID)
	hwList []string
	// 前の課題情報
	hwStatusPast map[string][]interface{} = make(map[string][]interface{}, 0)
	// 前の課題リスト (ID)
	hwListPast []string
)

// GetHomeworks はAPIから取得したJSONを収納する構造体
type GetHomeworks struct {
	Acquisition time.Time        `json:"acquisition"`
	Homeworks   []HomeworkStruct `json:"homeworks"`
}

// HomeworkStruct は1つの課題情報を収納する構造体
type HomeworkStruct struct {
	Subject string    `json:"subject"`
	Omitted string    `json:"omitted"`
	Name    string    `json:"name"`
	ID      string    `json:"id"`
	Due     time.Time `json:"due"`
}

// SetSchedule はTimeTreeに追加するスケジュールを格納する構造体
type SetSchedule struct {
	Data SetScheduleData `json:"data"`
}

// SetScheduleData はTimeTreeに追加するスケジュールの情報を格納する構造体
type SetScheduleData struct {
	Attributes    SetScheduleAttributes    `json:"attributes"`
	Relationships SetScheduleRelationships `json:"relationships"`
}

// SetScheduleAttributes は追加するスケジュールの属性を格納する構造体
type SetScheduleAttributes struct {
	Category      string    `json:"category"`
	Title         string    `json:"title"`
	AllDay        bool      `json:"all_day"`
	StartAt       time.Time `json:"start_at"`
	StartTimezone string    `json:"start_timezone"`
	EndAt         time.Time `json:"end_at"`
	EndTimezone   string    `json:"end_timezone"`
	Description   string    `json:"description"`
}

// SetScheduleRelationships は追加するスケジュールの関連するものを格納する構造体
type SetScheduleRelationships struct {
	Label SetScheduleLabel `json:"label"`
}

// SetScheduleLabel は追加するスケジュールのラベルを格納する構造体
type SetScheduleLabel struct {
	Data SetScheduleLabelData `json:"data"`
}

// SetScheduleLabelData は追加するスケジュールのラベルのデータを格納する構造体
type SetScheduleLabelData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ResSetSchedule はスケジュール追加リクエストを出したときに返ってくるデータを収納する構造体
type ResSetSchedule struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Category      string        `json:"category"`
			Title         string        `json:"title"`
			AllDay        bool          `json:"all_day"`
			StartAt       time.Time     `json:"start_at"`
			StartTimezone string        `json:"start_timezone"`
			EndAt         time.Time     `json:"end_at"`
			EndTimezone   string        `json:"end_timezone"`
			Recurrences   []interface{} `json:"recurrences"`
			Description   string        `json:"description"`
			Location      string        `json:"location"`
			URL           string        `json:"url"`
			UpdatedAt     time.Time     `json:"updated_at"`
			CreatedAt     time.Time     `json:"created_at"`
		} `json:"attributes"`
		Relationships struct {
			Creator struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"creator"`
			Label struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"label"`
			Attendees struct {
				Data []struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"attendees"`
		} `json:"relationships"`
	} `json:"data"`
}

func init() {
	var err error

	// GAEはタイムゾーン指定できないので、Go側でタイムゾーンを指定する
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

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

	// Firebaseから最後に起動していた時の課題情報を取得
	dbGetKadai(ctx, client)
}
