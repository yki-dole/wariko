package main

import (
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
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
	r.Static("/css", "./css")
	r.Static("/picture", "./picture")
	r.LoadHTMLGlob("views/*")
	r.GET("/login", indexHandler)
	r.GET("/", homeHandler)

	r.GET("/makeform", makeFormHandler)

	r.POST("/user", loginFormHandler)
	r.POST("/makeaccount", makeAccountHandler)

	r.Run(":" + port)
}

func check(er error) {
	if er != nil {
		panic(er)
	}
}
func makeFormHandler(c *gin.Context) {
	c.HTML(200, "make_form.html", gin.H{
		"errortxt": "",
	})
}
func homeHandler(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}
func indexHandler(c *gin.Context) {
	c.HTML(200, "form.html", gin.H{
		"errortxt": "",
	})
}
func makeAccountHandler(c *gin.Context) {
	var newForm userForm
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	c.Bind(&newForm)
	defer ci.Close()
	if (newForm.Name == "") || (newForm.Pass == "") {
		c.HTML(200, "make_form.html", gin.H{
			"errortxt": "error:Prease fill out",
		})
	} else {
		ci.Do("SELECT", 0)
		maked, err := ci.Do("EXISTS", newForm.Name)
		check(err)
		var i int64
		i = 0
		if maked == i {
			ci.Do("HSET", newForm.Name, "pass", newForm.Pass)
			urlListSize := 0
			ci.Do("HSET", newForm.Name, "urlListSize", urlListSize)
			c.HTML(200, "form.html", nil)
		} else {
			c.HTML(200, "make_form.html", gin.H{
				"errortxt": "error",
			})
		}

	}
}

func userHandler(c *gin.Context) {

	if user.Name == "" {
		c.HTML(200, "form.html", gin.H{
			"errortxt": "error:Prease login",
		})
	} else {
		text := "Hello  !!"
		c.HTML(200, "user.html", gin.H{

			"name": text,
		})
	}

}
func loginFormHandler(c *gin.Context) {

	var fakeForm userForm
	c.Bind(&fakeForm)
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	ci.Do("SELECT", 0)

	defer ci.Close()
	maked, err := ci.Do("EXISTS", fakeForm.Name)
	if err != nil {
		panic(err)
	}
	var i int64
	i = 1
	if maked == i {
		truePass, err := redis.String(ci.Do("HGET", fakeForm.Name, "pass"))
		check(err)
		if truePass == fakeForm.Pass {
			text := "Hello " + fakeForm.Name + " !!"
			user.Name = fakeForm.Name

			c.HTML(200, "user.html", gin.H{
				"name": text,
			})
		} else {
			c.HTML(200, "form.html", gin.H{
				"errortxt": truePass,
			})
		}
	} else {
		c.HTML(200, "form.html", gin.H{
			"errortxt": "error:There is not it name",
		})
	}
}
