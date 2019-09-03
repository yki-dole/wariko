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
	Id   string `form:"id"`
	Pass string `form:"password"`
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
	r.POST("/signin", loginHandler)
	r.GET("signup/error", makeAccountFormErrorHandler)
	r.GET("/", homeHandler)

	r.Run(":" + port)
}
func AccountCheck(id string, pass string) int {
	var one int64
	one = 1
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	ci.Do("SELECT", 0)
	defer ci.Close()
	existed, err := ci.Do("EXISTS", id)
	check(err)
	if existed == one {
		pass_true, err := redis.String(ci.Do("HGET", id, "pass"))
		check(err)
		if pass == pass_true {
			return 1
		}
	}

	return 0

}
func IsUserExist(id string, pass string, name string, sex int) int {
	ci, err := redis.DialURL(os.Getenv("REDIS_URL"))
	check(err)
	defer ci.Close()
	ci.Do("flushall")
	ci.Do("SELECT", 0)
	maked, err := ci.Do("EXISTS", id)
	check(err)
	var i int64
	i = 0
	if maked == i {
		ci.Do("HSET", id, "pass", pass)
		ci.Do("HSET", id, "NN", name)
		ci.Do("HSET", id, "sex", sex)
		ci.Do("HSET", id, "isPartner", "no")
		return 0
	}

	return 1
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
		isExit := IsUserExist(newForm.Id, newForm.Pass, newForm.Name, newForm.Sex)
		if isExit == 1 {
			c.Redirect(301, "/signup/error")

		} else {

			c.Redirect(301, "/signin")
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
func loginHandler(c *gin.Context) {
	var loginData loginForm
	c.Bind(&loginData)
	result := AccountCheck(loginData.Id, loginData.Pass)
	if result == 1 {
		c.HTML(200, "/login.html", nil)
	}
}
