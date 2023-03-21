package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/web"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &web.SigninForm{}
		if err := c.BindJSON(form); err != nil {
			return
		}
		user := auth.User{Username: form.Username}
		if err := db.Where(&user).First(&user).Error; err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": web.ErrInvalidUsernameOrPassword,
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

		id := uuid.New()
		sessionid := strings.ReplaceAll(id.String(), "-", "")

		// save new session to db
		newSession := &auth.Session{
			Key:        sessionid,
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
		web.SetCookieSessionid(c, sessionid, form.RememberMe)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": user.Username,
		})
	}
}

func Signout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if a session found, clear it in db
		if session, found := c.Get("session"); found {
			session := session.(*auth.Session)
			if err := db.Delete(session).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		web.SetCookieSessionid(c, "", false)
		c.JSON(http.StatusOK, gin.H{"code": 0})
	}
}
