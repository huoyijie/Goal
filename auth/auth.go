package auth

import (
	"fmt"
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username,
	Email string `gorm:"unique"`
	Password string `goal:"hidden"`

	IsSuperuser,
	IsActive bool
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

type Session struct {
	ID         uint   `gorm:"primaryKey"`
	Key        string `gorm:"unique"`
	UserID     uint
	User       User      `goal:"preload=Username"`
	ExpireDate time.Time `gorm:"index"`
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
