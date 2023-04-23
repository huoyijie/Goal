// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.36
// Please do not change anything in this file.
package class

import (
	"github.com/huoyijie/GoalGenerator/model"
)

type Teacher struct {
	model.Base
	Name string `binding:"required" goal:"<text>globalSearch,filter"`
}

func (*Teacher) Icon() string {
	return "teacher"
}

func (*Teacher) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Class"
	t["zh-CN"] = "班级"
	return t
}

func (*Teacher) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Teacher | teachers"
	t["zh-CN"] = "老师"
	return t
}

func (*Teacher) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":   "ID",
			"Name": "Name",
		},
		"zh-CN": {
			"ID":   "ID",
			"Name": "姓名",
		},
	}
}

func (m *Teacher) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
