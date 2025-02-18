/*
Copyright © 2025 alp1n3 1@alp1n3.dev

*/
package main

import (
    "os"
    "fmt"
    "log"

    "github.com/urfave/cli/v2"

    "github.com/alp1n3-eth/cast/cmd"

)

func main() {
    app := &cli.App{
        Commands: []*cli.Command{
            {
                Name:    "get",
                Aliases: []string{"GET"},
                Usage:   "send a get request to a url.",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("added task: ", cCtx.Args().First())
                    cmd.SendHTTP("get", cCtx.Args().First(), "test2")
                    return nil
                },
            },
            {
                Name:    "complete",
                Aliases: []string{"c"},
                Usage:   "complete a task on the list",
                Action: func(cCtx *cli.Context) error {
                    fmt.Println("completed task: ", cCtx.Args().First())
                    return nil
                },
            },
            {
                Name:    "template",
                Aliases: []string{"t"},
                Usage:   "options for task templates",
                Subcommands: []*cli.Command{
                    {
                        Name:  "add",
                        Usage: "add a new template",
                        Action: func(cCtx *cli.Context) error {
                            fmt.Println("new task template: ", cCtx.Args().First())
                            return nil
                        },
                    },
                    {
                        Name:  "remove",
                        Usage: "remove an existing template",
                        Action: func(cCtx *cli.Context) error {
                            fmt.Println("removed task template: ", cCtx.Args().First())
                            return nil
                        },
                    },
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
