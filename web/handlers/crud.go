package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/Goal/util"
	"github.com/huoyijie/Goal/web"
	"github.com/huoyijie/Goal/web/tag"
	"gorm.io/gorm"
)

func crud(c *gin.Context, action string, db *gorm.DB, enforcer *casbin.Enforcer) {
	mt, _ := c.Get("modelType")
	modelType := mt.(reflect.Type)

	record := reflect.New(modelType).Interface()
	if err := c.BindJSON(record); err != nil {
		return
	}

	var tx *gorm.DB
	switch action {
	case "post", "put":
		web.AutowiredCreator(c, action, record)
		switch r := record.(type) {
		case *auth.User:
			if action == "put" && r.Password == web.PASSWORD_PLACEHOLDER {
				o := &auth.User{}
				o.ID = r.ID
				if err := db.First(o).Error; err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				r.Password = o.Password
			} else {
				r.Password = util.BcryptHash(r.Password)
			}
		}
		tx = db.Save(record)
	case "delete":
		switch r := record.(type) {
		case *auth.Role:
			enforcer.DeletePermissionsForUser(r.RoleID())
		case *auth.User:
			enforcer.DeleteRolesForUser(r.Sub())
		}
		tx = db.Delete(record)
	}

	if err := tx.Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	web.RecordOpLog(db, c, record, action)

	c.JSON(http.StatusOK, web.Result{Data: record})
}

func preload(preloads []web.Column) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) (tx *gorm.DB) {
		tx = db
		for _, column := range preloads {
			tx = tx.Joins(column.Name)
		}
		return
	}
}

func pagiSort(model any, modelType reflect.Type, lazyParam *web.LazyParam) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) (tx *gorm.DB) {
		tx = db
		if lazyParam != nil {
			if lazyParam.Offset > 0 {
				tx = tx.Offset(lazyParam.Offset)
			}
			if lazyParam.Limit > 0 {
				tx = tx.Limit(lazyParam.Limit)
			}
			if lazyParam.SortField != "" {
				var table, field string
				if tmp := strings.Split(lazyParam.SortField, "."); len(tmp) == 2 {
					table = tmp[0]
					field = db.NamingStrategy.ColumnName("", tmp[1])
				} else {
					table = web.TableName(db, model, modelType)
					field = db.NamingStrategy.ColumnName("", tmp[0])
				}
				sortOrder := "asc"
				if lazyParam.SortOrder == -1 {
					sortOrder = "desc"
				}
				orderBy := fmt.Sprintf("`%s`.`%s` %s", table, field, sortOrder)
				tx = tx.Order(orderBy)
			}
		}
		return
	}
}

func filters(model any, modelType reflect.Type, session *auth.Session, mine bool, lazyParam *web.LazyParam) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) (tx *gorm.DB) {
		tx = db.Model(model)
		sb := strings.Builder{}

		var hasPrev bool
		if mine && reflect.ValueOf(model).Elem().FieldByName("Creator").IsValid() {
			sb.WriteRune('`')
			sb.WriteString(web.TableName(db, model, modelType))
			sb.WriteString("`.`creator` = ")
			sb.WriteString(strconv.FormatUint(uint64(session.UserID), 10))
			hasPrev = true
		}
		if lazyParam != nil && lazyParam.Filters != "" {
			var filters []any
			json.Unmarshal([]byte(lazyParam.Filters), &filters)
		loop:
			for _, filter := range filters {
				if hasPrev {
					sb.WriteString(" and ")
				}
				kv := filter.([]any)
				k := kv[0].(string)
				if k == "global" {
					continue loop
				}

				var table, field string
				if tmp := strings.Split(k, "."); len(tmp) == 2 {
					table = tmp[0]
					field = db.NamingStrategy.ColumnName("", tmp[1])
				} else {
					table = web.TableName(db, model, modelType)
					field = db.NamingStrategy.ColumnName("", tmp[0])
				}

				v := kv[1].(map[string]any)
				if value, ok := v["value"]; ok {
					switch value := value.(type) {
					case bool:
						var val int
						if value {
							val = 1
						}
						matchMode := v["matchMode"].(string)
						sb.WriteRune('`')
						sb.WriteString(table)
						sb.WriteString("`.`")
						sb.WriteString(field)
						sb.WriteRune('`')
						sb.WriteString(web.Convert(matchMode, fmt.Sprintf("%d", val)))
					}
				} else {
					operator := v["operator"].(string)
					constraints := v["constraints"].([]any)
					sb.WriteString("( ")
					var hasOp bool
					for _, constraint := range constraints {
						if hasOp {
							sb.WriteRune(' ')
							sb.WriteString(operator)
							sb.WriteRune(' ')
						}
						c := constraint.(map[string]any)
						value := fmt.Sprintf("%v", c["value"])
						matchMode := c["matchMode"].(string)
						sb.WriteRune('`')
						sb.WriteString(table)
						sb.WriteString("`.`")
						sb.WriteString(field)
						sb.WriteRune('`')
						sb.WriteString(web.Convert(matchMode, value))
						hasOp = true
					}
					sb.WriteString(" )")
				}
				hasPrev = true
			}
		}

		if sb.Len() > 0 {
			tx = tx.Where(sb.String())
		}
		return
	}
}

func crudGet(c *gin.Context, db *gorm.DB, mine bool) {
	model, _ := c.Get("model")
	mt, _ := c.Get("modelType")
	modelType := mt.(reflect.Type)
	session := web.GetSession(c)

	var lazyParam *web.LazyParam
	if web.IsLazy(model) {
		lazyParam = &web.LazyParam{}
		if err := c.Bind(lazyParam); err != nil {
			return
		}
	}

	secrets, preloads, _ := web.Reflect(modelType)

	var total int64
	if err := db.Scopes(filters(model, modelType, session, mine, lazyParam), preload(preloads)).Count(&total).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	records := reflect.New(reflect.SliceOf(modelType)).Interface()

	if err := db.Scopes(
		filters(model, modelType, session, mine, lazyParam),
		preload(preloads),
		pagiSort(model, modelType, lazyParam),
	).Find(records).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	recordsVal := reflect.ValueOf(records).Elem()
	web.SecureRecords(secrets, recordsVal)

	for _, c := range preloads {
		for i := 0; i < recordsVal.Len(); i++ {
			recordVal := recordsVal.Index(i)
			preloadVal := recordVal.FieldByName(c.Name)
			secrets, _, _ := web.Reflect(reflect.TypeOf(preloadVal.Interface()))
			web.SecureRecord(secrets, preloadVal)
		}
	}

	c.JSON(http.StatusOK, web.Result{Data: &web.RecordList{Total: uint(total), List: records}})
}

func CrudDataTable(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)
		_, _, columns := web.Reflect(modelType)

		session := web.GetSession(c)
		model, _ := c.Get("model")

		perms := &web.Perms{}
		permsVal := reflect.ValueOf(perms).Elem()
		for _, act := range web.Actions() {
			if web.Allow(session, web.Obj(model), act, enforcer) {
				f := permsVal.FieldByName(util.ToUpperFirstLetter(act))
				f.SetBool(true)
			}
		}

		c.JSON(http.StatusOK, web.Result{Data: web.DataTable{
			Lazy:    web.IsLazy(model),
			Columns: columns,
			Perms:   perms,
		}})
	}
}

func CrudGetMine(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		crudGet(c, db, true)
	}
}

func CrudGet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		crudGet(c, db, false)
	}
}

func CrudPost(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, "post", db, enforcer)
	}
}

func CrudPut(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, "put", db, enforcer)
	}
}

func CrudDelete(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, "delete", db, enforcer)
	}
}

func CrudBatchDelete(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		model, _ := c.Get("model")

		ids := []uint{}
		if err := c.BindJSON(&ids); err != nil {
			return
		}

		switch model.(type) {
		case *auth.Role:
			for _, id := range ids {
				role := auth.Role{}
				role.ID = id
				enforcer.DeletePermissionsForUser(role.RoleID())
				enforcer.DeleteRole(role.RoleID())
			}
		case *auth.User:
			for _, id := range ids {
				user := auth.User{}
				user.ID = id
				enforcer.DeleteRolesForUser(user.Sub())
				enforcer.DeleteUser(user.Sub())
			}
		}

		if err := db.Delete(model, ids).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		web.RecordOpLogs(db, c, ids, "delete")

		c.JSON(http.StatusOK, web.Result{})
	}
}

func CrudExist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)

		record := reflect.New(modelType).Interface()
		c.ShouldBindJSON(record)

		dbRecord := reflect.New(modelType).Interface()
		// todo refactor hardcode by `ID`
		err := db.Select("ID").Where(record).First(dbRecord).Error
		if err == nil {
			c.JSON(http.StatusOK, web.Result{Data: dbRecord})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, web.Result{})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func Upload(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)

		field, _ := c.Params.Get("field")
		f, _ := modelType.FieldByName(field)

		t := web.GetComponent(f)
		uploadPath := t.(*tag.File).UploadTo.Path

		file, _ := c.FormFile("file")

		ext := filepath.Ext(file.Filename)
		name, _ := strings.CutSuffix(file.Filename, ext)
		fname := fmt.Sprintf("%s.%s%s", name, util.RandString(8), ext)

		now := time.Now()
		filePath := filepath.Join(uploadPath, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fname)

		c.SaveUploadedFile(file, filePath)

		c.JSON(http.StatusOK, web.Result{Data: filePath})
	}
}
