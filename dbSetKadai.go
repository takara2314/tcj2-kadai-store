package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// dbSetKadai は課題情報をFirebaseに保存(更新)する関数
func dbSetKadai(ctx context.Context, client *firestore.Client, hwID string, hwData []interface{}) {
	// コレクションkadaisに課題情報を追加
	_, err := client.Collection("kadais").Doc(hwID).Set(ctx, map[string]interface{}{
		"subject":                  hwData[0],
		"omitted":                  hwData[1],
		"name":                     hwData[2],
		"id":                       hwData[3],
		"due":                      hwData[4],
		"timetree_name":            hwData[5],
		"timetree_j2a_schedule_id": hwData[6],
		"timetree_j2b_schedule_id": hwData[7],
	})
	if err != nil {
		log.Fatalf("エラーが発生しました: %s\n", err)
	}
}
