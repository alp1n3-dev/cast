package http_extractors

import (
	"net/http"
	"io"
	"fmt"

	"github.com/alp1n3-eth/cast/models"
)

/*
Meant to be run after the request has been sent.
*/

func BuildHTTPResponse (ogRequest models.HTTPRequest, resp *http.Response) (models.HTTPRequest, error) {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ogRequest, fmt.Errorf("failed to read response body: %w", err)
	}

	ogRequest.Response.Body = string(body)
	ogRequest.Response.StatusCode = resp.Status

	if ogRequest.Response.Headers == nil {
    ogRequest.Response.Headers = make(map[string]string)
	}
	for key, values := range resp.Header {
		for _, value := range values {
			//fmt.Printf("%s: %s\n", key, value)
			ogRequest.Response.Headers[key] = value
		}
	}


	return ogRequest, nil
}
