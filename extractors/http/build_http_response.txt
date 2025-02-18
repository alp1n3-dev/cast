package http_extractors

import (
	"net/http"
	"io"
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/models"
)

/*
Meant to be run after the request has been sent.
*/

func BuildHTTPResponse (ogRequest models.ExecutionResult, resp *http.Response) (models.ExecutionResult, error) {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ogRequest, fmt.Errorf("failed to read response body: %w", err)
	}

	ogRequest.Response.Body = body
	ogRequest.Response.StatusCode = resp.StatusCode
	ogRequest.Response.Status = resp.Status

	/* TODO: Review and rework, as headers type has changed in models/http.go. */
	//if ogRequest.Response.Headers == nil {
    	//ogRequest.Response.Headers = make(map[string]string)
	//}
	//for key, values := range resp.Header {
		//for _, value := range values {
			//fmt.Printf("%s: %s\n", key, value)
			//ogRequest.Response.Headers[key] = value
		//}
	//}
	ogRequest.Response.Headers = resp.Header



	return ogRequest, nil
}
