package cmd

import (
	"context"
	"fmt"

	//"log"
	"os"
	"strings"

	//httpcmd "github.com/alp1n3-eth/cast/cmd/actions"
	"github.com/alp1n3-eth/cast/internal/env"
	"github.com/alp1n3-eth/cast/internal/flags"
	"github.com/alp1n3-eth/cast/internal/http/executor"
	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/urfave/cli/v3"
	"github.com/valyala/fasthttp"
)

var (
	GetCommand = &cli.Command{ // HTTP Methods
		Name:    "get",
		Aliases: []string{"post", "put", "delete", "patch", "options", "trace", "head", "connect"},
		Usage:   "Send an HTTP request to a url.",
		Flags:   flags.GetFlags,
		Action:  GetAction,
		//Suggest: true,
	}
	FileCommand = &cli.Command{
		Name:   "file",
		Usage:  "Run HTTP requests from a provided file.",
		Flags:  flags.FileFlags,
		Action: FileAction,
		//Suggest: true,
	}
)

func Execute(ctx context.Context, args []string) error {
	app := &cli.Command{
		Name:    "cast",
		Version: "v0.4-alpha",
		//Authors: any[],
		Usage:     "make sending HTTP requests ezpz",
		UsageText: "placeholder text",
		ArgsUsage: "[Method] [Protocol + Host] <Flags>",
		Commands: []*cli.Command{
			GetCommand,
			FileCommand,
		},
	}

	return app.Run(ctx, args)
}

func GetAction(ctx context.Context, command *cli.Command) error {
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
			Color:             command.Bool("color"),
			More:              command.Bool("more"),
			FileUploadPath:    command.String("fileupload"),
			DownloadPath:      command.String("download"),
			CurlOutput:        command.Bool("curl"),
		},
	}

	userInputs.Request.Req.Header.SetMethod(strings.ToUpper(os.Args[1]))
	userInputs.Request.Req.SetRequestURI(command.Args().First())

	// TODO: Create the ability to encrypt the persistent ENVs
	if command.Bool("read-encrypted") {
		fmt.Print("Enter password: ")
		password, err := env.RetrievePasswordFromUser()
		if err != nil {
			fmt.Println("error retrieving password")
		}
		fmt.Println(password)
	}

	userInputs.Request.Req.SetBodyRaw([]byte(command.String("body")))

	for _, h := range command.StringSlice("header") {
		key, value, _ := strings.Cut(h, ":")
		if len(key) >= 1 {
			userInputs.Request.Req.Header.Set(key, value)
		}
	}
	if userInputs.Request.Req.Header.Peek("Content-Type") == nil {
		userInputs.Request.Req.Header.Add("Content-Type", "text/html")
	}

	replacementPair := flags.ParseReplacementValues(command.StringSlice("var"))

	executor.SendHTTP(replacementPair, userInputs)

	return nil
}

func FileAction(ctx context.Context, command *cli.Command) error {
	return nil
}

func EnvAction(ctx context.Context, command *cli.Command) error {
	return nil
}

func init() {
	//cli.CommandHelpTemplate = "get"
	//cli.SuggestCommand

}
