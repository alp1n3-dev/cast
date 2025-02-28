package executor

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/internal/http/parse"
	output "github.com/alp1n3-eth/cast/output/http"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

// Should assume all fields have been created and validated by the time they get here.
func SendRequest(request *models.Request) error {
	// TODO: Isolate unrelated logic to its own files. Find a better / more efficient way to "keep" / store the resp that isn't "expensive", so that printing doesn't need to be called from within this SendRequest().

	response := &models.Response{}
	// Going to be a flag later, based on if asserts are detected in the file when read.
	assertsRequired := false
	var err error

	// TODO: Add the ability to track byte sizes of responses.

	//req := fasthttp.AcquireRequest()
	//defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	startTime := time.Now()

	if request.CLI.RedirectsToFollow > 0 {
		if err = fasthttp.DoRedirects(request.Req, resp, request.CLI.RedirectsToFollow); err != nil {
			return err
		}
		logging.Logger.Info("Success: Followed redirect")
	} else {
		if err = fasthttp.Do(request.Req, resp); err != nil {
			return err
		}
		logging.Logger.Debug("Success: Request sent")
	}

	duration := time.Since(startTime).Milliseconds()
	response.Duration = int(duration)

	// TODO: Implement no-response for printOption flag
	//if *printOption == "no-response" {
	//return nil
	//}
	if request.CLI.PrintOptions != nil {
		if slices.Contains(request.CLI.PrintOptions, "response") {
			output.PrintHTTP(nil, resp, &request.CLI.Highlight, &request.CLI.PrintOptions)
		} else if slices.Contains(request.CLI.PrintOptions, "body") {
			// TODO:
		} else if slices.Contains(request.CLI.PrintOptions, "status") {
			// TODO:
		}

		if slices.Contains(request.CLI.PrintOptions, "duration") {
			fmt.Printf("\nRequest duration: %d ms\n", duration)
		}
	}
	output.PrintHTTP(nil, resp, &request.CLI.Highlight, &request.CLI.PrintOptions)

	/*if request.CLI.DownloadPath != "" {
	if err = saveResponseToFile(resp, request.CLI.DownloadPath); err != nil {
		logging.Logger.Error("Error saving response to file", "err", err)
		return err
	}
	logging.Logger.Info("Success: Response saved to file", "path", request.CLI.DownloadPath)
	}*/
	if request.CLI.DownloadPath != "" {
		if err := os.WriteFile(request.CLI.DownloadPath, resp.Body(), 0644); err != nil {
			return err
		}
	}

	// Blocking off the below section for later with a return that'll be hit
	if !assertsRequired {
		return nil
	}

	var assertion string // Placeholder for if the response values need to be saved and filtered
	if len(assertion) > 0 {
		*response = parse.BuildFastHTTPResponse(resp)
	}

	return nil
}

func saveResponseToFile(resp *fasthttp.Response, outputPath string) error {
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
