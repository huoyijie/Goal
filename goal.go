package goal

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/glebarez/sqlite"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"gorm.io/gorm"
)

var models = []any{
	&auth.User{},
	&auth.Role{},
	&auth.Session{},
}

func Register(modelList ...any) {
	models = append(models, modelList)
}

func Models() []any {
	sort.Slice(models, func(i, j int) bool {
		group1 := util.Group(models[i])
		group2 := util.Group(models[j])
		if group1 < group2 || (group1 == group2 && util.Item(models[i]) < util.Item(models[j])) {
			return true
		}
		return false
	})
	return models
}

var db *gorm.DB
var enforcer *casbin.Enforcer

//go:embed rbac_model.conf
var rbacModel string

func clearSessions(db *gorm.DB) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		util.LogFatal(db.Delete(&auth.Session{}, "expire_date < ?", time.Now()).Error)
	}
}

func OpenDB() (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(util.WorkDir(), "db.sqlite3")), &gorm.Config{})
	util.LogFatal(err)

	util.LogFatal(db.AutoMigrate(Models()...))

	go clearSessions(db)
	return
}

func Run(host string, port int) {
	model, err := model.NewModelFromString(rbacModel)
	util.LogFatal(err)

	db = OpenDB()
	adapter, err := gormadapter.NewAdapterByDB(db)
	util.LogFatal(err)

	enforcer, err = casbin.NewEnforcer(model, adapter)
	util.LogFatal(err)
	util.LogFatal(enforcer.LoadPolicy())

	router := newRouter()
	router.Run(fmt.Sprintf("%s:%d", host, port))
}
