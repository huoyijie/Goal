package auth

import (
	"fmt"
	"time"
)

type User struct {
	ID uint `goal:"<number>primary" gorm:"primaryKey"`

	Username string `validate:"required,alphanum,min=3,max=40" binding:"required,alphanum,min=3,max=40" goal:"<text>unique" gorm:"unique"`
	Email    string `validate:"required,email" binding:"required,email" goal:"<text>unique" gorm:"unique"`
	Password string `validate:"required,min=8" binding:"required,min=8" goal:"<password>secret,hidden"`

	IsSuperuser bool `goal:"<switch>readonly"`
	IsActive    bool `goal:"<switch>"`

	Creator uint `goal:"<number>autowired" gorm:"index"`
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

type Role struct {
	ID      uint   `goal:"<number>primary" gorm:"primaryKey"`
	Name    string `binding:"required,alphanum,min=3,max=40" goal:"<text>unique" gorm:"unique"`
	Creator uint   `goal:"<number>autowired" gorm:"index"`
}

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

type Session struct {
	ID         uint      `goal:"<number>primary,hidden" gorm:"primaryKey"`
	Key        string    `binding:"required,alphanum,len=32" goal:"<uuid>readonly,unique" gorm:"unique"`
	UserID     uint      `goal:"<number>autowired"`
	User       User      `binding:"required" goal:"<dropdown>postonly,belongTo=auth.User.Username,filter"`
	ExpireDate time.Time `binding:"required" goal:"<calendar>showTime,showIcon" gorm:"index"`
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
