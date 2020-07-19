package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 毎時指定した時間に TCJ2 Kadai Store API から課題一覧を取得
	go getRegularly([]int{1, 6, 11, 16, 21, 26, 31, 36, 41, 46, 51, 56})

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/line-callback", lineRequestFunc)

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
