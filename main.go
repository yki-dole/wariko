package main

import (
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type userMakeForm struct {
	Id   string `form:"user_id"`
	Name string `form:"user_name"`
	Pass string `form:"pass"`
	Sex  bool   `form:"sex"`
}
type userForm struct {
	Id   string `form:"user_id"`
	Name string `form:"user_name"`
	Sex  bool   `form:"sex"`
}
type loginForm struct {
	Id string `form:"user_id"`

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
	r.GET("/signin", loginFormHandler)
	r.GET("/", homeHandler)

	r.GET("/makeform", makeFormHandler)

	r.POST("/user", loginFormHandler)
	r.POST("/signup", makeAccountHandler)

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
	var newForm userMakeForm
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	c.Bind(&newForm)
	defer ci.Close()
	if (newForm.Id == "") || (newForm.Pass == "") {
		c.HTML(200, "make_form.html", gin.H{
			"errortxt": "error:Prease fill out",
		})
	} else {
		ci.Do("SELECT", 0)
		maked, err := ci.Do("EXISTS", newForm.Id)
		check(err)
		var i int64
		i = 0
		if maked == i {
			ci.Do("HSET", newForm.Id, "pass", newForm.Pass)
			ci.Do("HSET", newForm.Id, "NN", newForm.Name)
			ci.Do("HSET", newForm.Id, "sex", newForm.Sex)
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

	var fakeForm loginForm
	c.Bind(&fakeForm)
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	ci.Do("SELECT", 0)

	defer ci.Close()
	maked, err := ci.Do("EXISTS", fakeForm.Id)
	if err != nil {
		panic(err)
	}
	var i int64
	i = 1
	if maked == i {
		truePass, err := redis.String(ci.Do("HGET", fakeForm.Id, "pass"))
		check(err)
		if truePass == fakeForm.Pass {
			text := "Hello " + fakeForm.Id + " !!"

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
