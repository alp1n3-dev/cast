package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"github.com/alp1n3-eth/cast/models"
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
	var request models.HTTPRequest
	request.Request.URL = mockServer.URL
	request.Request.Method = "GET"
	resp, err := SendHTTPRequest(request)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}


	convResp, _ := io.ReadAll(resp.Body)
	fmt.Println("Returned response body from SendHTTPRequest: "+ string(convResp))
}
