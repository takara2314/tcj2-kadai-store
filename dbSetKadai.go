package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// dbSetKadai は課題情報をFirebaseに保存(更新)する関数
func dbSetKadai(ctx context.Context, client *firestore.Client, hwID string, hwData []interface{}) {
	// コレクションkadaisに課題情報を追加
	_, err := client.Collection("kadais2").Doc(hwID).Set(ctx, map[string]interface{}{
		"course":                   hwData[0],
		"subject":                  hwData[1],
		"subjectID":                hwData[2],
		"name":                     hwData[3],
		"id":                       hwData[4],
		"due":                      hwData[5],
		"timetree_name":            hwData[6],
		"timetree_j2a_schedule_id": hwData[7],
		"timetree_j2b_schedule_id": hwData[8],
	})
	if err != nil {
		log.Fatalf("エラーが発生しました: %v\n", err)
	}
}
