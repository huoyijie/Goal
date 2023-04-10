package auth

import (
	"fmt"

	"gorm.io/gorm/schema"
)

func (u User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

func (User) TableName() string {
	return "auth_users"
}

var _ schema.Tabler = User{}

func (r Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

func (Role) TableName() string {
	return "auth_roles"
}

var _ schema.Tabler = Role{}

func (s Session) Sub() string {
	return s.User.Sub()
}

func (Session) TableName() string {
	return "auth_sessions"
}

var _ schema.Tabler = Session{}
