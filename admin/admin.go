package admin

import (
	"time"

	"github.com/huoyijie/goal/auth"
)

type OperationLog struct {
	ID     uint `goal:"hidden" gorm:"primaryKey"`
	UserID uint `goal:"hidden"`
	User   auth.User `goal:"preload=Username"`
	Date   time.Time
	IPAddr,
	Group,
	Item,
	Action string
	ObjectID uint
}
