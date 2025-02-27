package cmd

import (
	"io"
	"os"
	"slices"

	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"
	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/apperrors"
)

func SendHTTP(headers, replacementVariables *map[string]string, CLIArgs *models.Request) {

	// TODO: Fix panic caused by apperrors.HandleExecutionError
	//apperrors.HandleExecutionError(
	//apperrors.Wrap(apperrors.ErrInvalidHeaderFormat, "random-header"))

	var err error

	if CLIArgs.CLI.Debug {
		logging.Init(true) // Activates debug mode.
	} else {
		logging.Init(false)
	}

	logging.Logger.Debugf("Debug: %t, Method: %s, URI: %s", CLIArgs.CLI.Debug, CLIArgs.CLI.Method, CLIArgs.CLI.URL)

	if CLIArgs.CLI.URL != "" {
		// Perform cli-based actions.

		result := &models.ExecutionResult{}

		if CLIArgs.CLI.FileUploadPath != "" {
			logging.Logger.Debugf("Upload file path: %s", CLIArgs.CLI.FileUploadPath)

			CLIArgs.CLI.Body = readFileIntoBody(&CLIArgs.CLI.FileUploadPath)
		}

		result.Request.Req = parse.BuildRequest(&CLIArgs.CLI.Method, &CLIArgs.CLI.URL, &CLIArgs.CLI.Body, headers)
		logging.Logger.Debugf("BuildRequest: %s", result.Request.Req)

		if len(*replacementVariables) > 0 {
			logging.Logger.Debugf("Replacement Variables: %s", replacementVariables)
			parse.SwapReqVals(result.Request.Req, replacementVariables)
			logging.Logger.Debug("Executed Successfully: SwapReqVals()")
		}

		logging.Logger.Debugf("Request being sent: %s", result.Request.Req)
		// Needs to be the one directly before sending it, as changes may happen in functions like SwapReqVals().

		if len(CLIArgs.CLI.PrintOptions) > 0 {
			if slices.Contains(CLIArgs.CLI.PrintOptions, "request") {
				output.PrintHTTP(result.Request.Req, nil, &CLIArgs.CLI.Highlight, &CLIArgs.CLI.PrintOptions)
			}
		}

		err = executor.SendRequest(result, &CLIArgs.CLI.Debug, &CLIArgs.CLI.Highlight, &CLIArgs.CLI.PrintOptions, &CLIArgs.CLI.RedirectsToFollow)
		if err != nil {
			//logging.Logger.Debugf("Highlight: %t, Print Option: %s", &CLIArgs.CLI.Highlight, &CLIArgs.CLI.PrintOptions)
			logging.Logger.Fatal("Error sending HTTP request")
		}

		logging.Logger.Debug("Executed Successfully: SendRequest()")

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
		logging.Logger.Fatal("Error opening file")
	}
	defer file.Close()

	fileContents, err = io.ReadAll(file)
	if err != nil {
		logging.Logger.Fatal("Error reading file")
	}

	logging.Logger.Debugf("fileContents: %x", fileContents)
	logging.Logger.Debug("Successfully read file contents for file upload.")

	return fileContents
}
