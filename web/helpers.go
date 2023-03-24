package web

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/admin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"gorm.io/gorm"
)

func Group(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(filepath.Base(t.PkgPath()))
}

func Item(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(t.Name())
}

func Obj(model any) string {
	t := reflect.TypeOf(model).Elem()
	group := strings.ToLower(filepath.Base(t.PkgPath()))
	item := strings.ToLower(t.Name())
	return fmt.Sprintf("%s.%s", group, item)
}

func Actions() []string {
	return []string{"post", "delete", "put", "get"}
}

func Allow(session *auth.Session, obj, act string, enforcer *casbin.Enforcer) bool {
	if session.User.IsSuperuser && (obj != Obj(&admin.OperationLog{}) || act == "get") {
		return true
	}
	if ok, err := enforcer.Enforce(session.Sub(), obj, act); err == nil && ok {
		return true
	}
	return false
}

func AllowAny(session *auth.Session, obj string, enforcer *casbin.Enforcer) bool {
	return Allow(session, obj, "get", enforcer) || Allow(session, obj, "post", enforcer) || Allow(session, obj, "put", enforcer) || Allow(session, obj, "delete", enforcer)
}

func ParseRoleID(roleID string) uint {
	idStr := strings.Split(roleID, "-")[1]
	id, _ := strconv.ParseUint(idStr, 10, 0)
	return uint(id)
}

func FieldKind(field reflect.StructField) string {
	fieldType := field.Type.Name()
	switch field.Type.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldType = "uint"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldType = "int"
	case reflect.Float32, reflect.Float64:
		fieldType = "float"
	}
	return fieldType
}

func GetGoalTag(field reflect.StructField) (secret, hidden bool, preloadField string) {
	goalTag := strings.Split(field.Tag.Get("goal"), ",")
	secret = util.Contains(goalTag, "secret")
	hidden = util.Contains(goalTag, "hidden")
	preloadField = util.GetWithPrefix(goalTag, "preload=")
	return
}

func GetGormTag(field reflect.StructField) (primary, unique bool) {
	gormTag := strings.Split(field.Tag.Get("gorm"), ",")
	primary = util.Contains(gormTag, "primaryKey")
	unique = util.Contains(gormTag, "unique")
	return
}

func GetBindingTag(field reflect.StructField) string {
	return field.Tag.Get("binding")
}

func Reflect(modelType reflect.Type) (secrets, hiddens, preloads, columns []Column) {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldType := FieldKind(field)
		primary, unique := GetGormTag(field)
		validateRule := GetBindingTag(field)
		secret, hidden, preloadField := GetGoalTag(field)

		column := Column{
			Name:         field.Name,
			Type:         fieldType,
			Secret:       secret,
			Hidden:       hidden,
			Primary:      primary,
			Unique:       unique,
			Preload:      preloadField != "",
			PreloadField: preloadField,
			ValidateRule: validateRule,
		}

		if column.Secret {
			secrets = append(secrets, column)
		}
		if column.Hidden {
			hiddens = append(hiddens, column)
		}
		if column.Preload {
			preloads = append(preloads, column)
		}
		columns = append(columns, column)
	}
	return
}

func GetSession(c *gin.Context) *auth.Session {
	if session, found := c.Get("session"); found {
		session := session.(*auth.Session)
		return session
	}
	return nil
}

func SetCookieSessionid(c *gin.Context, sessionid string, rememberMe bool) {
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

func ClearSessions(db *gorm.DB) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		util.LogFatal(db.Delete(&auth.Session{}, "expire_date < ?", time.Now()).Error)
	}
}

func RecordOpLog(db *gorm.DB, c *gin.Context, record any, action string) {
	s, _ := c.Get("session")
	session := s.(*auth.Session)
	pk := reflect.ValueOf(record).Elem().FieldByName("ID")
	opLog := admin.OperationLog{
		UserID:   session.UserID,
		Date:     time.Now(),
		IPAddr:   c.ClientIP(),
		Group:    Group(record),
		Item:     Item(record),
		Action:   action,
		ObjectID: uint(pk.Uint()),
	}
	db.Create(&opLog)
}

func RecordOpLogs(db *gorm.DB, c *gin.Context, ids []uint, action string) {
	s, _ := c.Get("session")
	session := s.(*auth.Session)
	model, _ := c.Get("model")
	var opLogs []admin.OperationLog
	for _, objID := range ids {
		opLogs = append(opLogs, admin.OperationLog{
			UserID:   session.UserID,
			Date:     time.Now(),
			IPAddr:   c.ClientIP(),
			Group:    Group(model),
			Item:     Item(model),
			Action:   action,
			ObjectID: objID,
		})
	}
	db.Create(&opLogs)
}
