package admin

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/huoyijie/GoalGenerator/model"
)

var items []any

func AddItems(models []any) {
	items = models
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

func (*OperationLog) TranslateGroupDynamicStrings() map[string]map[string]string {
	r := map[string]map[string]string{"en": {}, "zh-CN": {}}
	for _, i := range items {
		g := group(i)
		if _, found := r["en"][g]; !found {
			t := i.(model.Translate)
			tPkg := t.TranslatePkg()
			r["en"][g] = tPkg["en"]
			r["zh-CN"][g] = tPkg["zh-CN"]
		}
	}
	return r
}

func (*OperationLog) ItemDynamicStrings() (strings []string) {
	for _, i := range items {
		strings = append(strings, item(i))
	}
	return
}

func (*OperationLog) TranslateItemDynamicStrings() map[string]map[string]string {
	r := map[string]map[string]string{"en": {}, "zh-CN": {}}
	for _, i := range items {
		it := item(i)
		if _, found := r["en"][it]; !found {
			t := i.(model.Translate)
			tName := t.TranslateName()
			r["en"][it] = tName["en"]
			r["zh-CN"][it] = tName["zh-CN"]
		}
	}
	return r
}

func group(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(filepath.Base(t.PkgPath()))
}

func item(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(t.Name())
}
