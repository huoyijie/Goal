package auth

import (
	"fmt"
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username string `validate:"required,alphanum,min=3,max=40" binding:"required,alphanum,min=3,max=40" gorm:"unique"`
	Email    string `validate:"required,email" binding:"required,email" gorm:"unique"`
	Password string `goal:"hidden"`

	IsSuperuser,
	IsActive bool
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `binding:"required,alphanum,min=3,max=40" gorm:"unique"`
}

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

type Session struct {
	ID         uint      `gorm:"primaryKey"`
	Key        string    `binding:"required,alphanum,len=32" gorm:"unique"`
	UserID     uint      `binding:"required,min=1"`
	User       User      `binding:"-" goal:"preload=Username"`
	ExpireDate time.Time `binding:"required" gorm:"index"`
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
