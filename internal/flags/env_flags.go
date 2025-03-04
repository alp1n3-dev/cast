package flags

import (
	"github.com/urfave/cli/v3"
)

var (
	EnvAddFlag = &cli.StringFlag{
		Name:  "env",
		Value: "",
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	EnvPurgeFlag = &cli.BoolFlag{
		Name:  "purge",
		Value: false,
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	EnvPersistentAdd = &cli.BoolFlag{
		Name:  "persist",
		Value: false,
		Usage: "Runs the file without checking existing assertions on each response.",
	}

	EnvFlags = []cli.Flag{
		NoAssertsFlag,
	}
)
