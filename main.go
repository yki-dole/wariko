package main

import (
	"log"
	"net/http"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

type userMakeForm struct {
	Id   string `form:"id"`
	Name string `form:"name"`
	Pass string `form:"password"`
	Sex  int    `form:"sex"`
}
type userForm struct {
	Id   string `form:"user_id"`
	Name string `form:"user_name"`
	Sex  bool   `form:"sex"`
}

type loginForm struct {
	Id   string `form:"user_id"`
	Pass string `form:"pass"`
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
	r.Static("/js", "./js")
	r.Static("/picture", "./picture")
	r.LoadHTMLGlob("views/*")

	r.GET("/signin", loginAcsessHandler)       //アクセス時のハンドラ
	r.GET("/signup", makeAccountAcsessHandler) //アクセス時のハンドラ

	r.POST("/signup", makeAccountHandler) //ログインでPOST投げた時のハンドラ

	r.GET("signup/error", makeAccountFormErrorHandler)
	r.GET("/", homeHandler)

	r.POST("/user", loginFormHandler)

	r.Run(":" + port)
}
func IsUserExist(id string, pass string, name string, sex int) int {
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	defer ci.Close()
	ci.Do("SELECT", 0)
	maked, err := ci.Do("EXISTS", id)
	check(err)
	var i int64
	i = 0
	if maked == i {
		ci.Do("HSET", id, "pass", pass)
		ci.Do("HSET", id, "NN", name)
		ci.Do("HSET", id, "sex", sex)

		return 0
	}

	return 100
}

func check(er error) {
	if er != nil {
		panic(er)
	}
}

func makeAccountAcsessHandler(c *gin.Context) {
	c.HTML(200, "make_form.html", nil)
}
func makeAccountFormErrorHandler(c *gin.Context) {
	c.HTML(200, "make_form_error.html", nil)
}
func homeHandler(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}
func loginAcsessHandler(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}
func makeAccountHandler(c *gin.Context) {
	var newForm userMakeForm
	c.Bind(&newForm)
	if (newForm.Id == "") || (newForm.Pass == "") {
		c.Redirect(301, "/signup/errorfadsada")
	} else {
		c.Redirect(301, "/signin")
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
