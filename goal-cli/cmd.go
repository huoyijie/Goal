package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/huoyijie/goal"
	"github.com/huoyijie/goal/auth"
	"github.com/huoyijie/goal/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Username string `validate:"required,alphanum,min=3,max=40"`
	Email    string `validate:"required,email"`
}

func readPassword() (pw []byte) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	util.LogFatal(err)
	return
}

func main() {
	var user User
	flag.StringVar(&user.Username, "username", "", "username of the superuser")
	flag.StringVar(&user.Email, "email", "", "email of the superuser")
	flag.Parse()

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		util.LogFatal(err)
	}

	var rawPassword []byte
	for {
		fmt.Print("Password:")
		rawPassword = readPassword()
		fmt.Println()

		fmt.Print("Password (again):")
		rawPasswordAgain := readPassword()
		fmt.Println()

		if bytes.Equal(rawPassword, rawPasswordAgain) {
			break
		}

		// todo translate
		fmt.Println("Error: Your passwords didn't match.")
	}

	if err := validate.Var(string(rawPassword), "required,min=8"); err != nil {
		util.LogFatal(err)
	}

	bcryptHash, err := bcrypt.GenerateFromPassword(rawPassword, 14)
	util.LogFatal(err)

	db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
	util.LogFatal(err)

	db.AutoMigrate(goal.Models()...)

	superuser := &auth.User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    string(bcryptHash),
		DateJoined:  time.Now(),
		IsSuperuser: true,
		IsStaff:     true,
		IsActive:    true,
	}

	err = db.Create(superuser).Error
	util.LogFatal(err)
}
