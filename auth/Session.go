// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.14
// Please do not change anything in this file.
package auth

import (
    "github.com/huoyijie/GoalGenerator"
    "time"
)

type Session struct { 
    goalgenerator.Base

    Key string `gorm:"unique" binding:"required,alphanum,len=32" goal:"<uuid>unique,globalSearch,readonly"`
    UserID uint `goal:"<number>autowired,uint"`
    User User `binding:"required" goal:"<dropdown>belongTo=auth.User.Username,globalSearch,postonly,filter"`
    ExpireDate time.Time `gorm:"index" binding:"required" goal:"<calendar>sortable,showTime,showIcon"`
}
