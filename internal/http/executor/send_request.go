package executor

import (
	"fmt"
	"io"
	"os"
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
	response := &models.Response{}
	// Going to be a flag later, based on if asserts are detected in the file when read.
	assertsRequired := false
	var err error

	// TODO: Add the ability to track byte sizes of responses.

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	startTime := time.Now()

	if request.CLI.RedirectsToFollow > 0 {
		err = fasthttp.DoRedirects(request.Req, resp, request.CLI.RedirectsToFollow)
		logging.Logger.Info("Followed redirect\n")
	} else {
		err = fasthttp.Do(request.Req, resp)
	}
	if err != nil {
		//fmt.Printf("Client get failed: %s\n", err)
		logging.Logger.Error("Client get failed", "err", err)
		return err
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	response.Duration = duration

	// TODO: Implement no-response for printOption flag
	//if *printOption == "no-response" {
	//return nil
	//}

	output.PrintHTTP(nil, resp, &request.CLI.Highlight, &request.CLI.PrintOptions)

	printOptions := "duration"

	if strings.Contains(printOptions, "duration") {
		fmt.Printf("\nRequest duration: %d ms\n", duration.Milliseconds())
	}

	if request.CLI.DownloadPath != "" {
		err = saveResponseToFile(resp, request.CLI.DownloadPath)
		if err != nil {
			logging.Logger.Error("Error saving response to file", "err", err)
			return err
		}
		logging.Logger.Info("Response saved to file", "path", request.CLI.DownloadPath)
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
