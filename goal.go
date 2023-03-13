package goal

import (
	"path/filepath"
	"reflect"
	"sort"

	"github.com/huoyijie/goal/auth"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var models = []any{
	&auth.User{},
	&auth.Role{},
	&auth.Session{},
}

type item struct {
	Name string
}

type group struct {
	Name  string
	items []item
}

func groupList() (groups []group) {
	dict := map[string][]item{}

	for _, model := range models {
		elem := reflect.TypeOf(model).Elem()
		pkgName := cases.Title(language.Und).String(filepath.Base(elem.PkgPath()))
		if items, ok := dict[pkgName]; ok {
			dict[pkgName] = append(items, item{
				Name: elem.Name(),
			})
		} else {
			dict[pkgName] = []item{{Name: elem.Name()}}
		}
	}

	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, v := range keys {
		items := dict[v]
		sort.Slice(items, func(i, j int) (less bool) {
			less = items[i].Name < items[j].Name
			return
		})
		groups = append(groups, group{v, items})
	}
	return
}

func Register(modelList ...any) {
	models = append(models, modelList)
}

func Models() []any {
	return models
}
