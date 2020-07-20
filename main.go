package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 毎時指定した時間に TCJ2 Kadai Store API から課題一覧を取得
	go getRegularly([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45, 47, 49, 51, 53, 55, 57, 59})

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/line-callback", lineRequestFunc)
	r.GET("/version", versionRequestFunc)

	r.Run(":8080")
}

// homeRequestFunc は/アクセスされたときの処理
func homeRequestFunc(c *gin.Context) {
	c.String(404, "Please add line-callback or status path.")
}

// lineRequestFunc は/lineアクセス(LINE Webhook)されたときの処理
func lineRequestFunc(c *gin.Context) {
	c.String(200, "callbacked!")
}

// versionRequestFunc は/versionアクセスされたときの処理
func versionRequestFunc(c *gin.Context) {
	c.String(200, "TCJ2 Kadai Store - v0.1.1")
}
