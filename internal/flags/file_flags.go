package flags

import (
	"github.com/urfave/cli/v3"
)

var (
	IgnoreAssertsFlag = &cli.StringFlag{
		Name:  "ignore-assertions",
		Value: "",
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	FileFlags = []cli.Flag{
		IgnoreAssertsFlag,
	}
)
