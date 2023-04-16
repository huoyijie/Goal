package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huoyijie/Goal/web"
)

func Translate(c *gin.Context) {
	/**
		{
	    "cdn": {
	        "label": "CDN",
	        "resource": {
	            "label": "Resource | resources",
	            "ID": "ID",
	            "File": "File",
	            "Status": "Status",
	            "Level": "Level",
	            "options": {
	                "Status": {
	                    "tbd": "tbd",
	                    "on": "online",
	                    "off": "offline"
	                },
	                "Level": {
	                    "1": "no.1",
	                    "2": "no.2",
	                    "3": "no.3"
	                }
	            }
	        }
	    }
		}
	*/
	/**
		{
	    "cdn": {
	        "label": "CDN",
	        "resource": {
	            "label": "资源",
	            "ID": "ID",
	            "File": "文件",
	            "Status": "状态",
	            "Level": "等级",
	            "options": {
	                "Status": {
	                    "tbd": "待审核",
	                    "on": "上线",
	                    "off": "下线"
	                },
	                "Level": {
	                    "1": "级别1",
	                    "2": "级别2",
	                    "3": "级别3"
	                }
	            }
	        }
	    }
		}
	*/
	c.JSON(http.StatusOK, web.Result{Data: gin.H{
		"en": gin.H{
			"cdn": gin.H{
				"label": "CDN",
				"resource": gin.H{
					"label":  "Resource | resources",
					"ID":     "ID",
					"File":   "File",
					"Status": "Status",
					"Level":  "Level",
					"options": gin.H{
						"Status": gin.H{
							"tbd": "tbd",
							"on":  "online",
							"off": "offline",
						},
						"Level": gin.H{
							"1": "no.1",
							"2": "no.2",
							"3": "no.3",
						},
					},
				},
			},
		},
		"zh_CN": gin.H{
			"cdn": gin.H{
				"label": "CDN",
				"resource": gin.H{
					"label":  "资源",
					"ID":     "ID",
					"File":   "文件",
					"Status": "状态",
					"Level":  "等级",
					"options": gin.H{
						"Status": gin.H{
							"tbd": "待审核",
							"on":  "上线",
							"off": "下线",
						},
						"Level": gin.H{
							"1": "级别1",
							"2": "级别2",
							"3": "级别3",
						},
					},
				},
			},
		},
	}})
}
