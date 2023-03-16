package auth

import (
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username,
	Email string `gorm:"unique"`
	Password string `goal:"hidden"`

	Firstname,
	Lastname string

	DateJoined time.Time
	LastSignin time.Time `gorm:"null"`

	IsSuperuser,
	IsStaff,
	IsActive bool

	Roles []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Session struct {
	ID         string `gorm:"primaryKey;size:32"`
	UserID     uint
	User       User
	ExpireDate time.Time `gorm:"index"`
}
