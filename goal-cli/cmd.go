package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/huoyijie/goal/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func readPassword() (pw []byte) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	util.LogFatal(err)
	return
}

func main() {
	var (
		username, email string
	)
	flag.StringVar(&username, "username", "", "username of the superuser")
	flag.StringVar(&email, "email", "", "email of the superuser")
	flag.Parse()

	// todo regex validate
	if len(username) == 0 || len(email) == 0 {
		flag.PrintDefaults()
		return
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

	bcryptHash, err := bcrypt.GenerateFromPassword(rawPassword, 14)
	util.LogFatal(err)

	password := string(bcryptHash)
	fmt.Println(password)
}
