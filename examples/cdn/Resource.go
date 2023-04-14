// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.17
// Please do not change anything in this file.
package cdn

import (
    "github.com/huoyijie/GoalGenerator/model"
)

type Resource struct { 
    model.Base

    File string `gorm:"unique" binding:"required" goal:"<file>unique,postonly,globalSearch,filter,uploadTo=uploads"`
    Status string `binding:"required" goal:"<dropdown>filter,strings"`
    Level uint `binding:"required" goal:"<dropdown>filter,uints"`
}

func (*Resource) StatusStrings() []string {
    return []string{"tbd", "on", "off"}
}
func (*Resource) LevelUints() []uint {
    return []uint{1, 2, 3}
}

func (*Resource) Lazy() {}
var _ model.Lazy = (*Resource)(nil)
