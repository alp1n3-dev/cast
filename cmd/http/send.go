package cmd

import (
	"strings"
	"slices"
	//"fmt"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/valyala/fasthttp"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(method, urlVar string, body *string, headers *map[string]string, debug, highlight *bool, replacementVariables *map[string]string, printOption *[]string) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
    //apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))

    var err error

    method = strings.ToUpper(method)
    urlVar = strings.ToLower(urlVar)

    //printOption := "" // TODO: Placeholder currently, can be used to print response before request. Needs to have a flag created for it.

    if *debug{
        logging.Init(true) // Activates debug mode.
    } else if !*debug{
    	logging.Init(false)
    }

    logging.Logger.Debugf("Debug: %t, Method: %s, URI: %s", debug, method, urlVar)

	if urlVar != "" {
		// Perform cli-based actions.

		result := &models.ExecutionResult{}

		result.Request.Req = parse.BuildRequest(&method, &urlVar, body, headers)
		logging.Logger.Debugf("BuildRequest: %s", result.Request.Req)


		if len(*replacementVariables) > 0 {
			//fmt.Println(replacementVariables)
			logging.Logger.Debugf("Replacement Variables: %s", replacementVariables)
			parse.SwapReqVals(result.Request.Req, replacementVariables)
			logging.Logger.Debug("Executed Successfully: SwapReqVals()")
		}

		logging.Logger.Debugf("Request being sent: %s", result.Request.Req)
		// Needs to be the one directly before sending it, as changes may happen in functions like SwapReqVals().
		//if *printOption == "request" {
     		//output.PrintHTTP(result.Request.Req, nil, highlight)
      	//}
     	if len(*printOption) > 0 {
      		if slices.Contains(*printOption, "request") {
        		output.PrintHTTP(result.Request.Req, nil, highlight, printOption)
        	}
      	}



		// TODO: Get sendhttprequqest working again
		err = executor.SendRequest(result, debug, highlight, printOption)
		if err != nil {
			logging.Logger.Debugf("Result: %x, Highlight: %t, Print Option: %s", *result, *highlight, *printOption)
			logging.Logger.Fatal("Error sending HTTP request")
		}

		logging.Logger.Debug("Executed Successfully: SendRequest()")

			// TODO: Get printout of response working again
			// TODO: Get flags tied-in in order to provide body.
		return

		} else {
			logging.Logger.Fatal("Invalid command provided")
			return
		}
}
