package cmd

import (
	"net/http"
	"fmt"
	"os"
)

func SendRequest(req *http.Request) {
	client := &http.Client{}

	resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			os.Exit(1)
		}
	defer resp.Body.Close()

	PrintResponse(resp)

}
