// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.35
// Please do not change anything in this file.
package auth

import (
	"github.com/huoyijie/GoalGenerator/model"
	"time"
)

type Session struct {
	model.Base
	Key        string    `gorm:"unique" binding:"required,alphanum,len=32" goal:"<uuid>unique,readonly,globalSearch,filter"`
	UserID     uint      `goal:"<number>autowired,uint"`
	User       User      `binding:"required" goal:"<dropdown>postonly,filter,belongTo=auth.User.Username"`
	ExpireDate time.Time `gorm:"index" binding:"required" goal:"<calendar>sortable,desc,filter,showTime,showIcon"`
}

func (Session) TableName() string {
	return "auth_sessions"
}

func (*Session) Icon() string {
	return "ticket"
}

func (*Session) Purge() {}

func (*Session) Lazy() {}

func (*Session) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Auth"
	t["zh-CN"] = "认证授权"
	return t
}

func (*Session) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Session | sessions"
	t["zh-CN"] = "会话"
	return t
}

func (*Session) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":         "ID",
			"Key":        "Key",
			"UserID":     "User ID",
			"User":       "User",
			"ExpireDate": "Expire Date",
		},
		"zh-CN": {
			"ID":         "ID",
			"Key":        "会话",
			"UserID":     "用户ID",
			"User":       "用户",
			"ExpireDate": "过期时间",
		},
	}
}

func (m *Session) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
