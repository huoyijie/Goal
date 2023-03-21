package handlers

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/goal/web"
)

func Menus(models []any, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := web.GetSession(c)
		menus := []web.Menu{}
		for _, m := range models {
			if web.AllowAny(session, web.Obj(m), enforcer) {
				var found bool
				group := web.Group(m)
				for i := range menus {
					if menus[i].Name == group {
						found = true
						menus[i].Items = append(menus[i].Items, web.Menu{Name: web.Item(m), Items: nil})
						break
					}
				}
				if !found {
					menus = append(menus, web.Menu{Name: group, Items: []web.Menu{{Name: web.Item(m), Items: nil}}})
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": menus,
		})
	}
}
