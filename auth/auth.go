package auth

import (
	"fmt"
)

func (r *Role) RoleID() string {
	return fmt.Sprintf("role-%d", r.ID)
}

func (u *User) Sub() string {
	return fmt.Sprintf("user-%d", u.ID)
}

func (s *Session) Sub() string {
	return s.User.Sub()
}
