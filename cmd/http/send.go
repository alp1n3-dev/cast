package cmd

import (
	"fmt"
	//"net/http"

	"strings"
	//"fmt"
	//"bytes"
	//"io"

	//"fmt"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	//"github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
	//"github.com/alp1n3-eth/cast/models"
)

func SendHTTP(method, urlVar string, body string, headers map[string]string, debug, highlight bool) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
    //apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))


    //fmt.Println(debugMode)
    var methodPtr string
    var urlVarPtr string

    fmt.Println(headers)

    // Placeholder value
    printOption := ""

    if debug == true {
        logging.Init(true) // Debug mode is TRUE
    } else if debug != true {
    	logging.Init(false)
    }


		methodPtr = strings.ToUpper(method)
		//logging.Logger.Debug(requestBody)
		//validMethod := parse.ValidateMethod(methodVar)
		/*
		if !validMethod {
			// Perform file-based actions.
			if !parse.ValidateFile(methodVar) {
				logging.Logger.Fatal("No file extension detected")
			}
	 	*/
		//if !validMethod {

		//}
		if urlVar != "" {
			// Perform cli-based actions.
			//logging.Logger.Debug("Cli-Based Route")

			//method := methodVar
			urlVarPtr = strings.ToLower(urlVar)

			//var headers http.Header
			//var body io.Reader
			var err error
			//var result *models.ExecutionResult
			result := &models.ExecutionResult{}

			//if requestBody != "" {
				//body = bytes.NewBufferString(requestBody)
			//}

			//result.Request = parse.BuildRequest(&methodPtr, &urlVarPtr, body, headers)
			result.Request.Req = parse.BuildRequest(&methodPtr, &urlVarPtr, &body, &headers)
			logging.Logger.Debug("Executed Successfully: BuildRequest()")



			// TODO: Get sendhttprequqest working again
			err = executor.SendRequest(result, &debug, &highlight, &printOption)
			if err != nil {
				logging.Logger.Fatal("Error sending HTTP request")
			}



			// TODO: Get printout of response working again
			//fmt.Println(result.Response.Status)
			//fmt.Println(result.Response.Headers)
			//fmt.Println(result.Response.Body)

			//output.PrintResponse(result, highlight)
			//output.PrintHTTP(nil, result.Response, highlight)
			// TODO: Get flags tied-in in order to provide body.
			return

		} else {
			logging.Logger.Fatal("Invalid command provided")
			return
		}
}
