package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func OpenSqliteDB() (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(WorkDir(), "goal.db")), &gorm.Config{})
	LogFatal(err)
	return
}

func BcryptHash(rawPassword string) string {
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 14)
	LogFatal(err)
	return string(bcryptHash)
}
