/*
Copyright © 2025 alp1n3 1@alp1n3.dev
*/
package main

import (
	"net/http"
	"os"
	"fmt"
	"context"
	"log"
	"strings"
	"bytes"
	"io"

	//"net/http"

	"github.com/urfave/cli/v3" // docs: https://cli.urfave.org/v3/examples/subcommands/

	"github.com/alp1n3-eth/cast/cmd/http"
)

func main() {
	headers := &http.Header{}
	//bodyReader := &io.Reader
	var body io.Reader
	bodyReader := &body

    app := &cli.Command{
        Commands: []*cli.Command{
            {
                Name:    "get",
                Aliases: []string{"GET", "post", "put", "delete", "patch", "options", "trace", "head", "connect"},
                Usage:   "send an HTTP request to a url.",
                Flags: []cli.Flag{
                	&cli.StringFlag{
                 		Name: "body",
                   		Value: "",
                     	Usage: "HTTP request body",
                      	Aliases: []string{"B"},
                 },
                 	&cli.StringSliceFlag{
                  		Name: "header",
                     	Usage: "HTTP headers to include in the request",
                      	Aliases: []string{"H"},
                  },
                  &cli.BoolFlag{
                  		Name: "debug",
                     	Usage: "Enable debug output for the application",
                      	Aliases: []string{"D"},
                  },
                  &cli.BoolFlag{
                  		Name: "highlight",
                     	Usage: "Prettify the response body output with syntax highlighting",
                      	Aliases: []string{"HL"},
                  },
                },
                Action: func(ctx context.Context, command *cli.Command) error {
                    //fmt.Println("added task: ", command.Args().First())
                    //fmt.Println("Debug - All args:", os.Args)
                    //fmt.Println("Debug - First arg:", os.Args[1])
                    //fmt.Println("Debug - Context args:", command.Args().Slice())
                    //fmt.Println("Debug - Body flag:", command.String("body"))

                    //var bodyReader io.Reader

                    debug := command.Bool("debug")
                    highlight := command.Bool("highlight")
                    //fmt.Println(highlight)


                    bodyString := command.String("body")
                    if bodyString != "" {
                    	*bodyReader = bytes.NewBufferString(bodyString)
                    }


                    headerSlice := command.StringSlice("header")
                    if headerSlice != nil {
                    	//fmt.Println("headers not nil")


                      	for _, h := range headerSlice {
                            parts := strings.SplitN(h, ":", 2)
                            if len(parts) == 2 {
                                key := strings.TrimSpace(parts[0])
                                value := strings.TrimSpace(parts[1])

                                headers.Add(key, value)
                            }
                        }
                    }

                    //headers := make(http.Header)

                    //fmt.Println(body)
                    cmd.SendHTTP(os.Args[1], command.Args().First(), bodyReader, headers, debug, highlight)
                    return nil
                },
            },
        },
    }

    if err := app.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
