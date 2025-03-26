package cmd

import (
	"context"
	"fmt"
	"time"

	//"log"
	"os"
	"strings"

	//httpcmd "github.com/alp1n3-eth/cast/cmd/actions"
	"github.com/alp1n3-eth/cast/internal/assert"
	"github.com/alp1n3-eth/cast/internal/capture"
	"github.com/alp1n3-eth/cast/internal/env"
	"github.com/alp1n3-eth/cast/internal/executor"
	"github.com/alp1n3-eth/cast/internal/flags"
	output "github.com/alp1n3-eth/cast/internal/output"
	"github.com/alp1n3-eth/cast/internal/parse"
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
		Name:    "Cast",
		Version: "v0.4-alpha",
		//Authors: any[],
		Usage:     "Make sending HTTP requests & testing APIs EZPZ.\nSubmit any bugs as an issue in the repo: https://github.com/alp1n3-eth/cast/issues",
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
	filePath := command.Args().First()
	if filePath == "" {
		return fmt.Errorf("file path not provided")
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	//type customParser struct{}

	parser := parse.CustomParser{}
	castFile, err := parser.ParseToCastFile(fileContent)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	var results models.ResultOut

	startTime := time.Now()

	for i := 0; i < len(castFile.CtxMap); i++ {
		var reqCtx models.HTTPRequestContext
		reqCtx = castFile.CtxMap[i]

		var replacementPlaceholder map[string]string

		executor.SendHTTP(&replacementPlaceholder, &reqCtx)
		results.RequestTotal += 1

		assert.ValidateAssertions(&reqCtx.Response, reqCtx.Assertions, &results)

		if i < (len(castFile.CtxMap) - 1) {
			fmt.Println()
		}

		capture.Capture(&reqCtx)
		//fmt.Println("global vars")
		//fmt.Println(capture.GlobalVars)
	}
	results.Duration = int(time.Since(startTime).Milliseconds())

	output.FileRun(&results)

	return nil
}

func EnvAction(ctx context.Context, command *cli.Command) error {
	return nil
}

func init() {
	//cli.CommandHelpTemplate = "get"
	//cli.SuggestCommand

}
