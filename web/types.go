package web

import (
	"fmt"
	"strings"

	"github.com/huoyijie/Goal/web/tag"
)

const (
	ErrInvalidUsernameOrPassword int = -(iota + 10000)
	ErrInvalidPassword
)

type SigninForm struct {
	Username   string `binding:"required,alphanum,min=3,max=40"`
	Password   string `binding:"required,min=8"`
	RememberMe bool
}

type ChangePasswordForm struct {
	Password, NewPassword string `binding:"required,min=8"`
}

type Result struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type Component struct {
	Name string
	Tag  tag.Component
}

type Column struct {
	Name         string
	Component    Component
	ValidateRule string
}

type Perms struct {
	Post   bool `json:"post"`
	Delete bool `json:"delete"`
	Put    bool `json:"put"`
	Get    bool `json:"get"`
}

type DataTable struct {
	Lazy    bool     `json:"lazy"`
	Columns []Column `json:"columns"`
	Perms   *Perms   `json:"perms"`
}

type RecordList struct {
	Total uint `json:"total"`
	List  any  `json:"list"`
}

type Menu struct {
	Name  string `json:"label"`
	Items []Menu `json:"items"`
}

type Perm struct {
	Code,
	Group,
	Item,
	Action string
}

func NewPerm(obj, action string) Perm {
	arr := strings.Split(obj, ".")
	return Perm{
		Code:   fmt.Sprintf("%s:%s", obj, action),
		Group:  arr[0],
		Item:   arr[1],
		Action: action,
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

type LazyParam struct {
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
	SortField string `form:"sortField"`
	SortOrder int    `form:"sortOrder"`
	Filters   string `form:"filters"`
}

type Option struct {
	Label string `json:"label"`
	Value any `json:"value"`
}
