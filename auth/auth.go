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

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Session struct {
	ID         uint   `gorm:"primaryKey"`
	Key        string `gorm:"unique"`
	UserID     uint
	User       User      `goal:"preload=Username"`
	ExpireDate time.Time `gorm:"index"`
}

func (s *Session) Sub() string {
	return fmt.Sprintf("user-%d", s.UserID)
}
