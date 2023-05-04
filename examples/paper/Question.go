// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.37
// Please do not change anything in this file.
package paper

import (
	"github.com/huoyijie/GoalGenerator/model"
)

type Question struct {
	model.Base
	Label   string   `gorm:"unique" binding:"required" goal:"<text>unique,globalSearch,filter"`
	Choices []Choice `binding:"required" goal:"<inline>hasMany=paper.Choice"`
}

func (*Question) Icon() string {
	return "paper"
}

func (*Question) Purge() {}

func (*Question) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Paper"
	t["zh-CN"] = "试卷"
	return t
}

func (*Question) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Question | questions"
	t["zh-CN"] = "问题"
	return t
}

func (*Question) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":      "ID",
			"Label":   "Label",
			"Choices": "Choices",
		},
		"zh-CN": {
			"ID":      "ID",
			"Label":   "题干",
			"Choices": "选项",
		},
	}
}

func (m *Question) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
