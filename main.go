package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 毎時指定した時間に TCJ2 Kadai Store API から課題一覧を取得
	go getRegularly([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 24, 25, 26, 27, 28, 29, 30, 31, 36, 38, 41, 46, 49, 51, 56})

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/line-callback", lineRequestFunc)

	r.Run()
}

// homeRequestFunc は/アクセスされたときの処理
func homeRequestFunc(c *gin.Context) {
	c.String(404, "Please add line-callback or status path.")
}

// lineRequestFunc は/lineアクセス(LINE Webhook)されたときの処理
func lineRequestFunc(c *gin.Context) {
	c.String(200, "callbacked!")
}
