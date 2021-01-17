package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// dbGetKadai はFirebaseから課題情報を取得する関数
func dbGetKadai(ctx context.Context, client *firestore.Client) {
	var oneHwStatus map[string]interface{}

	iter := client.Collection("kadais2").Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}

		// 今のドキュメントのデータをmapとして取得する
		oneHwStatus = doc.Data()
		var nowID string = oneHwStatus["id"].(string)

		// 過去の課題情報としてメモリに格納する
		// ※時刻データはUTCで返ってくるので、JSTに変換する
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["course"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["subject"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["subjectID"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["name"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["id"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], timeDiffConv(oneHwStatus["due"].(time.Time)))
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["timetree_name"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["timetree_j2a_schedule_id"])
		hwStatusPast[nowID] = append(hwStatusPast[nowID], oneHwStatus["timetree_j2b_schedule_id"])

		// 過去の課題リスト(ID)としてメモリに格納する
		hwListPast = append(hwListPast, nowID)
	}
}
