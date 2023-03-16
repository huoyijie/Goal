package util

import (
	"log"
	"strings"
)

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
