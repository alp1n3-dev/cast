package output

import (
	"fmt"
	"strconv"

	"github.com/alp1n3-eth/cast/models"
	"github.com/fatih/color"
)

// PrintResponse will format and highlight the response
func PrintResponse(r models.HTTPRequest) error {
	// Print the status code in green for success and red for failure
	respStatusChar := r.Response.StatusCode[0:3]
	respStatusInt, err := strconv.Atoi(respStatusChar)
	if err != nil {
		fmt.Println("Error converting status code from char to int.")
	}
	if respStatusInt >= 200 && respStatusInt < 300 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgRed)
	}
	defer color.Unset()

	// Print the Status Code with status
	fmt.Printf("%s\n", r.Response.StatusCode)

	// Print headers (key-value)
	for key, value := range r.Response.Headers {

		fmt.Printf("%s: %s\n", key, value)

	}

	// Read and print the body (if any)
	//body, err := io.ReadAll(r.Response.Body)
	//if err != nil {
		//return fmt.Errorf("failed to read response body: %w", err)
		//}

	// Restore the body so that it can be used again after the return.
	//resp.Body = io.NopCloser(bytes.NewReader(body))

	// Print the body in blue
	color.Set(color.FgBlue)
	fmt.Println("\n" + r.Response.Body)
	color.Unset() // Unset the color formatting

	return nil
}
