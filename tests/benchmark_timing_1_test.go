package tests

import (
	"os/exec"
	"testing"
	"fmt"
	"bytes"
)

func BenchmarkCastBinary(b *testing.B) {
	//cmdPath := "../cast" // Path to your built binary
	//args := []string{"post", "https://echo.free.beeceptor.com", "--body", "test=test1"}

	//cmdPath := "http"
	//args := []string{"POST", "https://echo.free.beeceptor.com", "test=test1"}




	b.ResetTimer() // Reset the timer to ignore setup time

	for i := 0; i < b.N; i++ {
	//for i := 0; i < 5; i++ {
		//cmd := exec.Command("/opt/homebrew/bin/http", "POST", "echo.free.beeceptor.com",  "test==test1")

		cmd := exec.Command("../cast", "get", "https://www.google.com/")
		//cmd := exec.Command("../cast", "post", "https://echo.free.beeceptor.com",  "test=test1")

		//cmd := exec.Command("xh", "get", "https://www.google.com/")

		//cmd := exec.Command("hurl", "hurl_test.hurl")
		//cmd := exec.Command("hurl", "hurl_test_google.hurl")


		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Run the command and wait for completion
		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Stderr:", stderr.String())
			fmt.Println("Output:", out.String())
			b.Fatalf("Command: %v execution failed: %v", cmd, err)
		}
	}
}
