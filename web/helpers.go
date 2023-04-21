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
	"github.com/huoyijie/Goal/admin"
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/Goal/util"
	"github.com/huoyijie/Goal/web/tag"
	"github.com/huoyijie/GoalGenerator/model"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		if field.Name == "Base" {
			s, p, c := Reflect(field.Type)
			secrets = append(secrets, s...)
			preloads = append(preloads, p...)
			columns = append(columns, c...)
			continue
		}

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

		if d, ok := component.(*tag.Dropdown); ok && (d.BelongTo != nil || d.HasOne != nil) {
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

func SetCookieSessionid(c *gin.Context, sessionid string, rememberMe bool, domain string, secure bool) {
	// keep g_sessionid until the browser closed
	maxAge := 0

	if len(sessionid) == 0 {
		// sign out: delete g_sessionid right now
		maxAge = -1
	} else if rememberMe {
		// sign in: remember me was checked
		maxAge = 3 * 24 * 60 * 60
	}
	c.SetCookie("g_sessionid", sessionid, maxAge, "/", domain, secure, true)
}

func ClearSessions(db *gorm.DB) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		util.Log(db.Delete(&auth.Session{}, "expire_date < ?", time.Now()).Error)
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

func Icon(model any) string {
	m := reflect.ValueOf(model).MethodByName("Icon")
	out := m.Call([]reflect.Value{})
	return out[0].String()
}

func IsLazy(m any) bool {
	return reflect.TypeOf(m).Implements(reflect.TypeOf((*model.Lazy)(nil)).Elem())
}

func IsCtrl(m any) bool {
	return reflect.TypeOf(m).Implements(reflect.TypeOf((*model.Ctrl)(nil)).Elem())
}

func IsPurge(m any) bool {
	return reflect.TypeOf(m).Implements(reflect.TypeOf((*model.Purge)(nil)).Elem())
}

func IsTabler(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*schema.Tabler)(nil)).Elem())
}

func TableName(db *gorm.DB, model any, modelType reflect.Type) (table string) {
	if IsTabler(modelType) {
		m := reflect.ValueOf(model).Elem().MethodByName("TableName")
		table = m.Call([]reflect.Value{})[0].String()
	} else {
		table = db.NamingStrategy.TableName(modelType.Name())
	}
	return
}

func TableFieldName(db *gorm.DB, model any, modelType reflect.Type, filterField string) (table, field string) {
	if tmp := strings.Split(filterField, "."); len(tmp) == 2 {
		table = tmp[0]
		field = db.NamingStrategy.ColumnName("", tmp[1])
	} else {
		table = TableName(db, model, modelType)
		field = db.NamingStrategy.ColumnName("", tmp[0])
	}
	return
}

func Convert(matchMode, value string) (condition string) {
	switch matchMode {
	case "startsWith":
		condition = " LIKE '" + value + "%'"
	case "contains":
		condition = " LIKE '%" + value + "%'"
	case "notContains":
		condition = " NOT LIKE '%" + value + "%'"
	case "endsWith":
		condition = " LIKE '%" + value + "'"
	case "equals":
		condition = " = '" + value + "'"
	case "notEquals":
		condition = " != '" + value + "'"
	case "in":
		panic("not implement")
	case "lt":
		condition = " < " + value
	case "lte":
		condition = " <= " + value
	case "gt":
		condition = " > " + value
	case "gte":
		condition = " >= " + value
	case "between":
		panic("not implement")
	case "dateIs":
		d, _ := time.Parse(time.RFC3339, value)
		d1 := d.AddDate(0, 0, 1).Format(time.RFC3339)
		condition = " BETWEEN '" + value + "' and '" + d1 + "'"
	case "dateIsNot":
		d, _ := time.Parse(time.RFC3339, value)
		d1 := d.AddDate(0, 0, 1).Format(time.RFC3339)
		condition = " NOT BETWEEN '" + value + "' and '" + d1 + "'"
	case "dateBefore":
		condition = " < '" + value + "'"
	case "dateAfter":
		d, _ := time.Parse(time.RFC3339, value)
		d1 := d.AddDate(0, 0, 1).Format(time.RFC3339)
		condition = " > '" + d1 + "'"
	}
	return
}

func FilterClause(sb *strings.Builder, table, field, matchMode, value string) {
	sb.WriteRune('`')
	sb.WriteString(table)
	sb.WriteString("`.`")
	sb.WriteString(field)
	sb.WriteRune('`')
	sb.WriteString(Convert(matchMode, value))
}
