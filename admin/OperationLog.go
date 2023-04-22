// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.31
// Please do not change anything in this file.
package admin

import (
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/GoalGenerator/model"
	"time"
)

type OperationLog struct {
	model.Base
	UserID   uint      `gorm:"index" goal:"<number>hidden,uint"`
	User     auth.User `goal:"<dropdown>sortable,globalSearch,filter,belongTo=auth.User.Username"`
	Date     time.Time `gorm:"index" goal:"<calendar>sortable,desc,filter,showTime"`
	IP       string    `goal:"<text>filter"`
	Group    string    `goal:"<dropdown>filter,dynamicStrings"`
	Item     string    `goal:"<dropdown>filter,dynamicStrings"`
	Action   string    `goal:"<dropdown>filter,strings"`
	ObjectID uint      `goal:"<number>uint"`
}

func (OperationLog) TableName() string {
	return "admin_operation_logs"
}

func (*OperationLog) Icon() string {
	return "save"
}

func (*OperationLog) Lazy() {}

func (*OperationLog) Ctrl() {}

func (*OperationLog) TranslatePkg() map[string]string {
	t := map[string]string{}
	t["en"] = "Admin"
	t["zh-CN"] = "通用管理"
	return t
}

func (*OperationLog) TranslateName() map[string]string {
	t := map[string]string{}
	t["en"] = "Operation Log | operation logs"
	t["zh-CN"] = "操作日志"
	return t
}

func (*OperationLog) TranslateFields() map[string]map[string]string {
	return map[string]map[string]string{
		"en": {
			"ID":       "ID",
			"UserID":   "User ID",
			"User":     "User",
			"Date":     "Date",
			"IP":       "IP",
			"Group":    "Group",
			"Item":     "Item",
			"Action":   "Action",
			"ObjectID": "Object ID",
		},
		"zh-CN": {
			"ID":       "ID",
			"UserID":   "用户ID",
			"User":     "用户",
			"Date":     "时间",
			"IP":       "IP",
			"Group":    "组",
			"Item":     "项",
			"Action":   "动作",
			"ObjectID": "目标ID",
		},
	}
}

// Please implements this method in another file
// func (*OperationLog) GroupDynamicStrings() []string {
//     return []string{"option1", "option2"}
// }

// Please implements this method in another file
// func (*OperationLog) TranslateGroupDynamicStrings() map[string]map[string]string {
//     return map[string]map[string]string{"en": {"option1": "option 1", "option2": "option 2"}, "zh-CN": {"option1": "选项1", "option2": "选项2"},}
// }

// Please implements this method in another file
// func (*OperationLog) ItemDynamicStrings() []string {
//     return []string{"option1", "option2"}
// }

// Please implements this method in another file
// func (*OperationLog) TranslateItemDynamicStrings() map[string]map[string]string {
//     return map[string]map[string]string{"en": {"option1": "option 1", "option2": "option 2"}, "zh-CN": {"option1": "选项1", "option2": "选项2"},}
// }

func (*OperationLog) ActionStrings() []string {
	return []string{"post", "put", "delete"}
}

func (*OperationLog) TranslateActionStrings() map[string]map[string]string {
	return map[string]map[string]string{"en": {"post": "Add", "put": "Change", "delete": "Delete"}, "zh-CN": {"post": "新增", "put": "修改", "delete": "删除"}}
}

func (m *OperationLog) TranslateOptions() map[string]map[string]map[string]string {
	t := map[string]map[string]map[string]string{"en": {}, "zh-CN": {}}

	tGroup := m.TranslateGroupDynamicStrings()
	t["en"]["Group"] = tGroup["en"]
	t["zh-CN"]["Group"] = tGroup["zh-CN"]

	tItem := m.TranslateItemDynamicStrings()
	t["en"]["Item"] = tItem["en"]
	t["zh-CN"]["Item"] = tItem["zh-CN"]

	tAction := m.TranslateActionStrings()
	t["en"]["Action"] = tAction["en"]
	t["zh-CN"]["Action"] = tAction["zh-CN"]

	return t
}
