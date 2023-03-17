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
	Lastname string `goal:"hidden"`

	DateJoined,
	LastSignin time.Time `gorm:"null" goal:"hidden"`

	IsSuperuser,
	IsStaff,
	IsActive bool

	Roles []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type Session struct {
	ID         string `gorm:"primaryKey"`
	UserID     uint
	User       User      `goal:"preload=Username"`
	ExpireDate time.Time `gorm:"index"`
}
