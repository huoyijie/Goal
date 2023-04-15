// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.19
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

func (*Role) Purge() {}
