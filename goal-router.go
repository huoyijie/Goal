package goal

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Code int

const (
	ErrInvalidUsernameOrPassword Code = -(iota + 10000)
)

func authMiddleware(c *gin.Context) {
	if sessionid, err := c.Cookie("g_sessionid"); err == nil {
		session := &auth.Session{
			Key: sessionid,
		}
		if err := db.Preload("User").Where(session).First(session).Error; err == nil && time.Now().Before(session.ExpireDate) {
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
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

func authorizeMiddleware(c *gin.Context) {
	action := strings.ToLower(c.Request.Method)

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
		model, _ := c.Get("model")
		if ok, err := enforcer.Enforce(session.Sub(), util.Obj(model), action); err == nil && ok {
			c.Next()
			return
		}
	}
	c.AbortWithStatus(http.StatusForbidden)
}

func validateModelMiddleware(c *gin.Context) {
	group := c.Param("group")
	item := c.Param("item")

	var model any
	var modelType reflect.Type
	for _, m := range Models() {
		if strings.EqualFold(group, util.Group(m)) && strings.EqualFold(item, util.Item(m)) {
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

func setCookieSessionid(c *gin.Context, sessionid string, rememberMe bool) {
	// keep g_sessionid until the browser closed
	maxAge := 0

	if len(sessionid) == 0 {
		// sign out: delete g_sessionid right now
		maxAge = -1
	} else if rememberMe {
		// sign in: remember me was checked
		maxAge = 3 * 24 * 60 * 60
	}
	c.SetCookie("g_sessionid", sessionid, maxAge, "/", "127.0.0.1", false, true)
}

type SigninForm struct {
	Username   string `json:"username" binding:"required,alphanum,min=3,max=40"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
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
			"code": ErrInvalidUsernameOrPassword,
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
		Key:        sessionid,
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
	setCookieSessionid(c, sessionid, form.RememberMe)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user.Username,
	})
}

type Column struct {
	Name    string
	Type    string
	Primary bool
}

func crud(c *gin.Context, op byte) {
	mt, _ := c.Get("modelType")
	modelType := mt.(reflect.Type)

	record := reflect.New(modelType).Interface()
	if err := c.BindJSON(record); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var tx *gorm.DB
	switch op {
	case 1:
		tx = db.Create(record)
	case 2:
		tx = db.Save(record)
	case 3:
		tx = db.Delete(record)
	}

	if err := tx.Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": record,
	})
}

type Menu struct {
	Name  string `json:"label"`
	Items []Menu `json:"items"`
}

type Perm struct {
	Code,
	Name string
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:4000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           10 * time.Minute,
	}))
	router.Use(authMiddleware)
	adminGroup := router.Group("admin")

	anonymousGroup := adminGroup.Group("")
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
		setCookieSessionid(c, "", false)
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	signinRequiredGroup := adminGroup.Group("", signinRequiredMiddleware)
	signinRequiredGroup.GET("/menus", func(c *gin.Context) {
		session := getSession(c)
		menus := []Menu{}
		for _, m := range Models() {
			if util.AllowAny(session, util.Obj(m), enforcer) {
				var found bool
				group := util.Group(m)
				for i := range menus {
					if menus[i].Name == group {
						found = true
						menus[i].Items = append(menus[i].Items, Menu{util.Item(m), nil})
						break
					}
				}
				if !found {
					menus = append(menus, Menu{group, []Menu{{util.Item(m), nil}}})
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": menus,
		})
	})
	signinRequiredGroup.GET("/perms/:roleId", func(c *gin.Context) {
		if util.Allow(getSession(c), util.Obj(&auth.Role{}), "put", enforcer) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}, func(c *gin.Context) {
		var perms []Perm
		for _, m := range Models() {
			for _, act := range util.Actions() {
				perms = append(perms, Perm{
					Code: fmt.Sprintf("%s:%s", util.Obj(m), act),
					Name: fmt.Sprintf("%s %s %s", act, util.Group(m), util.Item(m)),
				})
			}
		}
		// roleId, _ := c.Params.Get("roleId");
		// enforcer
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": perms,
		})
	})

	modelGroup := signinRequiredGroup.Group("", validateModelMiddleware, authorizeMiddleware)

	modelGroup.GET("/:group/:item", func(c *gin.Context) {
		model, _ := c.Get("model")
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)

		var hiddens []string
		var columns []Column
		var preloads [][]string
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			goalTags := strings.Split(field.Tag.Get("goal"), ",")
			if util.Contains(goalTags, "hidden") {
				hiddens = append(hiddens, field.Name)
			} else {
				gormTags := strings.Split(field.Tag.Get("gorm"), ",")
				primary := util.Contains(gormTags, "primaryKey")
				columns = append(columns, Column{field.Name, field.Type.Name(), primary})
			}
			if fname := util.GetWithPrefix(goalTags, "preload="); fname != "" {
				preloads = append(preloads, []string{field.Name, fname})
			}
		}

		records := reflect.New(reflect.SliceOf(modelType)).Interface()
		tx := db.Model(model)
		for _, preload := range preloads {
			tx = tx.Preload(preload[0])
		}
		if err := tx.Find(records).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if len(hiddens) > 0 || len(preloads) > 0 {
			recordsVal := reflect.ValueOf(records).Elem()
			for i := 0; i < recordsVal.Len(); i++ {
				recordVal := recordsVal.Index(i)
				for _, hidden := range hiddens {
					hiddenField := recordVal.FieldByName(hidden)
					hiddenField.SetZero()
				}
				for _, preload := range preloads {
					preloadField := recordVal.FieldByName(preload[0])
					dstFF := preloadField.FieldByName(preload[1])
					dstVal := dstFF.Interface()
					preloadField.SetZero()
					dstFF.Set(reflect.ValueOf(dstVal))
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"records":  records,
				"columns":  columns,
				"preloads": preloads,
			},
		})
	})
	modelGroup.POST("/:group/:item", func(c *gin.Context) {
		crud(c, 1)
	})
	modelGroup.PUT("/:group/:item", func(c *gin.Context) {
		crud(c, 2)
	})
	modelGroup.DELETE("/:group/:item", func(c *gin.Context) {
		crud(c, 3)
	})
	modelGroup.DELETE("/:group/:item/batch", func(c *gin.Context) {
		model, _ := c.Get("model")

		ids := []uint{}
		if err := c.BindJSON(&ids); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := db.Delete(model, ids).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
	})
	return router
}
