package handlers

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/web"
	"gorm.io/gorm"
)

func crud(c *gin.Context, op byte, db *gorm.DB, enforcer *casbin.Enforcer) {
	mt, _ := c.Get("modelType")
	modelType := mt.(reflect.Type)

	record := reflect.New(modelType).Interface()
	if err := c.BindJSON(record); err != nil {
		return
	}

	var tx *gorm.DB
	switch op {
	case 1:
		tx = db.Create(record)
	case 2:
		tx = db.Save(record)
	case 3:
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

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": record,
	})
}

func CrudGet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		model, _ := c.Get("model")
		mt, _ := c.Get("modelType")
		modelType := mt.(reflect.Type)

		hiddens, preloads, columns := web.Reflect(modelType)

		records := reflect.New(reflect.SliceOf(modelType)).Interface()
		tx := db.Model(model)
		for _, column := range preloads {
			tx = tx.Preload(column.Name)
		}
		if err := tx.Find(records).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		recordsVal := reflect.ValueOf(records).Elem()
		for _, c := range hiddens {
			for i := 0; i < recordsVal.Len(); i++ {
				recordVal := recordsVal.Index(i)
				field := recordVal.FieldByName(c.Name)
				field.SetZero()
			}
		}

		for _, c := range preloads {
			for i := 0; i < recordsVal.Len(); i++ {
				recordVal := recordsVal.Index(i)
				preloadVal := recordVal.FieldByName(c.Name)
				// todo hardcode by `ID`
				pk := preloadVal.FieldByName("ID")
				pkVal := pk.Interface()
				preloadField := preloadVal.FieldByName(c.PreloadField)
				dstVal := preloadField.Interface()
				preloadVal.SetZero()
				pk.Set(reflect.ValueOf(pkVal))
				preloadField.Set(reflect.ValueOf(dstVal))
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"columns": columns,
				"records": records,
			},
		})
	}
}

func CrudPost(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, 1, db, enforcer)
	}
}

func CrudPut(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, 2, db, enforcer)
	}
}

func CrudDelete(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		crud(c, 3, db, enforcer)
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