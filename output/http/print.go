package output

import (
	//"fmt"
	"strconv"
	"os"
	"slices"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/pkg/logging"
)

// PrintResponse will format and highlight the response
func PrintHTTP(req *fasthttp.Request, resp *fasthttp.Response, highlight *bool, printOption *[]string) {
		var r string

		if req == nil {
			r = resp.String()
		} else if resp == nil {
			r = req.String()
		} else {
			logging.Logger.Debug("Error occurred in PrintHTTP() attempting to set req/resp")
		}

		//fmt.Println(*printOption)
		if len(*printOption) > 0 {
			if slices.Contains(*printOption, "status") {
				statusMsg := strconv.Itoa(resp.Header.StatusCode())
				statusMsg += " " + string(resp.Header.StatusMessage())
				printIt(&statusMsg, highlight)
				return
			}

			if slices.Contains(*printOption, "header-only") {
				// TODO: For later. Probably will require reworking it from a slice to a map.

			}
		}

	printIt(&r, highlight)


	logging.Logger.Debug("Resp/Req should have successfully printed")

    return
}

func printIt (r *string, highlight *bool) {
	//fmt.Println("reached inside HIGHLIGHT")
	if *highlight {
		lexer := lexers.Get("http")
    		if lexer == nil {
     		lexer = lexers.Fallback
     		}

      		err := quick.Highlight(
       			color.Output,
         		*r,
          		lexer.Config().Name, // Lexer name (e.g., "json", "html")
           		"terminal256",          // Formatter for CLI output
             	"tokyonight-moon",           // Syntax highlighting style
       		)
       		if err != nil {
       			logging.Logger.Warn("Colored output failed, printing response regularly. Error: %s", err)

         	} else {
          		return
          }
	}


          //fmt.Println(r) // printing it standard by default if highlight flag isn't passed.
	os.Stdout.Write([]byte(*r))
}
