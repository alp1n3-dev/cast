package cmd

import (
	//"fmt"
	"os"
	"slices"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(replacementVariables *map[string]string, HTTPCtx *models.HTTPRequestContext) {

	logging.Init(HTTPCtx.CmdArgs.Debug)

	logging.Logger.Debugf("Debug: %t, Method: %s, URI: %s", HTTPCtx.CmdArgs.Debug, HTTPCtx.CmdArgs.Method, HTTPCtx.CmdArgs.URL)

	//result := &models.ExecutionResult{}

	if HTTPCtx.CmdArgs.FileUploadPath != "" {
		logging.Logger.Debugf("Upload file path: %s", HTTPCtx.CmdArgs.FileUploadPath)

		HTTPCtx.CmdArgs.Body = parse.ReadFileIntoBody(&HTTPCtx.CmdArgs.FileUploadPath)
	}

	if len(*replacementVariables) > 0 {
		logging.Logger.Debugf("Replacement Variables: %s", replacementVariables)
		parse.SwapReqVals(HTTPCtx.Request.Req, replacementVariables)
		logging.Logger.Debug("Executed Successfully: SwapReqVals()")
	}

	logging.Logger.Debugf("Request being sent: %s", HTTPCtx.Request.Req)
	// Needs to be the one directly before sending it, as changes may happen in functions like SwapReqVals().

	if len(HTTPCtx.CmdArgs.PrintOptions) > 0 {
		if slices.Contains(HTTPCtx.CmdArgs.PrintOptions, "request") {
			output.OutputRequest(HTTPCtx.Request.Req, &HTTPCtx.CmdArgs)
		}
	}

	resp, err := executor.SendRequest(HTTPCtx)
	if err != nil {
		logging.Logger.Error("Error sending HTTP request", "err", err)
	}
	//defer fasthttp.ReleaseResponse(resp)
	//fmt.Println(resp)

	output.OutputResponse(resp, &HTTPCtx.CmdArgs)

	if HTTPCtx.CmdArgs.DownloadPath != "" {
		if err := os.WriteFile(HTTPCtx.CmdArgs.DownloadPath, resp.Body, 0644); err != nil {
			logging.Logger.Error(err)
			return
		}
	}

	logging.Logger.Debug("Executed Successfully: SendRequest()")

	return

}
