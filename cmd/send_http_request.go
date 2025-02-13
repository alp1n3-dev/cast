package cmd

import (
	"fmt"
	"net/http"
	"os"


	"github.com/spf13/cobra"

	//"github.com/spf13/viper"
	"strings"

	"github.com/alp1n3-eth/cast/extractors/http_extractors"
	"github.com/alp1n3-eth/cast/models"
	"github.com/alp1n3-eth/cast/output"
	//"github.com/alp1n3-eth/cast/models"
)

func init() {
	rootCmd.Args = cobra.MinimumNArgs(2)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {

		method := strings.ToUpper(args[0])
		url := args[1]
		//fmt.Println(method)
		//fmt.Println(url)

		//http_extractors.ValidateHTTP(method, url)
		request := http_extractors.BuildHTTPRequest(method, url)
		SendHTTPRequest(request)

		//makeRequest(method, url)

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
