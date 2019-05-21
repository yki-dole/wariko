package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type userForm struct {
	Name string `form:"user_name"`
	Pass string `form:"pass"`
}
type urlData struct {
	Url     string `form:"add_url"`
	Type    string `form:"type"`
	Title   string `form:"url_title"`
	TextKey string
}

type urlInputForm struct {
	Url   string `form:"add_url"`
	Type  string `form:"type"`
	Title string `form:"url_title"`
	Text  string `form:"url_text"`
	Value int    `form:"value"`
}

var user userForm

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	r := gin.Default()
	// redis-server でサーバーを起動
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "error404.html", nil)
	})

	//redisに接続
	r.LoadHTMLGlob("views/*")
	r.GET("/", homeHandler)

	r.Run(":" + port)
}
func homeHandler(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}
