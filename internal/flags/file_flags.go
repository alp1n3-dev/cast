package flags

import (
	"github.com/urfave/cli/v3"
)

var (
	FileDebugFlag = &cli.BoolFlag{
		Name:     "fdebug",
		Usage:    "Enable file debug output for the application",
		Value:    false,
		OnlyOnce: true,
		Category: "utility",
	}
	NoAssertsFlag = &cli.StringFlag{
		Name:  "no-asserts",
		Value: "",
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	FileFlags = []cli.Flag{
		NoAssertsFlag,
		FileDebugFlag,
	}
)
