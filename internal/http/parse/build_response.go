package parse

import (
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

func BuildFastHTTPResponse(response *fasthttp.Response) models.Response {
	var builtResponse models.Response

	builtResponse.StatusCode = response.StatusCode()

	builtResponse.Protocol = string(response.Header.Protocol())
	builtResponse.ContentType = string(response.Header.ContentType())
	builtResponse.Status = string(response.Header.StatusMessage())
	builtResponse.Body = response.Body()

	if builtResponse.Headers == nil {
		builtResponse.Headers = http.Header{}
	}

	// why is the function above and below both here? Not 100% sure what I was doing here.

	// TODO: Possible to log request and response to a .HAR file or something? What's the most efficient format and way to do this?

	response.Header.VisitAll(func(key, value []byte) {
		//fmt.Printf("Adding header: %s: %s\n", string(key), string(value))
		logging.Logger.Debugf("Building Response Header Key: %s, Value: %s", key, value)
		builtResponse.Headers.Add(string(key), string(value))
	})

	logging.Logger.Debug("Reached end of BuildFastHTTPResponse")

	return builtResponse
}
