package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/Goal/util"
	"github.com/huoyijie/Goal/web"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signin(db *gorm.DB, domain string, secure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &web.SigninForm{}
		if err := c.BindJSON(form); err != nil {
			return
		}
		user := auth.User{Username: form.Username, IsActive: true}
		if err := db.Where(&user).First(&user).Error; err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, web.Result{Code: web.ErrInvalidUsernameOrPassword})
			return
		}

		// if a session found, clear it in db
		if session := web.GetSession(c); session != nil {
			if err := db.Delete(session).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		id := uuid.New()
		sessionid := strings.ReplaceAll(id.String(), "-", "")

		// save new session to db
		newSession := &auth.Session{
			Key:        sessionid,
			User:       user,
			ExpireDate: time.Now().Add(3 * 24 * time.Hour),
		}
		if err := db.Create(newSession).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// save new session to request context
		c.Set("session", newSession)

		// save new sessionid to cookie
		web.SetCookieSessionid(c, sessionid, form.RememberMe, domain, secure)

		c.JSON(http.StatusOK, web.Result{Data: user.Username})
	}
}

func Signout(db *gorm.DB, domain string, secure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if a session found, clear it in db
		if session, found := c.Get("session"); found {
			session := session.(*auth.Session)
			if err := db.Delete(session).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
		web.SetCookieSessionid(c, "", false, domain, secure)
		c.JSON(http.StatusOK, web.Result{})
	}
}

func Userinfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := web.GetSession(c)
		c.JSON(http.StatusOK, web.Result{Data: session.User})
	}
}

func ChangePassword(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := &web.ChangePasswordForm{}
		if err := c.BindJSON(form); err != nil {
			return
		}
		session := web.GetSession(c)
		if bcrypt.CompareHashAndPassword([]byte(session.User.Password), []byte(form.Password)) != nil {
			c.JSON(http.StatusOK, web.Result{Code: web.ErrInvalidPassword})
			return
		}
		if err := db.Model(&session.User).Updates(&auth.User{Password: util.BcryptHash(form.NewPassword)}).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, web.Result{})
	}
}
