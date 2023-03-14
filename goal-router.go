package goal

import (
	"fmt"
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

func authMiddleware(c *gin.Context) {
	if sessionid, err := c.Cookie("sessionid"); err == nil {
		session := &auth.Session{
			ID: sessionid,
		}
		if err := db.Preload("User").First(session).Error; err == nil && time.Now().Before(session.ExpireDate) {
			c.Set("session", session)
		}
	}
	c.Next()
}

func anonymous(c *gin.Context) bool {
	_, found := c.Get("session")
	return !found
}

func getSession(c *gin.Context) *auth.Session {
	if session, found := c.Get("session"); found {
		session := session.(*auth.Session)
		return session
	}
	return nil
}

func signinRequiredMiddleware(c *gin.Context) {
	if anonymous(c) {
		contentType := c.GetHeader("Content-Type")
		if strings.EqualFold(contentType, "application/json") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Code": ErrUnauthorized,
			})
		} else {
			c.Redirect(http.StatusFound, "/admin/signin")
		}
	}
	c.Next()
}

func authorizeMiddleware(c *gin.Context) {
	action := c.Param("action")
	group := c.Param("group")
	item := c.Param("item")

	obj := fmt.Sprintf("%s.%s", group, item)

	session, found := c.Get("session")
	if found {
		// validate session
		session := session.(*auth.Session)
		// superuser
		if session.User.IsSuperuser {
			c.Next()
			return
		}
		// has permission
		if ok, err := enforcer.Enforce(session.UserID, obj, action); err == nil && ok {
			c.Next()
			return
		}
	}
	c.AbortWithStatus(http.StatusUnauthorized)
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

func signinHandler(c *gin.Context) {
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

	// todo save user last_signin
	// save new session to request context
	newSession.User = user
	c.Set("session", newSession)

	// save new sessionid to cookie
	setCookieSessionid(c, sessionid)

	c.JSON(http.StatusOK, gin.H{
		"Code": 0,
	})
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.SetHTMLTemplate(newTemplate())
	router.Use(authMiddleware)
	adminGroup := router.Group("admin")

	anonymousGroup := adminGroup.Group("")
	anonymousGroup.GET("signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.htm", gin.H{
			"SigninUrl": "/admin/signin",
			"HomeUrl":   "/admin/",
		})
	})
	anonymousGroup.POST("signin", signinHandler)
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

	signinRequiredGroup := adminGroup.Group("", signinRequiredMiddleware)
	signinRequiredGroup.GET("/", func(c *gin.Context) {
		session := getSession(c)
		groups := groupList()
		for _, group := range groups {
			for _, item := range group.Items {
				can := func(act string) bool {
					if session.User.IsSuperuser {
						return true
					}
					obj := strings.ToLower(fmt.Sprintf("%s.%s", group.Name, item.Name))
					ok, err := enforcer.Enforce(session.UserID, obj, act)
					return err == nil && ok
				}
				if can("add") {
					item.CanAdd = true
				}
				if can("delete") {
					item.CanDelete = true
				}
				if can("change") {
					item.CanChange = true
				}
				if can("get") {
					item.CanGet = true
				}
			}
		}
		c.HTML(http.StatusOK, "index.htm", gin.H{
			"Groups": groups,
		})
	})

	modelGroup := signinRequiredGroup.Group("", authorizeMiddleware)
	// 1.`/get/group/item`
	// 2.`/add/group/item`
	// 3.`/change/group/item/1`
	modelGroup.GET("/:action/:group/:item/*id", func(c *gin.Context) {
		action := c.Param("action")
		group := c.Param("group")
		item := c.Param("item")
		tmp := strings.Split(c.Param("id"), "/")
		id := tmp[1]

		badId := len(tmp) != 2
		badChange := action == "change" && id == ""
		badGetOrAdd := (action == "get" || action == "add") && id != ""
		if badId || badChange || badGetOrAdd {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.HTML(http.StatusOK, fmt.Sprintf("%s.htm", action), gin.H{
			"Action": action,
			"Group":  group,
			"Item":   item,
			"Id":     id,
		})
	})
	// 4.`/add/group/item`
	// 5.`/delete/group/item/1`
	// 6.`/change/group/item/1`
	modelGroup.POST("/:action/:group/:item/*id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Code": 0,
		})
	})

	return router
}
