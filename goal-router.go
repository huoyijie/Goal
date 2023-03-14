package goal

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"golang.org/x/crypto/bcrypt"
)

type Code int

const (
	ErrInvalidUsernameOrPassword Code = -(iota + 10000)
)

func authHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessionid, err := c.Cookie("sessionid"); err == nil {
			session := &auth.Session{
				ID: sessionid,
			}
			if err := db.First(session).Error; err == nil && time.Now().Before(session.ExpireDate) {
				c.Set("user", session.User)
			}
		}
		c.Next()
	}
}

func setCookieSessionid(c *gin.Context, sessionid string) {
	expireIn := 3 * 24 * 60 * 60
	if len(sessionid) == 0 {
		expireIn = -1
	}
	c.SetCookie("sessionid", sessionid, expireIn, "/", "127.0.0.1", false, true)
}

type SigninForm struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=40"`
	Password string `json:"password" binding:"required"`
}

func signinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &SigninForm{}
		if err := c.BindJSON(form); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		user := &auth.User{Username: form.Username}
		if err := db.First(user).Error; err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, gin.H{
				"Code": ErrInvalidUsernameOrPassword,
			})
			return
		}

		id, err := uuid.NewUUID()
		util.LogFatal(err)
		sessionid := strings.ToLower(strings.ReplaceAll(id.String(), "-", ""))
		setCookieSessionid(c, sessionid)

		c.JSON(http.StatusOK, gin.H{
			"Code": 0,
		})
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.SetHTMLTemplate(newTemplate())
	router.Use(authHandler())

	router.GET("signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.htm", gin.H{})
	})
	router.POST("signin", signinHandler())
	router.GET("signout", func(c *gin.Context) {
		setCookieSessionid(c, "")
		c.JSON(http.StatusOK, gin.H{"Code": 0})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.htm", gin.H{})
	})
	return router
}
