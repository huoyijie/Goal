package web

import (
	"fmt"
	"strings"
)

type Code int

const (
	ErrInvalidUsernameOrPassword Code = -(iota + 10000)
)

type SigninForm struct {
	Username   string `json:"username" binding:"required,alphanum,min=3,max=40"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
}

type Column struct {
	Name,
	Type string
	Hidden,
	Primary,
	Unique,
	Preload bool
	PreloadField,
	ValidateRule string
}

type Menu struct {
	Name  string `json:"label"`
	Items []Menu `json:"items"`
}

type Perm struct {
	Code,
	Name string
}

func NewPerm(obj, action string) Perm {
	arr := strings.Split(obj, ".")
	return Perm{
		Code: fmt.Sprintf("%s:%s", obj, action),
		Name: fmt.Sprintf("%s %s %s", arr[0], arr[1], action),
	}
}

func (p *Perm) Val() []string {
	return strings.Split(p.Code, ":")
}

type PermsParam struct {
	RoleID uint `uri:"roleID"`
}

type RolesParam struct {
	UserID uint `uri:"userID"`
}
