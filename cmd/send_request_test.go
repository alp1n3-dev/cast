package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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
	req, err := http.NewRequest("GET", mockServer.URL, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Call function (note: it prints output instead of returning values)
	SendRequest(req)
}
