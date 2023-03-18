package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/huoyijie/goal/auth"
)

func HomeDir() (homeDir string) {
	homeDir, err := os.UserHomeDir()
	LogFatal(err)
	return
}

func WorkDir() (workDir string) {
	workDir = filepath.Join(HomeDir(), ".goal")
	if _, err := os.Stat(workDir); os.IsNotExist(err) {
		LogFatal(os.Mkdir(workDir, 00744))
	}
	return
}

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func GetWithPrefix(elems []string, prefix string) string {
	for _, s := range elems {
		if r, found := strings.CutPrefix(s, prefix); found {
			return r
		}
	}
	return ""
}

func Group(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(filepath.Base(t.PkgPath()))
}

func Item(model any) string {
	t := reflect.TypeOf(model).Elem()
	return strings.ToLower(t.Name())
}

func Obj(model any) string {
	t := reflect.TypeOf(model).Elem()
	group := strings.ToLower(filepath.Base(t.PkgPath()))
	item := strings.ToLower(t.Name())
	return fmt.Sprintf("%s.%s", group, item)
}

func Actions() []string {
	return []string{"post", "delete", "put", "get"}
}

func Allow(session *auth.Session, obj, act string, enforcer *casbin.Enforcer) bool {
	if session.User.IsSuperuser {
		return true
	}
	if ok, err := enforcer.Enforce(session.Sub(), obj, act); err == nil && ok {
		return true
	}
	return false
}

func AllowAny(session *auth.Session, obj string, enforcer *casbin.Enforcer) bool {
	return Allow(session, obj, "get", enforcer) || Allow(session, obj, "post", enforcer) || Allow(session, obj, "put", enforcer) || Allow(session, obj, "delete", enforcer)
}
