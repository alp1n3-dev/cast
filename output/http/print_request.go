package output

import (
	"fmt"
	//"io"

	"github.com/valyala/fasthttp"
	"github.com/fatih/color"
	//"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"

	//"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"
)

// PrintResponse will format and highlight the response
func PrintRequest(r *fasthttp.Request, highlight *bool) {
	logging.Logger.Debug("Reached PrintRequest() 1")
	/*
	var completeRequest string
	request := r.Request
	//coloredOutput := false
	var err error
	err = nil

	logging.Logger.Debug("Reached PrintRequest() 2")

	// Parse method & URL
	completeRequest = request.Method + " "
	completeRequest += request.URL.String() + "\n"

	logging.Logger.Debug("Reached PrintRequest() 3")

	// Parse headers
	for name, values := range request.Headers {
    	for _, value := range values {
           	completeRequest += name + ": "
            completeRequest += value + "\n"
     	}
	}

	// Add break between headers and body (if it exists)
	completeRequest += "\n"

	logging.Logger.Debug("Reached PrintRequest() 4")

	if request.Body != nil {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			logging.Logger.Error(err)
		}
		completeRequest += string(body)
	}
	*/


	logging.Logger.Debug("Reached PrintRequest() 5")

	//fmt.Println(request)

	//mediaType, _, _ := mime.ParseMediaType(r.Response.ContentType)
	//mediaType = strings.Split(mediaType, "/")[1]
	//lexer := lexers.Get(mediaType)

	if *highlight{ // if coloredOutput flag is specified. Can be changed to true by default in configs.
		err := printFastHTTPColor(r)

		if err != nil {
			// If highlighting fails, print the raw response
			logging.Logger.Warn("Colored output failed, printing response regularly. Error: %s", err)
			fmt.Println(r)
		}
		return
	}

	fmt.Println(r)

	return
}

func printFastHTTPColor(r *fasthttp.Request) error {
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
		r.String(),
		lexer.Config().Name, // Lexer name (e.g., "json", "html")
		"terminal256",          // Formatter for CLI output
		"tokyonight-moon",           // Syntax highlighting style
	)
	if err != nil {
		return err
	}

	return nil
}
