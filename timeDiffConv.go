package main

import "time"

// timeDiffConv は時差変換をして返す関数
func timeDiffConv(tTime time.Time) (rTime time.Time) {
	// よりUTCらしくする
	rTime = tTime.UTC()

	// UTC → JST
	var jst *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)
	rTime = rTime.In(jst)

	return
}
