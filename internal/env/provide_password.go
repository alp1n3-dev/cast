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
