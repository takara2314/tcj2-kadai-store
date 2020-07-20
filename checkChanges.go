package main

import (
	"fmt"
	"time"
)

// checkChanges は前の課題情報と比べて変更点がないかをチェック
func checkChanges() (newHW []string, updateHW []string, deleteHW []string) {
	var bothHW []string

	// 過去の課題リストと現在の課題リストを比べたとき、
	// どちらかにしかないものが新規・削除である
	newHW, deleteHW, bothHW = checkEitherers()

	// 課題IDごとに変更点を確認
	for _, hwID := range bothHW {
		// 過去の課題情報と現在の課題情報を比べて内容が変更されていたら、内容変更
		if isUpdated := checkInfoUpdate(hwID); isUpdated {
			updateHW = append(updateHW, hwID)
		}
	}
	return
}

// checkEitherers は過去の課題リストと現在の課題リストを比べたとき、
// どちらかにしかないものを返したり、どちらにもあるもの返す関数
func checkEitherers() (newHW []string, deleteHW []string, bothHW []string) {
	var hwListMap map[string]int = make(map[string]int, 0)
	var hwListPastMap map[string]int = make(map[string]int, 0)

	fmt.Println("現在の課題リスト:", hwList)
	fmt.Println("過去の課題リスト:", hwListPast)

	// リスト内検索がしやすいように、リストからマップに変換
	for _, hwID := range hwList {
		hwListMap[hwID] = 1
	}
	for _, hwID := range hwListPast {
		hwListPastMap[hwID] = 1
	}

	// もしどちらかに存在しない場合、追加されたり消されているとみなす
	for _, hwID := range hwList {
		if _, exist := hwListPastMap[hwID]; !exist {
			newHW = append(newHW, hwID)
		} else {
			bothHW = append(bothHW, hwID)
		}
	}
	for _, hwID := range hwListPast {
		if _, exist := hwListMap[hwID]; !exist {
			deleteHW = append(deleteHW, hwID)
		}
	}

	return
}

// checkInfoUpdate は過去の課題情報と現在の課題情報を比べて内容が変更されていたら、trueを返す関数 (TimeTree関連の情報除く)
// 提出期限(要素番号:4)の比較は、ミリ秒誤差が生じるため、両方の差の整数部分が0であれば等しいとする
func checkInfoUpdate(hwID string) bool {
	for i := 0; i < len(hwStatus[hwID])-3; i++ {
		if hwStatus[hwID][i] != hwStatusPast[hwID][i] {
			if i == 4 {
				pastDueTime := hwStatusPast[hwID][i].(time.Time)
				nowDueTime := hwStatus[hwID][i].(time.Time)

				if durationSecond := pastDueTime.Sub(nowDueTime).Seconds(); int(durationSecond) != 0 {
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}
