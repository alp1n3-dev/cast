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
			fmt.Println()
			return
		}
	}

	//fmt.Println(r) // printing it standard by default if highlight flag isn't true.
	os.Stdout.Write([]byte(*r))

}

func OutputRequest(req *fasthttp.Request, args *models.CommandActions) error {
	reqStr := req.String() + "\n"
	print(&reqStr, &args.Color)

	return nil
}

func OutputResponse(resp *models.Response, args *models.CommandActions) {
	//fmt.Println("reached outputresponse in print.go")

	// Should "waterfall" down, as users may want to print just status + body, or body + bytes, etc. Waterfalling allows them combos that eventually add back up to being a properly formatted response + bytes, duration, etc.
	if len(args.PrintOptions) > 0 {
		if slices.Contains(args.PrintOptions, "status") {
			statusMsg := strconv.Itoa(resp.StatusCode)

			statusMsg += " " + resp.Status + "\n"

			print(&statusMsg, &args.Color)
			//return
		}

		if slices.Contains(args.PrintOptions, "header-only") {
			// TODO: For later. Probably will require reworking it from a slice to a map.

		}

		if slices.Contains(args.PrintOptions, "body") {
			// TODO: For later. Probably will require reworking it from a slice to a map.

		}

		if slices.Contains(args.PrintOptions, "no-response") {
			// TODO: For later. Probably will require reworking it from a slice to a map.

		}

		if slices.Contains(args.PrintOptions, "duration") {
			fmt.Printf("\nRequest duration: %d ms\n", resp.Duration)
		}

		if slices.Contains(args.PrintOptions, "bytes") {
			// TODO: For later. Probably will require reworking it from a slice to a map.

		}
		//print(resp, &args.Highlight)
		//os.Stdout.Write()

	}

	if !args.More {
		//resp.Body = resp.Body[:100]
		truncatedMsg := []byte("\n\033[36m[TRUNCATED]\033[0m\n")
		resp.Body = append(resp.Body[:120], truncatedMsg...)
		//resp.Body = resp.Body[:100] + []byte("\n[TRUNCATED]")
	}

	//fmt.Println("reached stdout in print.go")
	print(respToStr(resp), &args.Color)

	return
}

func respToStr(resp *models.Response) *string {
	var output string

	output = resp.Headers + string(resp.Body)

	return &output
}
