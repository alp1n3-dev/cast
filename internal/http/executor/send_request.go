package executor

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"

	//"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

// Should assume all fields have been created and validated by the time they get here.
func SendRequest(result *models.ExecutionResult, debug, highlight *bool, printOption *[]string, redirectsToFollow *int) error {
	// Going to be a flag later
	//followRedirect := true

	// Going to be a flag later, based on if asserts are detected in the file when read.
	assertsRequired := false
	var err error

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if *redirectsToFollow > 0 {
		err = fasthttp.DoRedirects(result.Request.Req, resp, *redirectsToFollow)
		logging.Logger.Info("Followed redirect\n")
	} else {
		err = fasthttp.Do(result.Request.Req, resp)
	}
	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return err
	}

	// TODO: Implement no-response for printOption flag
	//if *printOption == "no-response" {
	//return nil
	//}

	output.PrintHTTP(nil, resp, highlight, printOption)

	// Blocking off the below section for later with a return that'll be hit
	if !assertsRequired {
		return nil
	}

	var assertion string // Placeholder for if the response values need to be saved and filtered
	if len(assertion) > 0 {
		//fmt.Printf("reached assertion spot")
		result.Response = parse.BuildFastHTTPResponse(resp)
	}

	return nil
}
