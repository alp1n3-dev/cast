package output

import (
	"fmt"
	//"io"

	"github.com/fatih/color"
	"github.com/valyala/fasthttp"

	//"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"

	//"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/alp1n3-eth/cast/pkg/models"
)

// PrintResponse will format and highlight the response
func PrintHTTP(req *fasthttp.Request, resp *fasthttp.Response, highlight *bool) {
		var r string

		if req == nil {
			r = resp.String()
		} else if resp == nil {
			r = req.String()
		}

		if *highlight {
			lexer := lexers.Get("http")
    		if lexer == nil {
     		lexer = lexers.Fallback
     		}

      		err := quick.Highlight(
       			color.Output,
         		r,
          		lexer.Config().Name, // Lexer name (e.g., "json", "html")
           		"terminal256",          // Formatter for CLI output
             	"tokyonight-moon",           // Syntax highlighting style
       		)
       		if err != nil {
       			logging.Logger.Warn("Colored output failed, printing response regularly. Error: %s", err)

         	} else {
          		fmt.Printf("\n")
            	if resp == nil {
             		fmt.Printf("\n")
             	}
          		return
          }

		}

	fmt.Println(r) // printing it standard by default if highlight flag isn't passed.

    return
}
