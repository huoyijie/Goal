package tag

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/huoyijie/Goal/util"
)

var COMPONETS = []Component{
	&Calendar{},
	&Dropdown{},
	&Number{},
	&Text{},
	&Password{},
	&Uuid{},
	&Switch{},
	&File{},
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

// recursive marshal
func Marshal(tag Tag) (token string) {
	var arr []string
	t := reflect.TypeOf(tag).Elem()
	v := reflect.ValueOf(tag).Elem()
loop:
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fVal := v.FieldByName(f.Name)
		switch fVal.Kind() {
		case reflect.Bool:
			if fVal.Bool() {
				arr = append(arr, util.ToLowerFirstLetter(f.Name))
			}
		case reflect.Struct:
			if f.Name == "Base" {
				if tag := fVal.Addr().Interface().(Tag).Marshal(); tag != "" {
					arr = append(arr, tag)
				}
			}
		case reflect.Pointer:
			if fVal.IsNil() {
				continue loop
			}
			switch fVal.Elem().Kind() {
			case reflect.Int:
				arr = append(arr, fmt.Sprintf("%s=%d", util.ToLowerFirstLetter(f.Name), fVal.Elem().Int()))
			case reflect.Struct:
				if f.Name == "BelongTo" || f.Name == "HasOne" || f.Name == "UploadTo" {
					if tag := fVal.Interface().(Tag).Marshal(); tag != "" {
						arr = append(arr, tag)
					}
				}
			}
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

		t := reflect.TypeOf(tag).Elem()
		v := reflect.ValueOf(tag).Elem()
	loop:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fVal := v.FieldByName(f.Name)
			switch fVal.Kind() {
			case reflect.Bool:
				fn := util.ToLowerFirstLetter(f.Name)
				for _, s := range strings.Split(token, ",") {
					if s == fn {
						fVal.SetBool(true)
						break
					}
				}
			case reflect.Struct:
				if f.Name == "Base" {
					fVal.Addr().Interface().(Tag).Unmarshal(token)
				}
			case reflect.Pointer:
				if prefix := fmt.Sprintf("%s=", util.ToLowerFirstLetter(f.Name)); strings.Contains(token, prefix) {
					if fVal.IsNil() {
						fVal.Set(reflect.New(f.Type.Elem()))
					}
					switch fVal.Elem().Kind() {
					case reflect.Int:
						for _, v := range strings.Split(token, ",") {
							if t, found := strings.CutPrefix(v, prefix); found {
								if t != "" {
									fmt.Sscanf(t, "%d", fVal.Interface())
								}
								continue loop
							}
						}
					case reflect.Struct:
						if f.Name == "BelongTo" || f.Name == "HasOne" || f.Name == "UploadTo" {
							fVal.Interface().(Tag).Unmarshal(token)
						}
					}
				}
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
