package cmd

import (
	"fmt"
	"net/http"

	"strings"
	//"fmt"
	"bytes"
	"io"

	//"fmt"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	"github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
	//"github.com/alp1n3-eth/cast/models"
)

func SendHTTP(methodVar, urlVar, requestBody string) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
    //apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))

    //fmt.Println(debugMode)

    //if debugMode {
        logging.Init(true) // Debug mode is TRUE
    //}


		methodVar = strings.ToUpper(methodVar)
		logging.Logger.Debug(requestBody)
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

			method := methodVar
			url := strings.ToLower(urlVar)

			var headers http.Header
			var body io.Reader
			var err error

			if requestBody != "" {
				body = bytes.NewBufferString(requestBody)
			}

			var result models.ExecutionResult

			result.Request = parse.BuildRequest(method, url, body, headers)

			// TODO: Get sendhttprequqest working again
			logging.Logger.Debug("Request headers: ")
			//for k, v := range result.Request.Headers {
        		//fmt.Printf("Header field %q, Value %q\n", k, v)
        		//} // TODO: will panic if no headers provided
			result, err = executor.SendFastHTTPRequest(result)
			if err != nil {
				logging.Logger.Fatal("Error sending HTTP request")
			}

			// TODO: Get printout of response working again
			fmt.Println(result.Response.Status)
			fmt.Println(result.Response.Headers)
			fmt.Println(result.Response.Body)

			output.PrintResponse(result)
			// TODO: Get flags tied-in in order to provide body.

		} else {
			logging.Logger.Fatal("Invalid command provided")
		}
}
