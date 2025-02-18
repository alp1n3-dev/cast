package parse

import (
	"net/http"
	"io"
	"net/url"
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"
)

func BuildRequest (method, urlVal string, body io.Reader, headers http.Header) models.Request {
	var req models.Request
	var err error

	req.Method = method
	req.Body = body

	if headers != nil {
		req.Headers = make(http.Header)
		req.Headers = headers // http.Headers will panic if nil.
	}


	//logging.Logger.Debug("Assigned request method, body headers. About to assign the URL")

	req.URL, err = url.Parse(urlVal)
	if err != nil {
		fmt.Printf("%s", err)
		logging.Logger.Fatal("Unable to parse provided URL")
	}

	if req.Headers == nil {
		//logging.Logger.Debug("Adding Content-Type header")

		req.Headers = make(http.Header)
		req.Headers.Add("Content-Type", "text/html")
	}

	logging.Logger.Debug("Assigned request URL successfully")

	return req
}
