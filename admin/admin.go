package admin

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

var items []any

func group(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(filepath.Base(t.PkgPath()))
}

func item(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(t.Name())
}

func AddItems(models []any) {
	items = models
}

func (OperationLog) TableName() string {
	return "admin_operation_logs"
}

func (*OperationLog) GroupDynamicStrings() (strings []string) {
outer:
	for _, i := range items {
		g := group(i)
		for _, j := range strings {
			if g == j {
				continue outer
			}
		}
		strings = append(strings, g)
	}
	return
}

func (*OperationLog) ItemDynamicStrings() (strings []string) {
	for _, i := range items {
		strings = append(strings, fmt.Sprintf("%s.%s", group(i), item(i)))
	}
	return
}
