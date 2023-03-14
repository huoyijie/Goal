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
	ErrUnauthorized
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessionid, err := c.Cookie("sessionid"); err == nil {
			session := &auth.Session{
				ID: sessionid,
			}
			if err := db.First(session).Error; err == nil && time.Now().Before(session.ExpireDate) {
				c.Set("session", session)
			}
		}
		c.Next()
	}
}

func anonymous(c *gin.Context) bool {
	_, found := c.Get("session")
	return !found
}

func signinRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if anonymous(c) {
			contentType := c.GetHeader("Content-Type")
			if strings.EqualFold(contentType, "application/json") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Code": ErrUnauthorized,
				})
			} else {
				c.Redirect(http.StatusFound, "/signin")
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
		user := auth.User{Username: form.Username}
		if err := db.First(&user).Error; err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, gin.H{
				"Code": ErrInvalidUsernameOrPassword,
			})
			return
		}

		// if a session found, clear it in db
		if oldSession, found := c.Get("session"); found {
			oldSession := oldSession.(*auth.Session)
			if err := db.Delete(oldSession).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		id, err := uuid.NewUUID()
		util.LogFatal(err)
		sessionid := strings.ToLower(strings.ReplaceAll(id.String(), "-", ""))

		// save new session to db
		newSession := &auth.Session{
			ID:         sessionid,
			UserID:     user.ID,
			ExpireDate: time.Now().Add(3 * 24 * time.Hour),
		}
		if err := db.Create(newSession).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// save new session to request context
		newSession.User = user
		c.Set("session", newSession)

		// save new sessionid to cookie
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
	router.Use(authMiddleware())

	anonymousGroup := router.Group("")
	anonymousGroup.GET("signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.htm", gin.H{})
	})
	anonymousGroup.POST("signin", signinHandler())
	anonymousGroup.GET("signout", func(c *gin.Context) {
		// if a session found, clear it in db
		if session, found := c.Get("session"); found {
			session := session.(*auth.Session)
			if err := db.Delete(session).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		setCookieSessionid(c, "")
		c.JSON(http.StatusOK, gin.H{"Code": 0})
	})

	authGroup := router.Group("", signinRequiredMiddleware())
	authGroup.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.htm", gin.H{})
	})
	return router
}
