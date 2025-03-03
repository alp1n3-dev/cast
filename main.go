/*
Copyright © 2025 alp1n3 1@alp1n3.dev
*/
package main

import (
	"context"
	"log"
	"os"

	"github.com/alp1n3-eth/cast/cmd"
)

/*
"context"
"fmt"
"log"
"maps"
"os"
"strings"

"github.com/urfave/cli/v3" // docs: https://cli.urfave.org/v3/examples/subcommands/
"github.com/valyala/fasthttp"

cmd "github.com/alp1n3-eth/cast/cmd/http"
"github.com/alp1n3-eth/cast/internal/env"
"github.com/alp1n3-eth/cast/pkg/models"
*/
func main() {
	if err := cmd.Execute(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

/*
func main() {

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
						Name:  "debug",
						Usage: "Enable debug output for the application",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "highlight",
						Usage: "Prettify the response body output with syntax highlighting",
						Value: false,
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
					&cli.IntFlag{
						Name:    "redirect",
						Usage:   "A way to follow redirects up to < INT >.",
						Aliases: []string{"R"},
					},
					&cli.StringFlag{
						Name:    "download",
						Usage:   "Path to save the response body to a file.",
						Aliases: []string{"D"},
					},
					&cli.BoolFlag{
						Name:  "read-encrypted",
						Usage: "Read an encrypted KV store using a password.",
					},
					&cli.BoolFlag{
						Name:  "curl",
						Usage: "Output the request as a curl command instead of sending it.",
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {

					// Can modify to make testing take longer. Send request multiple times. Currently hardcoded to send it

					//userInputs := &models.Request{Req: fasthttp.AcquireRequest()}
					//defer fasthttp.ReleaseRequest(userInputs.Req)
					req := fasthttp.AcquireRequest()
					defer fasthttp.ReleaseRequest(req)

					userInputs := &models.HTTPRequestContext{
						Request: models.Request{
							Req: req,
						},

						CmdArgs: models.CommandActions{
							PrintOptions:      command.StringSlice("print"),
							RedirectsToFollow: int(command.Int("redirect")),
							Debug:             command.Bool("debug"),
							Highlight:         command.Bool("highlight"),
							FileUploadPath:    command.String("fileupload"),
							DownloadPath:      command.String("download"),
							CurlOutput:        command.Bool("curl"),
						},
					}

					userInputs.Request.Req.Header.SetMethod(strings.ToUpper(os.Args[1]))

					userInputs.Request.Req.SetRequestURI(command.Args().First())

					if command.Bool("read-encrypted") {
						fmt.Print("Enter password: ")
						password, err := env.RetrievePasswordFromUser()
						if err != nil {
							fmt.Println("error retrieving password")
						}
						fmt.Println(password)
					}

					//userInputs.CLI.PrintOptions = command.StringSlice("print")
					//userInputs.CLI.RedirectsToFollow = int(command.Int("redirect"))

					//userInputs.CLI.Debug = command.Bool("debug")

					//serInputs.CLI.Highlight = command.Bool("highlight")

					//userInputs.CLI.FileUploadPath = command.String("fileupload")

					//userInputs.CLI.DownloadPath = command.String("download")

					// Handle custom body

					userInputs.Request.Req.SetBodyRaw([]byte(command.String("body")))

					// Handle custom headers
					//headerSlice := command.StringSlice("header")
					//*headers = make(map[string]string)

					for _, h := range command.StringSlice("header") {
						key, value, _ := strings.Cut(h, ":")

						if len(key) >= 1 {
							//(headers)[key] = value
							//headers[key] = []byte(value)

							userInputs.Request.Req.Header.Set(key, value)
						}
					}
					if userInputs.Request.Req.Header.Peek("Content-Type") == nil {
						userInputs.Request.Req.Header.Add("Content-Type", "text/html")
					}



					replacementPair := parseReplacementValues(command.StringSlice("var"))

					//fmt.Println(replacementPair)

					cmd.SendHTTP(replacementPair, userInputs)

					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseReplacementValues(replacementSlice []string) *map[string]string {
	replacementPair := make(map[string]string)

	for _, h := range replacementSlice {
		if strings.Contains(h, ".env") {
			kvFileMap, _ := env.ReadKVFile(h)

			maps.Copy(replacementPair, *kvFileMap)
		} else {
			targetWord, value, _ := strings.Cut(h, "=")

			if len(targetWord) >= 1 {

				replacementPair[targetWord] = value
			}
		}
	}
	return &replacementPair
}

*/
