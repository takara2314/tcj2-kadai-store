package main

// checkChanges は前の課題情報と比べて変更点がないかをチェック
func checkChanges() (newHW []string, updateHW []string, deleteHW []string) {
	// 前の課題の数が0なら、新規追加確定
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
		// 過去の課題情報に、対象のIDが存在していなかったら、新規追加
		if _, exist := hwStatusPast[hwID]; !exist {
			newHW = append(newHW, hwID)
		}
	}
	return
}
