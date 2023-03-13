package auth

import (
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Username,
	Email string `gorm:"unique"`
	Password,

	Firstname,
	Lastname string

	DateJoined time.Time
	LastLogin  time.Time `gorm:"null"`

	IsSuperuser,
	IsStaff,
	IsActive bool

	Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type Session struct {
	ID         string `gorm:"primaryKey;size:32"`
	User       User
	ExpireDate time.Time
}
