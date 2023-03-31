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
	"github.com/huoyijie/goal/web/tag"
	"gorm.io/gorm"
)

const (
	PASSWORD_PLACEHOLDER string = "password placeholder"
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

func GetComponent(field reflect.StructField) tag.Component {
	token := field.Tag.Get("goal")
	if token != "" {
		for _, c := range tag.COMPONETS {
			if c.Is(token) {
				com := reflect.New(reflect.TypeOf(c).Elem())
				m := com.MethodByName("Unmarshal")
				m.Call([]reflect.Value{reflect.ValueOf(token)})
				return com.Interface().(tag.Component)
			}
		}
	}
	return nil
}

func GetBindingTag(field reflect.StructField) string {
	return field.Tag.Get("binding")
}

func Reflect(modelType reflect.Type) (secrets, preloads, columns []Column) {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		validateRule := GetBindingTag(field)

		component := GetComponent(field)
		if component == nil {
			continue
		}

		column := Column{
			field.Name,
			Component{
				component.Head(),
				component,
			},
			validateRule,
		}

		if component.(tag.IBase).Get().Autowired {
			continue
		}
		if component.(tag.IBase).Get().Secret {
			secrets = append(secrets, column)
		}

		if component.(tag.IBase).Get().BelongTo != nil {
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
	session := GetSession(c)
	pk := reflect.ValueOf(record).Elem().FieldByName("ID")
	opLog := admin.OperationLog{
		UserID:   session.UserID,
		Date:     time.Now(),
		IP:       c.ClientIP(),
		Group:    Group(record),
		Item:     Item(record),
		Action:   action,
		ObjectID: uint(pk.Uint()),
	}
	db.Create(&opLog)
}

func RecordOpLogs(db *gorm.DB, c *gin.Context, ids []uint, action string) {
	session := GetSession(c)
	model, _ := c.Get("model")
	var opLogs []admin.OperationLog
	for _, objID := range ids {
		opLogs = append(opLogs, admin.OperationLog{
			UserID:   session.UserID,
			Date:     time.Now(),
			IP:       c.ClientIP(),
			Group:    Group(model),
			Item:     Item(model),
			Action:   action,
			ObjectID: objID,
		})
	}
	db.Create(&opLogs)
}

func AutowiredCreator(c *gin.Context, action string, record any) {
	if action == "post" {
		creatorField := reflect.ValueOf(record).Elem().FieldByName("Creator")
		if creatorField.IsValid() {
			session := GetSession(c)
			creatorField.SetUint(uint64(session.UserID))
		}
	}
}

func SecureRecords(secrets []Column, recordsVal reflect.Value) {
	for i := 0; i < recordsVal.Len(); i++ {
		SecureRecord(secrets, recordsVal.Index(i))
	}
}

func SecureRecord(secrets []Column, recordVal reflect.Value) {
	for _, c := range secrets {
		field := recordVal.FieldByName(c.Name)
		if c.Component.Tag.Is("<password>") {
			field.Set(reflect.ValueOf(PASSWORD_PLACEHOLDER))
		} else {
			field.SetZero()
		}
	}
}
