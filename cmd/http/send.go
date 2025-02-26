package cmd

import (
	"io"
	"os"
	"slices"
	"strings"

	//"fmt"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/valyala/fasthttp"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(method, urlVar string, body *string, headers *map[string]string, debug, highlight *bool, replacementVariables *map[string]string, printOption *[]string, uploadFilePath *string) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
	//apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))

	var err error
	//var bodyByte *[]byte
	var bodyByte []byte

	method = strings.ToUpper(method)
	urlVar = strings.ToLower(urlVar)

	//printOption := "" // TODO: Placeholder currently, can be used to print response before request. Needs to have a flag created for it.

	// TODO: Put flag in main.go and ensure a valid value is passed to this function.
	//uploadFilePath := "tests/wordlists/random_endpoints.txt"

	if *debug {
		logging.Init(true) // Activates debug mode.
	} else if !*debug {
		logging.Init(false)
	}

	logging.Logger.Debugf("Debug: %t, Method: %s, URI: %s", *debug, method, urlVar)

	if urlVar != "" {
		// Perform cli-based actions.

		result := &models.ExecutionResult{}

		if *uploadFilePath != "" {
			logging.Logger.Debugf("Upload file path: %s", *uploadFilePath)

			bodyByte = readFileIntoBody(uploadFilePath)

		} else {
			bodyByte = []byte(*body)
		}

		result.Request.Req = parse.BuildRequest(&method, &urlVar, &bodyByte, headers)
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

func readFileIntoBody(uploadFilePath *string) []byte {
	var fileContents []byte

	// Check if file exists
	if _, err := os.Stat(*uploadFilePath); os.IsNotExist(err) {
		logging.Logger.Fatalf("File does not exist: %s", *uploadFilePath)
	}

	file, err := os.Open(*uploadFilePath)
	if err != nil {
		logging.Logger.Fatalf("Error opening file: %w", err)
	}
	defer file.Close()

	fileContents, err = io.ReadAll(file)
	if err != nil {
		logging.Logger.Fatalf("Error reading file: %w", err)
	}

	logging.Logger.Debugf("fileContents: %x", fileContents)
	logging.Logger.Debug("Successfully read file contents for file upload.")

	return fileContents
}
