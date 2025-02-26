/*
Copyright © 2025 alp1n3 1@alp1n3.dev
*/
package main

import (
	//"net/http"
	//"fmt"
	"os"
	//"fmt"
	"context"
	"log"
	"strings"
	"sync"

	//"bytes"
	//"io"

	//"runtime/pprof"
	//"runtime/trace"

	//"net/http"

	"github.com/urfave/cli/v3" // docs: https://cli.urfave.org/v3/examples/subcommands/
	"github.com/valyala/fasthttp"

	cmd "github.com/alp1n3-eth/cast/cmd/http"
	//"github.com/alp1n3-eth/cast/internal/http/parse"
	//"github.com/alp1n3-eth/cast/pkg/logging"
	//"github.com/alp1n3-eth/cast/parse"
)

func main() {
	//f, _ := os.Create("cpu.prof")
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	//s, _ := os.Create("trace.out")
	//trace.Start(s)
	//defer trace.Stop()

	//headers := &http.Header{}
	//bodyReader := &io.Reader
	//var body io.Reader
	//bodyReader := &body

	app := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "get",
				Aliases: []string{"GET", "post", "put", "delete", "patch", "options", "trace", "head", "connect"},
				Usage:   "send an HTTP request to a url.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "body",
						Value:   "",
						Usage:   "HTTP request body",
						Aliases: []string{"B"},
					},
					&cli.StringSliceFlag{
						Name:    "header",
						Usage:   "HTTP headers to include in the request",
						Aliases: []string{"H"},
					},
					&cli.BoolFlag{
						Name:    "debug",
						Usage:   "Enable debug output for the application",
						Aliases: []string{"D"},
					},
					&cli.BoolFlag{
						Name:    "highlight",
						Usage:   "Prettify the response body output with syntax highlighting",
						Aliases: []string{"HL"},
					},
					&cli.StringSliceFlag{
						Name:    "var",
						Usage:   "A text file to be iteratively used in any portion of the request to insert values.",
						Aliases: []string{"V"},
					},
					&cli.StringSliceFlag{
						Name:    "print",
						Usage:   "A text file to be iteratively used in any portion of the request to insert values.",
						Aliases: []string{"P"},
					},
					&cli.StringFlag{
						Name:    "file",
						Usage:   "A way to include a file in the request's body.",
						Aliases: []string{"F"},
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {
					//fmt.Println("added task: ", command.Args().First())
					//fmt.Println("Debug - All args:", os.Args)
					//fmt.Println("Debug - First arg:", os.Args[1])
					//fmt.Println("Debug - Context args:", command.Args().Slice())
					//fmt.Println("Debug - Body flag:", command.String("body"))

					// Can modify to make testing take longer. Send request multiple times. Currently hardcoded to send it

					//for i := 0; i <= 1; i++ {
					var debug bool
					var highlight bool
					var bodyStr string
					var uploadFilePath string

					printOption := command.StringSlice("print")

					replacementPair := make(map[string]string)
					headers := make(map[string]string)

					request := fasthttp.Request{}

					var wg sync.WaitGroup
					wg.Add(4)

					go func() {
						defer wg.Done()

						// Handle debug and highlight options
						debug = command.Bool("debug")
						highlight = command.Bool("highlight")
						uploadFilePath = command.String("file")
					}()

					go func() {
						defer wg.Done()

						// Handle custom body
						bodyStr = command.String("body")
						request.SetBody([]byte(bodyStr)) // TODO: What's going on? Prebuilding request here?
					}()

					go func() {
						defer wg.Done()

						// Handle custom headers
						headerSlice := command.StringSlice("header")
						//*headers = make(map[string]string)

						for _, h := range headerSlice {
							key, value, _ := strings.Cut(h, ":")
							//fmt.Print("reached headerslice main.go")
							if len(key) >= 1 {
								//key := strings.TrimSpace(parts[0])
								//value := strings.TrimSpace(parts[1])
								(headers)[key] = value
							}
						}
					}()

					go func() {
						defer wg.Done()

						// Handle replacement variables
						replacementSlice := command.StringSlice("var")
						//*replacementPair = make(map[string]string)

						for _, h := range replacementSlice {
							targetWord, value, _ := strings.Cut(h, "=")
							//fmt.Print("reached headerslice main.go")
							if len(targetWord) >= 1 {
								//key := strings.TrimSpace(parts[0])
								//value := strings.TrimSpace(parts[1])
								(replacementPair)[targetWord] = value
							}
						}

					}()

					/*
					   for _, h := range headerSlice {
					       parts := strings.SplitN(h, ":", 2)
					       //fmt.Print("reached headerslice main.go")
					       if len(parts) == 2 {
					           key := strings.TrimSpace(parts[0])
					           value := strings.TrimSpace(parts[1])
					           headers[key] = value
					       }
					   }
					*/

					/*
					   for _, h := range replacementSlice {
					   	//fmt.Print("reached replacementslice main.go")
					       parts := strings.SplitN(h, "=", 2)
					       if len(parts) == 2 {
					           key := strings.TrimSpace(parts[0])
					           value := strings.TrimSpace(parts[1])
					           replacementPair[key] = value
					       }
					   }
					*/

					wg.Wait()

					cmd.SendHTTP(os.Args[1], command.Args().First(), &bodyStr, &headers, &debug, &highlight, &replacementPair, &printOption, &uploadFilePath)

					//}
					// ^ Ending brace for profiling pprof

					// KEEP THIS
					//cmd.SendHTTP(os.Args[1], command.Args().First(), bodyStr, headers, debug, highlight, replacementPair)

					//f, _ := os.Create("mem.prof")
					//pprof.WriteHeapProfile(f)
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
