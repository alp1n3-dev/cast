package http_extractors

// This is to be used post-validation (after validate_http.go).

import (
	"net/url"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

func BuildHTTPRequest (method, urlString string) models.ExecutionResult {
	ValidateHTTP(method, urlString)

	var request models.ExecutionResult
	var err error

	request.Request.Method = method
	request.Request.URL, err = url.Parse(urlString)
	if err != nil {
		logging.Logger.Fatal("URL unable to be read into request structure.")
	}

	return request
}
