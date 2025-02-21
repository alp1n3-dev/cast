package executor

import (
	"fmt"
	//"net/http"

	//"net/url"
	//"io"
	//"os"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	//"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/logging"
)

// Should assume all fields have been created and validated by the time they get here.



func SendRequest(result *models.ExecutionResult, debug, highlight *bool, printOption *string) error {

	if *debug || *printOption == "request" {
    	//logging.Logger.Debug("debug == true; printing request ->")
     	output.PrintHTTP(result.Request.Req, nil, highlight)
    }

	req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    //req.SetRequestURI(r.Request.URL.String())
    //req.Header.SetMethod(r.Request.Method)

    //if r.Request.Body != nil {

      //req.SetBodyStream(r.Request.Body, -1)
      //}




    //for name, values := range r.Request.Headers {
       // Loop over all values for the name.
       	//for _, value := range values {
           	//fmt.Printf("%s: %s\n", name, value)
            //req.Header.Add(name, value)
           	//}
            //}



	resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    err := fasthttp.Do(result.Request.Req, resp)
    if err != nil {
        fmt.Printf("Client get failed: %s\n", err)
        return err
    }

    if *printOption == "no-response" {
    	return nil
    }

    output.PrintHTTP(nil, resp, highlight)
    result.Response = parse.BuildFastHTTPResponse(resp)

	return nil
}
