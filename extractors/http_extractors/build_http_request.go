package http_extractors

// This is to be used post-validation (after validate_http.go).

import (
	"github.com/alp1n3-eth/cast/models"
)

func BuildHTTPRequest (method, url string) models.HTTPRequest {
	ValidateHTTP(method, url)

	var request models.HTTPRequest

	request.Request.Method = method
	request.Request.URL = url

	return request
}
