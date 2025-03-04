package flags

import (
	"maps"
	"strings"

	"github.com/alp1n3-eth/cast/internal/env"
	"github.com/urfave/cli/v3"
)

var (
	BodyFlag = &cli.StringFlag{
		Name:     "body",
		Value:    "",
		Usage:    "HTTP request body",
		Aliases:  []string{"B"},
		Category: "http",
	}
	HeaderFlag = &cli.StringSliceFlag{
		Name:     "header",
		Usage:    "HTTP headers to include in the request",
		Aliases:  []string{"H"},
		Category: "http",
	}
	DebugFlag = &cli.BoolFlag{
		Name:     "debug",
		Usage:    "Enable debug output for the application",
		Value:    false,
		OnlyOnce: true,
		Category: "utility",
	}
	ColorFlag = &cli.BoolFlag{
		Name:     "color",
		Usage:    "Prettify the response body output with syntax highlighting",
		Value:    false,
		OnlyOnce: true,
		Category: "utility",
	}
	MoreFlag = &cli.BoolFlag{
		Name:     "more",
		Usage:    "Print the entire response body instead of only the first 200 characters.",
		Value:    false,
		OnlyOnce: true,
		Category: "utility",
	}
	VarFlag = &cli.StringSliceFlag{
		Name:     "var",
		Usage:    "A text file to be iteratively used in any portion of the request to insert values.",
		Aliases:  []string{"V"},
		Category: "utility",
	}
	PrintFlag = &cli.StringSliceFlag{
		Name:     "print",
		Usage:    "A text file to be iteratively used in any portion of the request to insert values.",
		Aliases:  []string{"P"},
		Category: "utility",
	}
	FileFlag = &cli.StringFlag{
		Name:      "file",
		Usage:     "A way to include a file in the request's body.",
		Aliases:   []string{"F"},
		TakesFile: true,
		OnlyOnce:  true,
		Category:  "utility",
	}
	RedirectFlag = &cli.IntFlag{
		Name:     "redirect",
		Usage:    "A way to follow redirects up to < INT >.",
		Aliases:  []string{"R"},
		OnlyOnce: true,
		Category: "http",
	}
	DownloadFlag = &cli.StringFlag{
		Name:      "download",
		Usage:     "Path to save the response body to a file.",
		Aliases:   []string{"D"},
		TakesFile: true,
		OnlyOnce:  true,
	}
	ReadEncryptedFlag = &cli.BoolFlag{
		Name:     "read-encrypted",
		Usage:    "Read an encrypted KV store using a password.",
		OnlyOnce: true,
		Category: "env",
	}
	CurlFlag = &cli.BoolFlag{
		Name:     "curl",
		Usage:    "Output the request as a curl command instead of sending it.",
		OnlyOnce: true,
		Category: "export",
	}

	GetFlags = []cli.Flag{
		BodyFlag,
		HeaderFlag,
		DebugFlag,
		ColorFlag,
		MoreFlag,
		VarFlag,
		PrintFlag,
		FileFlag,
		RedirectFlag,
		DownloadFlag,
		ReadEncryptedFlag,
		CurlFlag,
	}
)

func ParseReplacementValues(replacementSlice []string) *map[string]string {
	replacementPair := make(map[string]string)

	for _, h := range replacementSlice {
		if strings.Contains(h, ".env") {
			kvFileMap, _ := env.ReadKVFile(h)

			maps.Copy(replacementPair, *kvFileMap)
		} else {
			targetWord, value, _ := strings.Cut(h, "=")

			if len(targetWord) >= 1 {
				replacementPair[targetWord] = value
			}
		}
	}
	return &replacementPair
}
