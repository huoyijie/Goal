package tag

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/huoyijie/goal/util"
)

var COMPONETS = []Component{
	&Calendar{},
	&Dropdown{},
	&Number{},
	&Text{},
	&Password{},
	&Uuid{},
	&Switch{},
}

// define a Component
type Component interface {
	Head() string
	Is(string) bool
}

// define a Tag
type Tag interface {
	Match(string) bool
	Marshal() string
	Unmarshal(string)
}

// define a iface that return itself
type IBase interface {
	Get() *Base
}

// base
type Base struct {
	Autowired `json:",omitempty"`
	Secret    `json:",omitempty"`
	Hidden    `json:",omitempty"`
	Postonly  `json:",omitempty"`
	Readonly  `json:",omitempty"`
	Primary   `json:",omitempty"`
	Unique    `json:",omitempty"`
	BelongTo  *BelongTo `json:",omitempty"`
}

// Get implements IBase
func (b *Base) Get() *Base {
	return b
}

// Marshal implements Tag
func (b *Base) Marshal() string {
	return Marshal(b)
}

// Match implements Tag
func (b *Base) Match(token string) bool {
	return true
}

// Unmarshal implements Tag
func (b *Base) Unmarshal(token string) {
	Unmarshal(token, b)
}

var _ Tag = (*Base)(nil)
var _ IBase = (*Base)(nil)

// Get `key` of a tag
func Key(tag Tag) string {
	return util.ToLowerFirstLetter(reflect.TypeOf(tag).Elem().Name())
}

// Get `key=` of a tag
func Prefix(tag Tag) string {
	return fmt.Sprintf("%s=", Key(tag))
}

// Parse property of string
func ParseString(tag Tag, token string) (propVal string) {
	if tag.Match(token) {
		prefix := Prefix(tag)
		for _, v := range strings.Split(token, ",") {
			if t, found := strings.CutPrefix(v, prefix); found {
				return t
			}
		}
	}
	return
}

// Parse property of int
func ParseInt(tag Tag, token string) (propVal int) {
	if propValStr := ParseString(tag, token); propValStr != "" {
		fmt.Sscanf(propValStr, "%d", &propVal)
	}
	return
}

// recursive marshal
func Marshal(tag Tag) (token string) {
	var arr []string
	elem := reflect.ValueOf(tag).Elem()
	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		if f.Kind() != reflect.Pointer {
			f = f.Addr()
		}
		if f.IsNil() {
			continue
		}
		tag := f.Interface().(Tag)
		if t := tag.Marshal(); t != "" {
			arr = append(arr, t)
		}
	}
	if c, ok := tag.(Component); ok {
		token = c.Head()
	}
	token += strings.Join(arr, ",")
	return
}

// recursive unmarshal
func Unmarshal(token string, tag Tag) {
	if tag.Match(token) {
		if c, ok := tag.(Component); ok {
			token, _ = strings.CutPrefix(token, c.Head())
		}

		elem := reflect.ValueOf(tag).Elem()
		for i := 0; i < elem.NumField(); i++ {
			f := elem.Field(i)
			if f.Kind() != reflect.Pointer {
				f = f.Addr()
			}
			tag := f.Interface().(Tag)
			if tag.Match(token) {
				if f.IsNil() {
					f.Set(reflect.New(reflect.TypeOf(f.Interface()).Elem()))
					tag = f.Interface().(Tag)
				}
				tag.Unmarshal(token)
			}
		}
	}
}

// Get <head> of a component
func ComponentHead(c Component) string {
	return fmt.Sprintf("<%s>", Key(c.(Tag)))
}

// Get to know if the token is a component
func IsComponent(c Component, token string) bool {
	return strings.HasPrefix(token, c.Head())
}
