package auth

import (
	"fmt"
)

func (*User) TableName() string {
	return "auth_users"
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

func (*Role) TableName() string {
	return "auth_roles"
}

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

func (*Session) TableName() string {
	return "auth_sessions"
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
