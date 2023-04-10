package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/huoyijie/Goal"
	"github.com/huoyijie/Goal/auth"
	"github.com/huoyijie/Goal/util"
	"golang.org/x/term"
)

func readPassword() (pw []byte) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	util.LogFatal(err)
	return
}

func main() {
	var user auth.User
	flag.StringVar(&user.Username, "username", "", "username of the superuser")
	flag.StringVar(&user.Email, "email", "", "email of the superuser")
	flag.Parse()

	validate := validator.New()

	if err := validate.Var(user.Username, "required,alphanum,min=3,max=40"); err != nil {
		util.LogFatal(err)
	}

	if err := validate.Var(user.Email, "required,email"); err != nil {
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

	super := &auth.User{
		Username:    user.Username,
		Email:       user.Email,
		Password:    util.BcryptHash(string(rawPassword)),
		IsSuperuser: true,
		IsActive:    true,
	}

	goal.NewGoal(util.OpenSqliteDB()).NewSuper(super)
}
