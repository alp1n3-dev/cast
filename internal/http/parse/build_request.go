package parse

import (
	"net/http"
	"io"
	"net/url"
	"fmt"

	"sync"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"
)

func BuildRequest (method, urlVal *string, body *io.Reader, headers *http.Header) models.Request {
	var req models.Request
	var err error

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	go func() {
		defer wg.Done()
		req.Method = *method
		req.Body = *body
	}()

	go func() {
		defer wg.Done()
		logging.Logger.Debug(headers)
		if headers != nil {
			//req.Headers = make(http.Header)
			logging.Logger.Debug(headers)
			mu.Lock()
			req.Headers = *headers // http.Headers will panic if nil.
			mu.Unlock()

			if _, ok := req.Headers["Content-Type"]; !ok {
			mu.Lock()
    		req.Headers.Add("Content-Type", "text/html")
      		mu.Unlock()
			}
		}
	}()

	wg.Wait()




	logging.Logger.Debug(req.Headers)


	//logging.Logger.Debug("Assigned request method, body headers. About to assign the URL")

	req.URL, err = url.Parse(*urlVal)
	if err != nil {
		fmt.Printf("%s", err)
		logging.Logger.Fatal("Unable to parse provided URL")
	}

	/*
	if req.Headers == nil {
		//logging.Logger.Debug("Adding Content-Type header")

		req.Headers = make(http.Header)
		req.Headers.Add("Content-Type", "text/html")
	}
	*/

	logging.Logger.Debug("Assigned request URL successfully")

	return req
}
