package output

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/fatih/color"
	"github.com/valyala/fasthttp"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

func print(r *string, highlight *bool) {
	//fmt.Println("reached inside HIGHLIGHT")
	logging.Logger.Debug("Break: Reached printIt() start")

	if *highlight {
		lexer := lexers.Get("http")
		if lexer == nil {
			lexer = lexers.Fallback
		}

		err := quick.Highlight(
			color.Output,
			*r,
			lexer.Config().Name, // Lexer name (e.g., "json", "html")
			"terminal256",       // Formatter for CLI output
			"tokyonight-moon",   // Syntax highlighting style
		)
		if err != nil {
			logging.Logger.Warn("Colored output failed, printing response regularly. Error: %s", err)

		} else {
			//fmt.Printf("<-- BREAK -->")
			fmt.Println()
			return
		}
	}

	//fmt.Println(r) // printing it standard by default if highlight flag isn't true.
	os.Stdout.Write([]byte(*r))
	return
}

func OutputRequest(req *fasthttp.Request, args *models.CommandActions) error {
	reqStr := req.String() + "\n"
	print(&reqStr, &args.Color)
	//fmt.Printf("<-- BREAK -->")

	return nil
}

func OutputResponse(resp *models.Response, args *models.CommandActions) {
	//fmt.Println("reached outputresponse in print.go")

	// Should "waterfall" down, as users may want to print just status + body, or body + bytes, etc. Waterfalling allows them combos that eventually add back up to being a properly formatted response + bytes, duration, etc.
	if len(args.PrintOptions) > 0 {
		// If a print option is used, it essentially allows them to build their own response. If something isn't mentioned via the CLI, it won't be included in the printout when these options are used.
		if slices.Contains(args.PrintOptions, "nothing") {
			return
		}

		if slices.Contains(args.PrintOptions, "status") {
			statusMsg := strconv.Itoa(resp.StatusCode)

			statusMsg += " " + resp.Status + "\n"

			print(&statusMsg, &args.Color)
			//return
		}

		if slices.Contains(args.PrintOptions, "headers") {
			var headerStr string
			for key, value := range resp.Headers {
				headerStr += key + ": " + value
			}

			print(&headerStr, &args.Color)
		}

		if slices.Contains(args.PrintOptions, "body") {

			bodyStr := string(resp.Body)

			print(&bodyStr, &args.Color)
		}

		if slices.Contains(args.PrintOptions, "duration") {
			fmt.Printf("\nRequest duration: %d ms\n", resp.Duration)
		}

		if slices.Contains(args.PrintOptions, "bytes") {
			// TODO: For later. Probably will require reworking it from a slice to a map.

			fmt.Printf("\nRequest Body in Bytes:\n%d", resp.Body)
		}

		if slices.Contains(args.PrintOptions, "truncate") {
			// TODO: For later. Probably will require reworking it from a slice to a map.
			truncatedMsg := []byte("\n\033[36m[TRUNCATED]\033[0m\n\n")
			resp.Body = append(resp.Body[:120], truncatedMsg...)
		}
		//print(resp, &args.Highlight)
		//os.Stdout.Write()
		return

	}

	// Turn off truncate-by-default
	// Prev: if !args.More {
	// TODO: Will need to fix, to make it an optional flag.
	// WARNING: Need to modify, will break if response is too small. (AKA just a "hello" in the response body, for example.)
	// Moved to "truncate" print option above.
	/*
		if args.More {

			truncatedMsg := []byte("\n\033[36m[TRUNCATED]\033[0m\n\n")
			resp.Body = append(resp.Body[:120], truncatedMsg...)

		}
	*/

	//fmt.Println("reached stdout in print.go")

	// Default. Print full response.
	print(respToStr(resp), &args.Color)

	return
}

func respToStr(resp *models.Response) *string {
	var output string

	output += strconv.Itoa(resp.StatusCode) + " " + resp.Status + "\n"

	for key, value := range resp.Headers {
		output += fmt.Sprintf("%s: %s\n", key, value)
	}
	output += "\n"

	//output = resp.Headers + string(resp.Body)
	output += string(resp.Body)

	return &output
}
