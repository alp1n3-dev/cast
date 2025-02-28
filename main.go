/*
Copyright © 2025 alp1n3 1@alp1n3.dev
*/
package main

import (
	"os"

	"context"
	"log"
	"strings"

	"github.com/urfave/cli/v3" // docs: https://cli.urfave.org/v3/examples/subcommands/
	"github.com/valyala/fasthttp"

	cmd "github.com/alp1n3-eth/cast/cmd/http"
	"github.com/alp1n3-eth/cast/pkg/models"
)

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
						Name:    "highlight",
						Usage:   "Prettify the response body output with syntax highlighting",
						Value:   false,
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
						Aliases: []string{"FU"},
					},
					&cli.IntFlag{
						Name:    "redirect",
						Usage:   "A way to follow redirects up to < INT >.",
						Aliases: []string{"RD"},
					},
					&cli.StringFlag{
						Name:    "download",
						Usage:   "Path to save the response body to a file.",
						Aliases: []string{"DL"},
					},
				},
				Action: func(ctx context.Context, command *cli.Command) error {

					// Can modify to make testing take longer. Send request multiple times. Currently hardcoded to send it

					userInputs := &models.Request{Req: fasthttp.AcquireRequest()}
					defer fasthttp.ReleaseRequest(userInputs.Req)

					userInputs.Req.Header.SetMethod(strings.ToUpper(os.Args[1]))

					userInputs.Req.SetRequestURI(command.Args().First())

					userInputs.CLI.PrintOptions = command.StringSlice("print")
					userInputs.CLI.RedirectsToFollow = int(command.Int("redirect"))

					replacementPair := make(map[string]string)

					userInputs.CLI.Debug = command.Bool("debug")

					userInputs.CLI.Highlight = command.Bool("highlight")

					userInputs.CLI.FileUploadPath = command.String("fileupload")

					userInputs.CLI.DownloadPath = command.String("download")

					// Handle custom body

					userInputs.Req.SetBodyString(command.String("body"))

					// Handle custom headers
					headerSlice := command.StringSlice("header")
					//*headers = make(map[string]string)

					for _, h := range headerSlice {
						key, value, _ := strings.Cut(h, ":")

						if len(key) >= 1 {
							//(headers)[key] = value
							//headers[key] = []byte(value)

							userInputs.Req.Header.Set(key, value)
						}
					}
					if userInputs.Req.Header.Peek("Content-Type") == nil {
						userInputs.Req.Header.Add("Content-Type", "text/html")
					}

					// Handle replacement variables
					replacementSlice := command.StringSlice("var")
					//*replacementPair = make(map[string]string)

					for _, h := range replacementSlice {
						targetWord, value, _ := strings.Cut(h, "=")

						if len(targetWord) >= 1 {

							(replacementPair)[targetWord] = value

						}
					}

					cmd.SendHTTP(&replacementPair, userInputs)

					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
