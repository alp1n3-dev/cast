package executor

import (
	//"fmt"

	"os"
	"slices"
	"strings"

	//"github.com/alp1n3-eth/cast/internal/http/executor"

	"github.com/alp1n3-eth/cast/internal/capture"
	output "github.com/alp1n3-eth/cast/internal/output"
	"github.com/alp1n3-eth/cast/internal/parse"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/valyala/fasthttp"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(replacementVariables *map[string]string, HTTPCtx *models.HTTPRequestContext) {

	//logging.Init(HTTPCtx.CmdArgs.Debug) // turning this one off for now

	//logging.Init(true)                                                             // temporarily turning debug on
	//HTTPCtx.CmdArgs.PrintOptions = append(HTTPCtx.CmdArgs.PrintOptions, "request") // temporarily turning print request on as well

	logging.Logger.Debugf("req uri directed fasthttp: %s", string(HTTPCtx.Request.Req.RequestURI()))

	if len(string(HTTPCtx.Request.Req.RequestURI())) > 4 {
		HTTPCtx.CmdArgs.URL = string(HTTPCtx.Request.Req.RequestURI())
		HTTPCtx.CmdArgs.Method = string(HTTPCtx.Request.Req.Header.Method())
	}

	logging.Logger.Debugf("Debug: %t, Method: %s, URI: %s", HTTPCtx.CmdArgs.Debug, HTTPCtx.CmdArgs.Method, HTTPCtx.CmdArgs.URL) // appears fine still here

	/*
		hostHeader := HTTPCtx.Request.Req.Header.Peek("Host")
		if len(hostHeader) < 4 {
			//logging.Logger.Debug("Host Header Being Set")
			HTTPCtx.Request.Req.SetRequestURI(string(hostHeader))
			logging.Logger.Debugf("Host Header Set: %s", string(hostHeader))
		}
	*/

	requestURI := string(HTTPCtx.Request.Req.RequestURI())
	if len(requestURI) <= 3 {
		logging.Logger.Error("The request URI appears to be invalid", "err", requestURI)
		return
	}

	logging.Logger.Debug(requestURI)
	if !strings.Contains(requestURI, "http") && (string(requestURI[0]) != "/") {
		logging.Logger.Debug("HTTP Being Added - URI Being Edited")
		//logging.Logger.Warnf("Did you want 'https://' inserted before the URI (%s)? [y/n]: ", requestURI)
		//var userChoice string
		//fmt.Scanln(&userChoice)
		userChoice := "y" // temporarily hardcoding it
		if userChoice == "y" {
			http := "http://" // temporary so that I can see traffic w/o cert/proxy.
			//http := "https://"
			http += requestURI
			HTTPCtx.Request.Req.SetRequestURI(http)
		}

	}

	if HTTPCtx.CmdArgs.CurlOutput {
		curlCmd := generateCurlCommand(HTTPCtx.Request.Req)

		logging.Logger.Debug(curlCmd)
		return
	}

	//result := &models.ExecutionResult{}

	if HTTPCtx.CmdArgs.FileUploadPath != "" {
		logging.Logger.Debugf("Upload file path: %s", HTTPCtx.CmdArgs.FileUploadPath)

		HTTPCtx.CmdArgs.Body = parse.ReadFileIntoBody(&HTTPCtx.CmdArgs.FileUploadPath)
	}

	if len(*replacementVariables) > 0 && (*replacementVariables != nil) {
		logging.Logger.Debugf("Replacement Variables: %s", replacementVariables)
		parse.SwapReqVals(HTTPCtx.Request.Req, replacementVariables)
		logging.Logger.Debug("Executed Successfully: SwapReqVals()")
	}

	if len(capture.GlobalVars) > 0 && (capture.GlobalVars != nil) {
		logging.Logger.Debugf("Global Replacement Variables: %s", capture.GlobalVars)
		parse.SwapReqVals(HTTPCtx.Request.Req, &capture.GlobalVars)
		logging.Logger.Debug("Executed Successfully: SwapReqVals() - Global Replacement")
	}

	logging.Logger.Debugf("Request being sent: %s", HTTPCtx.Request.Req)
	// Needs to be the one directly before sending it, as changes may happen in functions like SwapReqVals().

	if len(HTTPCtx.CmdArgs.PrintOptions) > 0 {
		if slices.Contains(HTTPCtx.CmdArgs.PrintOptions, "request") {
			output.OutputRequest(HTTPCtx.Request.Req, &HTTPCtx.CmdArgs)
		}
	}

	resp, err := SendRequest(HTTPCtx)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			logging.Logger.Error("Lookup failed: no such host. Double-check the entered URI & its extension.")
			return
		}
		logging.Logger.Error("Error sending HTTP request", "err", err)
		return
	}
	//defer fasthttp.ReleaseResponse(resp)
	//fmt.Println(resp)

	HTTPCtx.Response = *resp

	output.OutputResponse(resp, &HTTPCtx.CmdArgs)

	if HTTPCtx.CmdArgs.DownloadPath != "" {
		if err := os.WriteFile(HTTPCtx.CmdArgs.DownloadPath, resp.Body, 0644); err != nil {
			logging.Logger.Error("Problem writing body to download file", "err", requestURI)
			logging.Logger.Error(err)
			return
		}
	}

	logging.Logger.Debug("Executed Successfully: SendRequest()")

	return

}

func generateCurlCommand(req *fasthttp.Request) string {
	var curlCmdStr string

	curlCmdStr = "curl -X " + (string(req.Header.Method())) + " " + (req.URI().String())

	req.Header.VisitAll(func(key, value []byte) {
		curlCmdStr += (" -H '") + string(key) + ": " + string(value) + "'"
	})

	if len(req.Body()) > 0 {
		curlCmdStr += " -d '" + string(req.Body()) + "'"
	}

	/*

		if replacementVariables != nil {
			for key, value := range *replacementVariables {
				curlCmdStr += "--variable '" + key + "=" + value + "'"
			}
		}

	*/

	return curlCmdStr
}

type customParser struct{}

func (p *customParser) RequestSequence(castFile *models.CastFile) {
	// The cast file is handed off to RequestSequence AFTER being read from the file via Koanf. No request building has been done.

	// This function will handle loading the variables, building the requests, sending them, asserting on them, then loop. If syntax is incorrect, it won't be known until that specific request is hit and an error is produced.

	// Use SendHTTP in a loop, pass it the individual requestCTX for each entry.

}
