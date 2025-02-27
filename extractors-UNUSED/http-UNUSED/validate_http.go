package http_extractors

import (
	//"log"
	"errors"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"net/url"
)

func ValidateHTTP(method, url string) {
	validMethod := isValidMethod(method)
	if validMethod != true {

		err := errors.New("Invalid Method.")
		logging.Logger.Fatal(err)
	}

	validURL := isHTTPURL(url)
	if validURL != true {
		err := errors.New("Invalid URL.")
		logging.Logger.Fatal(err)
	}
}

func isValidMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "OPTION", "TRACE", "HEAD":
		// Valid method, do nothing
		return true
	default:
		return false
	}
}

func isHTTPURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
