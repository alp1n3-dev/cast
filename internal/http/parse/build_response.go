package parse

import (
	"net/http"
	//"fmt"
	//"io"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/logging"
)

/*
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
*/

func BuildFastHTTPResponse(response *fasthttp.Response) models.Response {
	var builtResponse models.Response
	//var err error

	//logging.Logger.Debug("BuildFastHTTPResponse - Point 1")
	//logging.Logger.Debug(builtResponse.StatusCode)
	//logging.Logger.Debug(response.StatusCode())

	builtResponse.StatusCode = response.StatusCode()

	builtResponse.Protocol = string(response.Header.Protocol())
	builtResponse.ContentType = string(response.Header.ContentType())

	//logging.Logger.Debug(builtResponse.StatusCode)
	//logging.Logger.Debug(response.StatusCode())

	builtResponse.Status = string(response.Header.StatusMessage())

	builtResponse.Body = response.Body()



    //fmt.Println(string(body[:]))
	//if err != nil {
		//fmt.Println("Problem reading response body in BuildResponse")
		//}
	//fmt.Println(io.ReadAll(response.Body))
	if builtResponse.Headers == nil {
			builtResponse.Headers = http.Header{}
	}

	// why is the function above and below both here? Not 100% sure what I was doing here.



	response.Header.VisitAll(func(key, value []byte){
		//fmt.Printf("Adding header: %s: %s\n", string(key), string(value))
		logging.Logger.Debugf("Building Response Header Key: %s, Value: %s", key, value)
		builtResponse.Headers.Add(string(key), string(value))
	})

	logging.Logger.Debug("Reached end of BuildFastHTTPResponse")

	return builtResponse
}
