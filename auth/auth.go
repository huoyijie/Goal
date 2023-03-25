package auth

import (
	"fmt"
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username string `validate:"required,alphanum,min=3,max=40" binding:"required,alphanum,min=3,max=40" gorm:"unique"`
	Email    string `validate:"required,email" binding:"required,email" gorm:"unique"`
	Password string `validate:"required,min=8" binding:"required,min=8" goal:"secret,hidden"`

	IsSuperuser bool `goal:"readonly"`
	IsActive    bool

	Creator uint `goal:"autowired" gorm:"index"`
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

type Role struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `binding:"required,alphanum,min=3,max=40" gorm:"unique"`
	Creator uint   `goal:"autowired" gorm:"index"`
}

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

type Session struct {
	ID         uint      `goal:"hidden" gorm:"primaryKey"`
	Key        string    `binding:"required,alphanum,len=32" goal:"uuid,readonly" gorm:"unique"`
	UserID     uint      `binding:"required,min=1" goal:"ref=auth.User.Username,hidden,postonly"`
	User       User      `binding:"-" goal:"preload=Username"`
	ExpireDate time.Time `binding:"required" gorm:"index"`
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
