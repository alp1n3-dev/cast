package cmd

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"bytes"
)

func PrintResponse(resp *http.Response) (*http.Response, error) {

	// Print the major and minor protocol, and response status.
	fmt.Printf("HTTP/%d.%d %s\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	fmt.Println()

		// Read the response body.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			os.Exit(1)
		}

		// Print the body of the response.
		fmt.Println(string(body))

		// Restore the body so that it can be used again after the return.
		resp.Body = io.NopCloser(bytes.NewReader(body))

		return resp, nil
}
