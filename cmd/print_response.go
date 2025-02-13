package cmd

import (
	"fmt"
	"net/http"
	"io"
	"bytes"

	"github.com/fatih/color"

)

// PrintResponse will format and highlight the response
func PrintResponse(resp *http.Response) (*http.Response, error) {
	// Print the status code in green for success and red for failure
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgRed)
	}
	defer color.Unset()

	// Print the Status Code with status
	fmt.Printf("%s\n", resp.Status)

	// Print headers (key-value)
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Read and print the body (if any)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Restore the body so that it can be used again after the return.
	resp.Body = io.NopCloser(bytes.NewReader(body))

	// Print the body in blue
	color.Set(color.FgBlue)
	fmt.Println("\n" + string(body))
	color.Unset() // Unset the color formatting

	return resp, nil
}
