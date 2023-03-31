package handlers

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"github.com/huoyijie/goal/web"
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
				o := &auth.User{ID: r.ID}
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

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": record,
	})
}

func crudGet(c *gin.Context, db *gorm.DB, mine bool) {
	model, _ := c.Get("model")
	mt, _ := c.Get("modelType")
	modelType := mt.(reflect.Type)
	session := web.GetSession(c)

	secrets, preloads, _ := web.Reflect(modelType)

	records := reflect.New(reflect.SliceOf(modelType)).Interface()
	tx := db.Model(model)
	if mine {
		tx = tx.Where("creator = ?", session.UserID)
	}
	for _, column := range preloads {
		tx = tx.Joins(column.Name)
	}
	if err := tx.Find(records).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": records,
	})
}

func CrudColumns(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)
		_, _, columns := web.Reflect(modelType)

		session := web.GetSession(c)
		model, _ := c.Get("model")
		perms := gin.H{}
		for _, act := range web.Actions() {
			perms[act] = web.Allow(session, web.Obj(model), act, enforcer)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"columns": columns,
				"perms":   perms,
			},
		})
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
				role := auth.Role{ID: id}
				enforcer.DeletePermissionsForUser(role.RoleID())
				enforcer.DeleteRole(role.RoleID())
			}
		case *auth.User:
			for _, id := range ids {
				user := auth.User{ID: id}
				enforcer.DeleteRolesForUser(user.Sub())
				enforcer.DeleteUser(user.Sub())
			}
		}

		if err := db.Delete(model, ids).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		web.RecordOpLogs(db, c, ids, "delete")

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
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
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": dbRecord})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
