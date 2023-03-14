package goal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/auth"
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

func newRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.SetHTMLTemplate(newTemplate())
	router.Use(authHandler())

	router.GET("signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.htm", gin.H{})
	})
	// router.POST("signin", signinHandler())
	// router.GET("signout", signoutHandler())
	// router.GET("/", homeHandler())
	return router
}
