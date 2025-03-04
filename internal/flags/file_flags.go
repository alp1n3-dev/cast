package flags

import (
	"github.com/urfave/cli/v3"
)

var (
	NoAssertsFlag = &cli.StringFlag{
		Name:  "no-asserts",
		Value: "",
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	FileFlags = []cli.Flag{
		NoAssertsFlag,
	}
)
