// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.36
// Please do not change anything in this file.
package auth

import (
	"github.com/huoyijie/GoalGenerator/model"
)

type Role struct {
	model.Base
	Name string `gorm:"unique" binding:"required,alphanum,min=3,max=40" goal:"<text>unique,sortable,globalSearch,filter"`
}

func (Role) TableName() string {
	return "auth_roles"
}

func (*Role) Icon() string {
	return "key"
}

func (*Role) Purge() {}

func (*Role) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Auth"
	t["zh-CN"] = "认证授权"
	return t
}

func (*Role) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Role | roles"
	t["zh-CN"] = "角色"
	return t
}

func (*Role) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":   "ID",
			"Name": "Name",
		},
		"zh-CN": {
			"ID":   "ID",
			"Name": "名称",
		},
	}
}

func (m *Role) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
