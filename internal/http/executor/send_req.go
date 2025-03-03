package executor

import (
	"fmt"

	"time"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

// Should assume all fields have been created and validated by the time they get here.
func SendRequest(HTTPCtx *models.HTTPRequestContext) (*models.Response, error) {
	// TODO: Isolate unrelated logic to its own files. Find a better / more efficient way to "keep" / store the resp that isn't "expensive", so that printing doesn't need to be called from within this SendRequest().

	var storeResp models.Response

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	//response := &models.Response{}
	// Going to be a flag later, based on if asserts are detected in the file when read.
	assertsRequired := false
	var err error

	// TODO: Add the ability to track byte sizes of responses.

	//req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(HTTPCtx.Request.Req)

	//resp := fasthttp.AcquireResponse()
	//defer fasthttp.ReleaseResponse(resp)

	startTime := time.Now()

	if HTTPCtx.CmdArgs.RedirectsToFollow > 0 {
		if err = fasthttp.DoRedirects(HTTPCtx.Request.Req, resp, HTTPCtx.CmdArgs.RedirectsToFollow); err != nil {
			return &storeResp, err
		}
		logging.Logger.Info("Success: Followed redirect")
	} else {
		if err = fasthttp.Do(HTTPCtx.Request.Req, resp); err != nil {
			return &storeResp, err
		}
		logging.Logger.Debug("Success: Request sent")
	}

	//bytetype := []byte(client.Resp.String())
	//var RespStore models.Response

	//client.RespStore = []byte(client.Resp.String())
	//storeResp.Headers = resp.Header.String()
	//storeResp.StatusCode = resp.StatusCode()
	//storeResp.Status = string(resp.Header.StatusMessage())
	//storeResp.Body = resp.Body()

	duration := time.Since(startTime).Milliseconds()
	//storeResp.Duration = int(duration)

	buildResp(resp, duration, &storeResp)

	//fmt.Println(resp)

	// Blocking off the below section for later with a return that'll be hit
	if assertsRequired {
		fmt.Println("reached inside of the if assertsReq")
		return &storeResp, nil
	}

	var assertion string // Placeholder for if the response values need to be saved and filtered
	if len(assertion) > 0 {
		//*response = parse.BuildFastHTTPResponse(resp)
		//TODO: do actions following assertion logic
	}

	//var respStore fasthttp.Response
	//respStore = *resp
	fmt.Println("reached end in send_rest.go")
	//fmt.Println(storeResp)

	return &storeResp, nil
}

/*
func SaveResponseToFile(resp *fasthttp.Response, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(string(resp.Body()))) // Convert body to io.Reader
	if err != nil {
		return fmt.Errorf("failed to write response body to file: %w", err)
	}

	return nil
}
*/

func buildResp(resp *fasthttp.Response, duration int64, storeResp *models.Response) {

	storeResp.Headers = resp.Header.String()
	storeResp.StatusCode = resp.StatusCode()
	storeResp.Status = string(resp.Header.StatusMessage())
	storeResp.Body = resp.Body()
	storeResp.Duration = int(duration)

	return
}
