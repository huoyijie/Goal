// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.37
// Please do not change anything in this file.
package country

import (
	"github.com/huoyijie/GoalGenerator/model"
)

type Identify struct {
	model.Base
	NO       string `gorm:"unique" binding:"required,alphanum,len=18" goal:"<text>unique,globalSearch,filter"`
	PeopleID uint   `gorm:"unique" goal:"<number>unique,autowired,uint"`
}

func (*Identify) Icon() string {
	return "id-card"
}

func (*Identify) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Country"
	t["zh-CN"] = "国家"
	return t
}

func (*Identify) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Identify | identifies"
	t["zh-CN"] = "身份证"
	return t
}

func (*Identify) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":       "ID",
			"NO":       "NO.",
			"PeopleID": "People ID",
		},
		"zh-CN": {
			"ID":       "ID",
			"NO":       "号码",
			"PeopleID": "公民ID",
		},
	}
}

func (m *Identify) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
