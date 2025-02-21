package executor

import (
	"fmt"
	"net/http"

	//"net/url"
	//"io"
	"os"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/logging"
)

// Should assume all fields have been created and validated by the time they get here.
func SendHTTPRequest(r models.ExecutionResult) (models.ExecutionResult, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Request.Method, r.Request.URL.String(), r.Request.Body)
	if err != nil {
			fmt.Println("Error creating request:", err)
			return r, nil
	}

	logging.Logger.Debug(r.Request.Headers)
	req.Header = r.Request.Headers

	// TODO: Only way host can be set is req.Host = "domain.tld" apparently. Check into it.

	resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			os.Exit(1)
		}
	defer resp.Body.Close()




	//rPair, err := http_extractors.BuildHTTPResponse(r, resp)
	r.Response = parse.BuildResponse(resp)

	// Print the Response
	// Returns *http.Response & error
	//err = output.PrintResponse(r)

	return r, nil
}

func SendFastHTTPRequest(r *models.ExecutionResult, debug, highlight *bool) (*models.ExecutionResult, error) {
	//logging.Logger.Debug("SendFastHTTPRequest - Point 1")

	//var body []byte

	req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)

    //logging.Logger.Debug("SendFastHTTPRequest - Point 2")

    req.SetRequestURI(r.Request.URL.String())
    req.Header.SetMethod(r.Request.Method)

    if r.Request.Body != nil {
   		//body, _ = io.ReadAll(r.Request.Body)
     	//req.SetBody(body)

      req.SetBodyStream(r.Request.Body, -1)
    }


    //logging.Logger.Debug("SendFastHTTPRequest - Point 3")

    for name, values := range r.Request.Headers {
       // Loop over all values for the name.
       	for _, value := range values {
           	//fmt.Printf("%s: %s\n", name, value)
            req.Header.Add(name, value)
        	}
	}

	//logging.Logger.Debug("SendFastHTTPRequest - Point 4")

	resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    err := fasthttp.Do(req, resp)
    if err != nil {
        fmt.Printf("Client get failed: %s\n", err)
        return r, err
    }

    if *debug {
    	logging.Logger.Debug("Debug true, printing request")
    	//logging.Logger.Debug(req)
     	output.PrintRequest(req, highlight)
    }


    //logging.Logger.Debug("SendFastHTTPRequest - Point 5")


   	//body = resp.Body()
    //fmt.Println(string(body[:]))

    //logging.Logger.Debug("SendFastHTTPRequest - Point 6")

    r.Response = parse.BuildFastHTTPResponse(resp)

    //logging.Logger.Debug("SendFastHTTPRequest - Point 7")

	return r, nil
}
