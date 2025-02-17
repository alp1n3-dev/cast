package cmd

// Test is deprecated currently due to ongoing program architect issues.

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

// TestSendRequest using a mock HTTP server
func TestSendRequest(t *testing.T) {
	// Mock server response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "ok"}`)
	}))
	defer mockServer.Close()

	// Create request
	//req, err := http.NewRequest("GET", mockServer.URL, nil)
	//if err != nil {
		//t.Fatalf("Error creating request: %v", err)
		//}

	// Call function (note: it prints output instead of returning values)
	var request models.ExecutionResult
	var err error
	request.Request.URL, err = url.Parse(mockServer.URL)
	if err != nil {
		logging.Logger.Fatal("Unable to assign mock server URL to request.")
	}
	request.Request.Method = "GET"
	request, err = SendHTTPRequest(request)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}


	//convResp, _ := io.ReadAll(resp.Body)
	//fmt.Println("Returned response body from SendHTTPRequest: "+ string(convResp))
}
