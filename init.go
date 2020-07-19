package main

import (
	"time"

	_ "github.com/lib/pq"
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

// AddSchedule はTimeTreeに追加するスケジュールを格納する構造体
type AddSchedule struct {
	Data AddScheduleData `json:"data"`
}

// AddScheduleData はTimeTreeに追加するスケジュールの情報を格納する構造体
type AddScheduleData struct {
	Attributes    AddScheduleAttributes    `json:"attributes"`
	Relationships AddScheduleRelationships `json:"relationships"`
}

// AddScheduleAttributes は追加するスケジュールの属性を格納する構造体
type AddScheduleAttributes struct {
	Category      string    `json:"category"`
	Title         string    `json:"title"`
	AllDay        bool      `json:"all_day"`
	StartAt       time.Time `json:"start_at"`
	StartTimezone string    `json:"start_timezone"`
	EndAt         time.Time `json:"end_at"`
	EndTimezone   string    `json:"end_timezone"`
}

// AddScheduleRelationships は追加するスケジュールの関連するものを格納する構造体
type AddScheduleRelationships struct {
	Label AddScheduleLabel `json:"label"`
}

// AddScheduleLabel は追加するスケジュールのラベルを格納する構造体
type AddScheduleLabel struct {
	Data AddScheduleLabelData `json:"data"`
}

// AddScheduleLabelData は追加するスケジュールのラベルのデータを格納する構造体
type AddScheduleLabelData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ResAddSchedule はスケジュール追加リクエストを出したときに返ってくるデータを収納する構造体
type ResAddSchedule struct {
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
}
