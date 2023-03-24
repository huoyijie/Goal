package middlewares

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/web"
	"gorm.io/gorm"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:4000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           10 * time.Minute,
	})
}

func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessionid, err := c.Cookie("g_sessionid"); err == nil {
			session := &auth.Session{
				Key: sessionid,
			}
			if err := db.Joins("User").Where(session).First(session).Error; err == nil && time.Now().Before(session.ExpireDate) {
				c.Set("session", session)
			}
		}
		c.Next()
	}
}

func SigninRequired(c *gin.Context) {
	if _, found := c.Get("session"); !found {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

func ValidateModel(models []any) gin.HandlerFunc {
	return func(c *gin.Context) {
		group := c.Param("group")
		item := c.Param("item")

		var model any
		var modelType reflect.Type
		for _, m := range models {
			if strings.EqualFold(group, web.Group(m)) && strings.EqualFold(item, web.Item(m)) {
				model = m
				modelType = reflect.TypeOf(m).Elem()
				break
			}
		}

		if model == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Set("model", model)
		c.Set("modelType", modelType)
		c.Next()
	}
}

func Authorize(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if session, found := c.Get("session"); found {
			session := session.(*auth.Session)
			model, _ := c.Get("model")
			action := strings.ToLower(c.Request.Method)

			if web.Allow(session, web.Obj(model), action, enforcer) {
				c.Next()
				return
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func CanChangePerms(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if web.Allow(web.GetSession(c), web.Obj(&auth.Role{}), "put", enforcer) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func CanChangeRoles(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if web.Allow(web.GetSession(c), web.Obj(&auth.User{}), "put", enforcer) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
