package cmd

import (
	"strings"
	//"fmt"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/valyala/fasthttp"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(method, urlVar string, body string, headers map[string]string, debug, highlight bool, replacementVariables map[string]string) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
    //apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))

    var err error
    var methodPtr string
    var urlVarPtr string
    methodPtr = strings.ToUpper(method)
    urlVarPtr = strings.ToLower(urlVar)

    printOption := "" // TODO: Placeholder currently, can be used to print response before request. Needs to have a flag created for it.

    if debug{
        logging.Init(true) // Activates debug mode.
    } else if !debug{
    	logging.Init(false)
    }

	if urlVar != "" {
		// Perform cli-based actions.

		result := &models.ExecutionResult{}

		result.Request.Req = parse.BuildRequest1(&methodPtr, &urlVarPtr, &body, &headers)
		logging.Logger.Debug("Executed Successfully: BuildRequest()")

		if len(replacementVariables) > 0 {
			//fmt.Println(replacementVariables)
			parse.SwapReqVals(result.Request.Req, &replacementVariables)
			logging.Logger.Debug("Executed Successfully: SwapReqVals()")
		}

		// TODO: Get sendhttprequqest working again
		err = executor.SendRequest(result, &debug, &highlight, &printOption)
		if err != nil {
			logging.Logger.Fatal("Error sending HTTP request")
		}
		logging.Logger.Debug("Executed Successfully: SendRequest")

			// TODO: Get printout of response working again
			// TODO: Get flags tied-in in order to provide body.
		return

		} else {
			logging.Logger.Fatal("Invalid command provided")
			return
		}
}
