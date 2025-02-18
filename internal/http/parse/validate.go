package parse

import (
	"slices"

	"github.com/alp1n3-eth/cast/pkg/logging"
)

func ValidateMethod(method string) bool {
	methodList := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTION", "TRACE"}

	if slices.Contains(methodList, method) {
		//fmt.Println("[!] Method and URL Provided")
		//logging.Logger.Info("Method Provided")
		logging.Logger.Debug("Method provided")
		// TODO: Parse and send custom headers
		// TODO: Double-check that multiple cookies being set doesn't run into any issues
		// TODO: Add the ability to add a custom body to POSTS
       return true
	}
	logging.Logger.Debug("Method not provided, checking for file extension")
	return false
}

func ValidateFile(arg string) bool {
	if len(arg) > 5 && arg[len(arg)-5:] == ".http" {
		return true
	}
	return false
}
