package goal

import (
	_ "embed"
	"sort"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/admin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"github.com/huoyijie/goal/web"
	"github.com/huoyijie/goal/web/handlers"
	"github.com/huoyijie/goal/web/middlewares"
	"gorm.io/gorm"
)

//go:embed rbac_model.conf
var rbacModel string

type Goal interface {
	NewSuper(*auth.User)
	Router() *gin.Engine
}

func NewGoal(db *gorm.DB, models ...any) Goal {
	m := []any{
		&auth.User{},
		&auth.Role{},
		&auth.Session{},
		&admin.OperationLog{},
	}
	m = append(m, models...)

	model, err := model.NewModelFromString(rbacModel)
	util.LogFatal(err)

	util.LogFatal(db.AutoMigrate(m...))

	adapter, err := gormadapter.NewAdapterByDB(db)
	util.LogFatal(err)

	enforcer, err := casbin.NewEnforcer(model, adapter)
	util.LogFatal(err)
	util.LogFatal(enforcer.LoadPolicy())

	goal := &goal_web_t{
		db:       db,
		enforcer: enforcer,
		models:   m,
	}
	return goal
}

type goal_web_t struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	models   []any
}

// NewSuper implements Goal
func (gw *goal_web_t) NewSuper(super *auth.User) {
	util.LogFatal(gw.db.Create(super).Error)
}

// Router implements Goal
func (gw *goal_web_t) Router() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(middlewares.Cors())
	router.Use(middlewares.Auth(gw.db))
	// `/admin`
	adminGroup := router.Group("admin")

	anonymousGroup := adminGroup.Group("")
	// `/admin/signin`
	anonymousGroup.POST("signin", handlers.Signin(gw.db))
	// `/admin/signout`
	anonymousGroup.GET("signout", handlers.Signout(gw.db))

	signinRequiredGroup := adminGroup.Group("", middlewares.SigninRequired)
	// `/admin/menus`
	signinRequiredGroup.GET("menus", handlers.Menus(gw.getModels(), gw.enforcer))

	// `/admin/perms`
	permsGroup := signinRequiredGroup.Group("perms", middlewares.CanChangePerms(gw.enforcer))
	// `/admin/perms/:roleID`
	permsGroup.GET(":roleID", handlers.GetPerms(gw.getModels(), gw.enforcer))
	permsGroup.PUT(":roleID", handlers.ChangePerms(gw.db, gw.enforcer))

	// `/admin/roles`
	rolesGroup := signinRequiredGroup.Group("roles", middlewares.CanChangeRoles(gw.enforcer))
	// `/admin/roles/:userID`
	rolesGroup.GET(":userID", handlers.GetRoles(gw.db, gw.enforcer))
	rolesGroup.PUT(":userID", handlers.ChangeRoles(gw.db, gw.enforcer))

	// `/admin/crud`
	crudGroup := signinRequiredGroup.Group("crud", middlewares.ValidateModel(gw.getModels()))
	// `/admin/crud/:group/:item`
	modelGroup := crudGroup.Group(":group/:item")

	// `/admin/crud/:group/:item/perms`
	modelGroup.GET("perms", handlers.CrudPerms(gw.enforcer))
	modelGroup.GET("mine", handlers.CrudGetMine(gw.db))

	AuthorizeGroup := modelGroup.Group("", middlewares.Authorize(gw.enforcer))
	// `/admin/crud/:group/:item`
	AuthorizeGroup.GET("", handlers.CrudGet(gw.db))
	AuthorizeGroup.POST("", handlers.CrudPost(gw.db, gw.enforcer))
	AuthorizeGroup.PUT("", handlers.CrudPut(gw.db, gw.enforcer))
	AuthorizeGroup.DELETE("", handlers.CrudDelete(gw.db, gw.enforcer))
	AuthorizeGroup.DELETE("batch", handlers.CrudBatchDelete(gw.db, gw.enforcer))
	AuthorizeGroup.POST("exist", handlers.CrudExist(gw.db))

	go web.ClearSessions(gw.db)
	return router
}

func (gw *goal_web_t) getModels() []any {
	models := make([]any, len(gw.models))
	copy(models, gw.models)
	sort.Slice(models, func(i, j int) bool {
		group1 := web.Group(models[i])
		group2 := web.Group(models[j])
		if group1 < group2 || (group1 == group2 && web.Item(models[i]) < web.Item(models[j])) {
			return true
		}
		return false
	})
	return models
}

var _ Goal = (*goal_web_t)(nil)
