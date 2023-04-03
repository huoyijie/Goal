package util

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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

func ToLowerFirstLetter(str string) string {
	a := []rune(str)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func ToUpperFirstLetter(str string) string {
	a := []rune(str)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}