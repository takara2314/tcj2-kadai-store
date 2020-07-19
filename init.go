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

func init() {
	var err error

	// GAEはタイムゾーン指定できないので、Go側でタイムゾーンを指定する
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}
