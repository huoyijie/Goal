// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.36
// Please do not change anything in this file.
package paper

import (
	"github.com/huoyijie/GoalGenerator/model"
)

type Choice struct {
	model.Base
	Content    string `binding:"required" goal:"<text>globalSearch,filter"`
	QuestionID uint   `goal:"<number>autowired,uint"`
}

func (*Choice) Icon() string {
	return "choice"
}

func (*Choice) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Paper"
	t["zh-CN"] = "试卷"
	return t
}

func (*Choice) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Choice | choices"
	t["zh-CN"] = "选项"
	return t
}

func (*Choice) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":         "ID",
			"Content":    "Content",
			"QuestionID": "Question ID",
		},
		"zh-CN": {
			"ID":         "ID",
			"Content":    "内容",
			"QuestionID": "问题ID",
		},
	}
}

func (m *Choice) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	return t
}
