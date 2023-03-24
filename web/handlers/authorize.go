package handlers

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/admin"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/web"
	"gorm.io/gorm"
)

func GetPerms(models []any, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := web.PermsParam{}
		if err := c.BindUri(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		perms := []web.Perm{}
		for _, m := range models {
		inner:
			for _, act := range web.Actions() {
				switch m.(type) {
				case *admin.OperationLog:
					if act != "get" {
						continue inner
					}
				}
				perms = append(perms, web.NewPerm(web.Obj(m), act))
			}
		}

		role := auth.Role{ID: param.RoleID}
		permissions := enforcer.GetPermissionsForUser(role.RoleID())
		rolePerms := []web.Perm{}
		for _, p := range permissions {
			rolePerms = append(rolePerms, web.NewPerm(p[1], p[2]))
		}

		availablePerms := []web.Perm{}
	outer:
		for _, p1 := range perms {
			for _, p2 := range rolePerms {
				if p2.Code == p1.Code {
					continue outer
				}
			}
			availablePerms = append(availablePerms, p1)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": [][]web.Perm{availablePerms, rolePerms},
		})
	}
}

func ChangePerms(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := web.PermsParam{}
		if err := c.BindUri(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		role := auth.Role{ID: param.RoleID}
		var selected []web.Perm
		if err := c.BindJSON(&selected); err != nil {
			return
		}
		var permissions [][]string
		for _, perm := range selected {
			permissions = append(permissions, perm.Val())
		}
		enforcer.DeletePermissionsForUser(role.RoleID())
		enforcer.AddPermissionsForUser(role.RoleID(), permissions...)
		c.JSON(http.StatusOK, gin.H{"code": 0})
	}
}

func GetRoles(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := web.RolesParam{}
		if err := c.BindUri(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var roles []auth.Role
		if err := db.Find(&roles).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user := auth.User{ID: param.UserID}
		userRoles, err := enforcer.GetRolesForUser(user.Sub())
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		selected := []uint{}
		for _, p := range userRoles {
			selected = append(selected, web.ParseRoleID(p))
		}
		selectedRoles := []auth.Role{}
		if len(selected) > 0 {
			if err := db.Find(&selectedRoles, selected).Error; err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		availableRoles := []auth.Role{}
	outer:
		for _, p1 := range roles {
			for _, p2 := range selectedRoles {
				if p2.ID == p1.ID {
					continue outer
				}
			}
			availableRoles = append(availableRoles, p1)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": [][]auth.Role{availableRoles, selectedRoles},
		})
	}
}

func ChangeRoles(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := web.RolesParam{}
		if err := c.BindUri(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		user := auth.User{ID: param.UserID}
		var selected []auth.Role
		if err := c.BindJSON(&selected); err != nil {
			return
		}
		var roles []string
		for _, role := range selected {
			roles = append(roles, role.RoleID())
		}
		enforcer.DeleteRolesForUser(user.Sub())
		enforcer.AddRolesForUser(user.Sub(), roles)
		c.JSON(http.StatusOK, gin.H{"code": 0})
	}
}
