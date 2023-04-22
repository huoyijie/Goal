package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huoyijie/Goal/web"
	"github.com/huoyijie/GoalGenerator/model"
)

var langs = []string{"en", "zh-CN"}

func Translate(models []any) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := gin.H{}
		for _, m := range models {
			translate := m.(model.Translate)
			tPkg := translate.TranslatePkg()
			tName := translate.TranslateName()
			tFields := translate.TranslateFields()
			tOptions := translate.TranslateOptions()

			pkg := web.Group(m)
			item := web.Item(m)
			for _, lang := range langs {
				if _, found := t[lang]; !found {
					t[lang] = gin.H{}
				}
				if _, found := t[lang].(gin.H)[pkg]; !found {
					t[lang].(gin.H)[pkg] = gin.H{
						"label": tPkg[lang],
					}
				}
				t[lang].(gin.H)[pkg].(gin.H)[item] = gin.H{
					"label":   tName[lang],
					"fields":  tFields[lang],
					"options": tOptions[lang],
				}
			}
		}
		c.JSON(http.StatusOK, web.Result{Data: t})
	}
}
