// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.12
// Please do not change anything in this file.
package auth

import (
    "time"
)

type Session struct { 

    ID uint `gorm:"primaryKey" goal:"<number>primary,hidden,uint"`
    Key string `gorm:"unique" binding:"required,alphanum,len=32" goal:"<uuid>unique,globalSearch,readonly"`
    UserID uint `goal:"<number>autowired,uint"`
    User User `binding:"required" goal:"<dropdown>belongTo=auth.User.Username,globalSearch,postonly,filter"`
    ExpireDate time.Time `gorm:"index" binding:"required" goal:"<calendar>sortable,desc,showTime,showIcon"`
}