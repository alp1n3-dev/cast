package cmd

import (
	"fmt"
	"net/http"
	"os"
	"slices"


	"github.com/spf13/cobra"

	//"github.com/spf13/viper"
	"strings"

	"github.com/alp1n3-eth/cast/extractors/http_extractors"
	"github.com/alp1n3-eth/cast/models"
	"github.com/alp1n3-eth/cast/output"
	"github.com/alp1n3-eth/cast/pkg/logging"

	//"github.com/alp1n3-eth/cast/models"
)

func init() {
	rootCmd.Args = cobra.MatchAll(cobra.OnlyValidArgs, cobra.MinimumNArgs(1))
	rootCmd.Run = func(cmd *cobra.Command, args []string) {

		logging.Init(true) // Debug mode is TRUE

		methodList := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTION", "TRACE"}
		argZero := strings.ToLower(args[0])

		if slices.Contains(methodList, strings.ToUpper(argZero)) {
			//fmt.Println("[!] Method and URL Provided")
			logging.Logger.Info("Method and URL Provided")
			method := strings.ToUpper(args[0])
			url := args[1]
			// TODO: Parse and send custom headers
			// TODO: Double-check that multiple cookies being set doesn't run into any issues
			// TODO: Add the ability to add a custom body to POSTS
			request := http_extractors.BuildHTTPRequest(method, url)
     		SendHTTPRequest(request)
		} else if len(argZero) > 5 && argZero[len(argZero)-5:] == ".http"{
			//fmt.Println("[!] HTTP File Provided")
			logging.Logger.Info("HTTP File Provided")

			// Parse the provided file.

			// Using the ParsedHTTPFile struct, create an array of HTTP Requests
		} else {
			logging.Logger.Fatal("Error reading user args. Initial arg user provided: " + args[0])
		}
	}
}

func SendHTTPRequest(r models.HTTPRequest) (models.HTTPRequest, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Request.Method, r.Request.URL, nil)
	if err != nil {
			fmt.Println("Error creating request:", err)
			return r, nil
	}

	resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			os.Exit(1)
		}
	defer resp.Body.Close()


	rPair, err := http_extractors.BuildHTTPResponse(r, resp)
	// Print the Response
	// Returns *http.Response & error
	err = output.PrintResponse(rPair)

	return rPair, nil
}
