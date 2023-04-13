// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.17
// Please do not change anything in this file.
package auth

import (
    "github.com/huoyijie/GoalGenerator/model"
    "time"
)

type Session struct { 
    model.Base

    Key string `gorm:"unique" binding:"required,alphanum,len=32" goal:"<uuid>unique,readonly,globalSearch,filter"`
    UserID uint `goal:"<number>autowired,uint"`
    User User `binding:"required" goal:"<dropdown>postonly,filter,belongTo=auth.User.Username"`
    ExpireDate time.Time `gorm:"index" binding:"required" goal:"<calendar>sortable,desc,filter,showTime,showIcon"`
}


func (*Session) Lazy() {}
var _ model.Lazy = (*Session)(nil)
