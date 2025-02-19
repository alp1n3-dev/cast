package output

import (
	"fmt"
	"strconv"
	//"mime"
	//"strings"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"

	"github.com/fatih/color"
	//"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
)

// PrintResponse will format and highlight the response
func PrintResponse(r models.ExecutionResult) {
	var completeResponse string
	coloredOutput := false
	var err error
	err = nil
	// Print the status code in green for success and red for failure
	respStatusInt := r.Response.StatusCode

	//logging.Logger.Debug(respStatusInt)
	//respStatusInt, err := strconv.Atoi(respStatusChar)
	//if err != nil {
		//fmt.Println("Error converting status code from char to int.")
		//}

	//color.Set(color.FgHiBlue)
	//fmt.Printf("%s ", r.Response.Protocol)
	completeResponse += r.Response.Protocol + " "


	if respStatusInt >= 200 && respStatusInt < 300 {
		//color.Set(color.FgGreen)
	} else {
		//color.Set(color.FgRed)
	}
	//defer color.Unset()

	// Print the Status Code with status
	//fmt.Printf("%d %s\n", r.Response.StatusCode, r.Response.Status)
	completeResponse += strconv.Itoa(r.Response.StatusCode) + " "
	completeResponse += r.Response.Status + "\n"

	// Print headers (key-value)
	//for key, value := range r.Response.Headers {
		//fmt.Printf("%s: %s\n", key, value)
		//}

	//color.Set(color.FgBlue)
	for name, values := range r.Response.Headers {
    // Loop over all values for the name.
    	for _, value := range values {
        	//fmt.Printf("%s: %s\n", name, value)
         	//color.Set(color.FgHiMagenta)
          	//fmt.Printf("%s: ", name)
           	completeResponse += name + ": "

           	//color.Unset()
           	//fmt.Printf("%s\n", value)
            completeResponse += value + "\n"
     	}
	}

	completeResponse += "\n"

	// Read and print the body (if any)
	//body, err := io.ReadAll(r.Response.Body)
	//if err != nil {
		//return fmt.Errorf("failed to read response body: %w", err)
		//}

	// Restore the body so that it can be used again after the return.
	//resp.Body = io.NopCloser(bytes.NewReader(body))

	// Print the body in blue
	//color.Set(color.FgCyan)
	//fmt.Println("\n" + string(r.Response.Body) + "\n")
	//color.Unset() // Unset the color formatting

	//fmt.Println(r.Response.ContentType)
	completeResponse += string(r.Response.Body) + "\n"

	//mediaType, _, _ := mime.ParseMediaType(r.Response.ContentType)
	//mediaType = strings.Split(mediaType, "/")[1]
	//lexer := lexers.Get(mediaType)

	if coloredOutput{ // if coloredOutput flag is specified. Can be changed to true by default in configs.
		err = printResponseColor(&completeResponse)

		if err != nil {
			// If highlighting fails, print the raw response
			logging.Logger.Warn("Colored output failed, printing response regularly. Error: %s", err)
			fmt.Println(completeResponse)
		}
		return
	}

	fmt.Println(completeResponse)

	return
}

func printResponseColor(completeResponse *string) error {
	lexer := lexers.Get("http")
	if lexer == nil {
		// Default to plain text if no lexer is found
		lexer = lexers.Fallback
	}

	//fmt.Println(mediaType)

	// Use Chroma to highlight output
	err := quick.Highlight(
		color.Output, // Use color-capable output
		//string(r.Response.Body),
		*completeResponse,
		lexer.Config().Name, // Lexer name (e.g., "json", "html")
		"terminal256",          // Formatter for CLI output
		"tokyonight-moon",           // Syntax highlighting style
	)
	if err != nil {
		return err
	}

	return nil
}
