package main

import (
	"time"
)

// checkChanges は前の課題情報と比べて変更点がないかをチェック
func checkChanges() (newHW []string, updateHW []string, deleteHW []string) {
	// 前の課題の数が0なら、すべて新規追加確定
	if len(hwListPast) == 0 {
		// 今の課題情報を過去のものにする
		hwStatusPast = hwStatus
		hwListPast = hwList
		// 現在の課題情報は全て新規追加
		newHW = hwList

		return
	}

	// 課題IDごとに変更点を確認
	for _, hwID := range hwList {
		// 過去の課題情報に、現在のIDが存在していなかったら、新規追加
		// 現在の課題情報に、過去のIDが存在していなかったら、削除

		// 過去の課題情報と現在の課題情報を比べて内容が変更されていたら、内容変更
		if isUpdated := checkInfoUpdate(hwID); isUpdated {
			updateHW = append(updateHW, hwID)
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
