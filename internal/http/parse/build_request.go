package parse

import (
	"sync"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/valyala/fasthttp"
)

func BuildRequest(method, urlStr *string, body *[]byte, headers *map[string][]byte) *fasthttp.Request {
	req := &fasthttp.Request{}

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		logging.Logger.Debugf("Setting method: %s", *method)
		req.Header.SetMethod(*method)
	}()

	go func() {
		defer wg.Done()
		logging.Logger.Debugf("Setting body: %x", *body)
		//if body != nil {
		req.SetBody(*body)
		//req.SetBodyStream(body, -1)
		//}
	}()

	go func() {
		defer wg.Done()
		logging.Logger.Debugf("Setting uri: %s", *urlStr)
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
				//fmt.Println("reached headers")
				// Loop over all values for the name.
				logging.Logger.Debugf("Setting Header: Key: %s, Val: %s", key, value)
				//req.Header.Add(key, string(value))
				//req.Header.Set(key, value)

				req.Header.SetBytesV(key, value)
			}
		}
		if req.Header.Peek("Content-Type") == nil {
			req.Header.Add("Content-Type", "text/html")
		}
	}()

	wg.Wait()

	return req
}
