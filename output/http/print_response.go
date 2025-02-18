package output

import (
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/models"
	//"github.com/alp1n3-eth/cast/pkg/logging"

	"github.com/fatih/color"
)

// PrintResponse will format and highlight the response
func PrintResponse(r models.ExecutionResult) error {
	// Print the status code in green for success and red for failure
	respStatusInt := r.Response.StatusCode

	//logging.Logger.Debug(respStatusInt)
	//respStatusInt, err := strconv.Atoi(respStatusChar)
	//if err != nil {
		//fmt.Println("Error converting status code from char to int.")
		//}
	if respStatusInt >= 200 && respStatusInt < 300 {
		color.Set(color.FgGreen)
	} else {
		color.Set(color.FgRed)
	}
	defer color.Unset()

	// Print the Status Code with status
	fmt.Printf("%d %s\n", r.Response.StatusCode, r.Response.Status)

	// Print headers (key-value)
	//for key, value := range r.Response.Headers {
		//fmt.Printf("%s: %s\n", key, value)
		//}
	color.Set(color.FgBlue)
	for name, values := range r.Response.Headers {
    // Loop over all values for the name.
    	for _, value := range values {
        	fmt.Printf("%s: %s\n", name, value)
     	}
	}

	// Read and print the body (if any)
	//body, err := io.ReadAll(r.Response.Body)
	//if err != nil {
		//return fmt.Errorf("failed to read response body: %w", err)
		//}

	// Restore the body so that it can be used again after the return.
	//resp.Body = io.NopCloser(bytes.NewReader(body))

	// Print the body in blue
	color.Set(color.FgCyan)
	fmt.Println("\n" + string(r.Response.Body) + "\n")
	color.Unset() // Unset the color formatting

	return nil
}
