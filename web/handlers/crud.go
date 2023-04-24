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

func purge(record any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) (tx *gorm.DB) {
		tx = db
		if web.IsPurge(record) {
			tx = tx.Unscoped()
		}
		return
	}
}

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
		tx = db.Scopes(purge(record)).Delete(record)
	}

	if err := tx.Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	web.RecordOpLog(db, c, record, action)

	c.JSON(http.StatusOK, web.Result{Data: record})
}

func join(joins []web.Column) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) (tx *gorm.DB) {
		tx = db
		for _, column := range joins {
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
				kv := filter.([]any)
				k := kv[0].(string)
				v := kv[1].(map[string]any)

				if k == "global" {
					_, _, columns := web.Reflect(modelType)
					for _, c := range columns {
						base := c.Component.Tag.(tag.IBase).Get()
						if base.GlobalSearch {
							filterField := c.Name
							if d, ok := c.Component.Tag.(*tag.Dropdown); ok {
								if d.BelongTo != nil {
									filterField = fmt.Sprintf("%s.%s", d.BelongTo.Name, d.BelongTo.Field)
								} else if d.HasOne != nil {
									filterField = fmt.Sprintf("%s.%s", d.HasOne.Name, d.HasOne.Field)
								}
							}

							if hasPrev {
								sb.WriteString(" and ")
							}

							// todo: global search support only one column that is text
							value := v["value"].(string)
							table, field := web.TableFieldName(db, model, modelType, filterField)
							matchMode := v["matchMode"].(string)
							web.FilterClause(&sb, table, field, matchMode, value)
							hasPrev = true
							continue loop
						}
					}
				}

				table, field := web.TableFieldName(db, model, modelType, k)

				if value, ok := v["value"]; ok {
					switch value := value.(type) {
					case bool:
						if hasPrev {
							sb.WriteString(" and ")
						}
						var val int
						if value {
							val = 1
						}
						matchMode := v["matchMode"].(string)
						web.FilterClause(&sb, table, field, matchMode, fmt.Sprintf("%d", val))
						hasPrev = true
					case string:
						if hasPrev {
							sb.WriteString(" and ")
						}
						matchMode := v["matchMode"].(string)
						web.FilterClause(&sb, table, field, matchMode, value)
						hasPrev = true
					case int, uint:
						if hasPrev {
							sb.WriteString(" and ")
						}
						matchMode := v["matchMode"].(string)
						web.FilterClause(&sb, table, field, matchMode, fmt.Sprintf("%d", value))
						hasPrev = true
					case float64:
						// todo float64 eqauls maybe wrong sometimes
						if hasPrev {
							sb.WriteString(" and ")
						}
						matchMode := v["matchMode"].(string)
						web.FilterClause(&sb, table, field, matchMode, strconv.FormatFloat(value, 'f', -1, 64))
						hasPrev = true
					}
				} else {
					operator := v["operator"].(string)
					constraints := v["constraints"].([]any)
					if operator != "" && len(constraints) > 0 {
						if hasPrev {
							sb.WriteString(" and ")
						}
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
							web.FilterClause(&sb, table, field, matchMode, value)
							hasOp = true
						}
						sb.WriteString(" )")
						hasPrev = true
					}
				}
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

	secrets, joins, _ := web.Reflect(modelType)

	var total int64
	if err := db.Scopes(filters(model, modelType, session, mine, lazyParam), join(joins)).Count(&total).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	records := reflect.New(reflect.SliceOf(modelType)).Interface()

	if err := db.Scopes(
		filters(model, modelType, session, mine, lazyParam),
		join(joins),
		pagiSort(model, modelType, lazyParam),
	).Find(records).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	recordsVal := reflect.ValueOf(records).Elem()
	web.SecureRecords(secrets, recordsVal)

	for _, c := range joins {
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
			Ctrl:    web.IsCtrl(model),
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

		if err := db.Scopes(purge(model)).Delete(model, ids).Error; err != nil {
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
		err := db.Unscoped().Select("ID").Where(record).First(dbRecord).Error
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

func Select(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		model, _ := c.Get("model")
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)

		field, _ := c.Params.Get("field")
		f, _ := modelType.FieldByName(field)
		t := web.GetComponent(f)
		d := t.(*tag.Dropdown)
		dt := reflect.TypeOf(d).Elem()
		dVal := reflect.ValueOf(d).Elem()

		var m string
		for i := 0; i < dt.NumField(); i++ {
			f := dt.Field(i)
			if fVal := dVal.FieldByName(f.Name); fVal.Kind() == reflect.Bool && fVal.Bool() {
				m = fmt.Sprintf("%s%s", field, f.Name)
				break
			}
		}

		getOptions := reflect.ValueOf(model).MethodByName(m)
		if !getOptions.IsValid() {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var data []web.Option
		out := getOptions.Call([]reflect.Value{})
		for i := 0; i < out[0].Len(); i++ {
			o := out[0].Index(i).Interface()
			data = append(data, web.Option{
				Label: fmt.Sprintf("%v", o),
				Value: o,
			})
		}
		c.JSON(http.StatusOK, web.Result{Data: data})
	}
}
