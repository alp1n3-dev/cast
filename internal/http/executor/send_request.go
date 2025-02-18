package executor

import (
	"net/http"
	"fmt"
	"os"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/alp1n3-eth/cast/internal/http/parse"
)

func SendHTTPRequest(r models.ExecutionResult) (models.ExecutionResult, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Request.Method, r.Request.URL.String(), nil)
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




	//rPair, err := http_extractors.BuildHTTPResponse(r, resp)
	r.Response = parse.BuildResponse(resp)

	// Print the Response
	// Returns *http.Response & error
	//err = output.PrintResponse(r)

	return r, nil
}
