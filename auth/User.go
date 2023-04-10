// auto generated by https://github.com/huoyijie/GoalGenerator/releases/tag/v0.0.14
// Please do not change anything in this file.
package auth

import (
    "github.com/huoyijie/GoalGenerator"
)

type User struct { 
    goalgenerator.Base

    Username string `gorm:"unique" binding:"required,alphanum,min=3,max=40" goal:"<text>unique,sortable,globalSearch"`
    Email string `gorm:"unique" binding:"required,email" goal:"<text>unique,sortable,globalSearch"`
    Password string `binding:"required,min=8" goal:"<password>secret,hidden"`
    IsSuperuser bool `goal:"<switch>readonly"`
    IsActive bool `goal:"<switch>"`
}
