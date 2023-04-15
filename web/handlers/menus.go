package handlers

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/huoyijie/Goal/web"
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
						menus[i].Items = append(menus[i].Items, web.Menu{Name: web.Item(m), Icon: web.Icon(m)})
						break
					}
				}
				if !found {
					menus = append(menus, web.Menu{Name: group, Items: []web.Menu{{Name: web.Item(m), Icon: web.Icon(m)}}})
				}
			}
		}
		c.JSON(http.StatusOK, web.Result{Data: menus})
	}
}
