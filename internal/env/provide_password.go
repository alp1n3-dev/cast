package env

import (
	"fmt"
	"syscall"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"golang.org/x/term"
)

func RetrievePasswordFromUser() (string, error) {

	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logging.Logger.Fatal("Failed to read password")
		return "", err
	}

	fmt.Println()

	passwordStr := string(password)

	return passwordStr, nil
}

func CreateNewPassword() (string, error) {
	fmt.Print("Please provide a new password: ")
	passOne, _ := RetrievePasswordFromUser()

	fmt.Print("Please re-enter the password: ")
	passTwo, _ := RetrievePasswordFromUser()

	if passOne != passTwo {
		logging.Logger.Warn("Passwords do not match")
		// TODO: Implement error return here.
		return "", nil
	} else {
		return passOne, nil
	}
}
