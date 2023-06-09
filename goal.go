package goal

import (
	_ "embed"
	"sort"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/Goal/admin"
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/Goal/util"
	"github.com/huoyijie/Goal/web"
	"github.com/huoyijie/Goal/web/handlers"
	"github.com/huoyijie/Goal/web/middlewares"
	"gorm.io/gorm"
)

//go:embed rbac_model.conf
var rbacModel string

type Cookie struct {
	Domain string
	Secure bool
}

type Config struct {
	AllowOrigins, TrustedProxies []string
	Cookie
}

type Goal interface {
	CreateSuper(*auth.User)
	Router() *gin.Engine
}

func New(config Config, db *gorm.DB, models ...any) Goal {
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
		Config:   config,
		db:       db,
		enforcer: enforcer,
		models:   m,
	}
	admin.AddItems(goal.getModels())
	return goal
}

type goal_web_t struct {
	Config
	enforcer *casbin.Enforcer
	db       *gorm.DB
	models   []any
}

// CreateSuper implements Goal
func (gw *goal_web_t) CreateSuper(super *auth.User) {
	util.LogFatal(gw.db.Create(super).Error)
}

// Router implements Goal
func (gw *goal_web_t) Router() *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.Cors(gw.AllowOrigins))
	router.Use(middlewares.Auth(gw.db))
	// `/admin`
	adminGroup := router.Group("admin")

	anonymousGroup := adminGroup.Group("")
	// `/admin/signin`
	anonymousGroup.POST("signin", handlers.Signin(gw.db, gw.Domain, gw.Secure))
	// `/admin/signout`
	anonymousGroup.GET("signout", handlers.Signout(gw.db, gw.Domain, gw.Secure))
	anonymousGroup.GET("locale", handlers.Translate(gw.getModels()))

	signinRequiredGroup := adminGroup.Group("", middlewares.SigninRequired)
	// `/admin/menus`
	signinRequiredGroup.GET("menus", handlers.Menus(gw.getModels(), gw.enforcer))
	signinRequiredGroup.GET("userinfo", handlers.Userinfo())
	signinRequiredGroup.POST("changepw", handlers.ChangePassword(gw.db))

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

	// `/admin/crud/:group/:item/columns`
	modelGroup.GET("datatable", handlers.CrudDataTable(gw.enforcer))
	modelGroup.GET("mine", handlers.CrudGetMine(gw.db))

	AuthorizeGroup := modelGroup.Group("", middlewares.Authorize(gw.enforcer))
	// `/admin/crud/:group/:item`
	AuthorizeGroup.GET("", handlers.CrudGet(gw.db))
	AuthorizeGroup.POST("", handlers.CrudPost(gw.db, gw.enforcer))
	AuthorizeGroup.PUT("", handlers.CrudPut(gw.db, gw.enforcer))
	AuthorizeGroup.DELETE("", handlers.CrudDelete(gw.db, gw.enforcer))
	AuthorizeGroup.DELETE("batch", handlers.CrudBatchDelete(gw.db, gw.enforcer))
	AuthorizeGroup.POST("exist", handlers.CrudExist(gw.db))
	AuthorizeGroup.POST("upload/:field", handlers.Upload(gw.db))
	AuthorizeGroup.GET("select/:field", handlers.Select(gw.db))

	go web.ClearSessions(gw.db)
	router.SetTrustedProxies(gw.TrustedProxies)
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
