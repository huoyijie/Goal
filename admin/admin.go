package admin

import (
	"time"

	"github.com/huoyijie/goal/auth"
)

type OperationLog struct {
	ID          uint      `goal:"<number>primary,hidden" gorm:"primaryKey"`
	UserID      uint      `goal:"<number>hidden" gorm:"index"`
	User        auth.User `goal:"<dropdown>belongTo=auth.User.Username,sortable"`
	Date        time.Time `goal:"<calendar>showTime,sortable,desc" gorm:"index"`
	IP          string    `goal:"<text>"`
	Group, Item string    `goal:"<text>hidden"`
	Action      string    `goal:"<text>"`
	ObjectID    uint      `goal:"<number>"`
}
