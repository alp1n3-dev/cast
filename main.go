/*
Copyright © 2025 alp1n3 1@alp1n3.dev

*/
package main

import (
    "os"
    //"fmt"
    "log"
    "context"
    //"net/http"

    "github.com/urfave/cli/v3"

    "github.com/alp1n3-eth/cast/cmd/http"
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
                 		Name: "body",
                   		Value: "",
                     	Usage: "HTTP request body",
                      	Aliases: []string{"b"},
                 },
                 	&cli.StringFlag{
                  		Name: "header",
                    	Value: "",
                     	Usage: "HTTP headers to include in the request",
                      	//Aliases: []string{"h"},
                  },
                },
                Action: func(ctx context.Context, command *cli.Command) error {
                    //fmt.Println("added task: ", command.Args().First())
                    //fmt.Println("Debug - All args:", os.Args)
                    //fmt.Println("Debug - First arg:", os.Args[1])
                    //fmt.Println("Debug - Context args:", command.Args().Slice())
                    //fmt.Println("Debug - Body flag:", command.String("body"))
                    body := command.String("body")
                    headers := command.StringMap("header")
                    //fmt.Println(body)
                    cmd.SendHTTP(os.Args[1], command.Args().First(), body, headers)
                    return nil
                },
            },
        },
    }

    if err := app.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
