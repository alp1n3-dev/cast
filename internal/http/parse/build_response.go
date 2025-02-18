package parse

import (
	"net/http"
	"fmt"
	"io"


	"github.com/alp1n3-eth/cast/pkg/models"
)

func BuildResponse(response *http.Response) models.Response {
	var builtResponse models.Response
	var err error

	builtResponse.StatusCode = response.StatusCode
	builtResponse.Status = response.Status
	builtResponse.Body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Problem reading response body in BuildResponse")
	}
	//fmt.Println(io.ReadAll(response.Body))
	builtResponse.Headers = response.Header

	return builtResponse
}
