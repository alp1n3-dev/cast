package http_extractors

import (
	"log"
	"errors"
	"net/url"
)

func ValidateHTTP(method, url string) {
	validMethod := isValidMethod(method)
	if validMethod != true {
		err := errors.New("Invalid Method.")
		log.Fatalln(err)
	}

	validURL := isHTTPURL(url)
	if validURL != true {
		err := errors.New("Invalid URL.")
		log.Fatalln(err)
	}
}

func isValidMethod (method string) bool {
	switch method {
		case "GET", "POST", "PUT", "DELETE", "OPTION", "TRACE", "HEAD":
    		// Valid method, do nothing
    		return true
    	default:
   			return false
	}
}

func isHTTPURL (str string) bool {
	u, err := url.Parse(str)
    return err == nil && u.Scheme != "" && u.Host != ""
}
