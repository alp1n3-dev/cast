package parse

import (
	"fmt"
	//"io"
	//"net/http"
	//"net/url"


	"sync"
	//"sync/Mutex"

	"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
)
/*
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
/*
	logging.Logger.Debug("Assigned request URL successfully")

	return req
}
*/
func BuildRequest1 (method, urlStr, body *string, headers *map[string]string) (*fasthttp.Request) {
	req := &fasthttp.Request{}

	logging.Logger.Debug("BuildRequest point 1")

	var wg sync.WaitGroup
	//var mu sync.Mutex
	wg.Add(4)

	go func() {
		defer wg.Done()
		logging.Logger.Debug("Setting method")
		req.Header.SetMethod(*method)
	}()

	go func() {
		defer wg.Done()
		logging.Logger.Debug("Setting body")
		if body != nil {
			req.SetBody([]byte(*body))
			//req.SetBodyStream(body, -1)
		}
	}()



	go func() {
		defer wg.Done()
		logging.Logger.Debug("Setting uri")
		uri := fasthttp.AcquireURI()
		defer fasthttp.ReleaseURI(uri)
		uri.Parse(nil, []byte(*urlStr))
		req.SetURI(uri)
	}()

	go func() {
		defer wg.Done()
		logging.Logger.Debug("Setting headers")
		if headers != nil {
			for key, value := range *headers {
				fmt.Println("reached headers")
    		// Loop over all values for the name.
            req.Header.Add(key, string(value))
			}
		}
		if req.Header.Peek("Content-Type") == nil {
			req.Header.Add("Content-Type", "text/html")
		}
	}()

	wg.Wait()

	return req
}

func BuildRequest2 (method, urlStr, body *string, headers *map[string]string) (*fasthttp.Request) {
	req := &fasthttp.Request{}

	logging.Logger.Debug("BuildRequest point 1")

	req.Header.SetMethod(*method)

	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)
	uri.Parse(nil, []byte(*urlStr))
	req.SetURI(uri)

	logging.Logger.Debug("BuildRequest point 2")
	//fmt.Println(headers)

	if headers != nil {
		for key, value := range *headers {
			fmt.Println("reached headers")
    		// Loop over all values for the name.
            req.Header.Add(key, string(value))
		}
	}

	logging.Logger.Debug("BuildRequest point 3")

	if req.Header.Peek("Content-Type") == nil {
		req.Header.Add("Content-Type", "text/html")
	}

	if body != nil {
		req.SetBody([]byte(*body))
		//req.SetBodyStream(body, -1)
	}

	return req
}
