package tests

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

// CommandBenchmark defines a command to be benchmarked
type CommandBenchmark struct {
	Name string
	Args []string
}

func BenchmarkHTTPClients(b *testing.B) {
	commands := []CommandBenchmark{
		//{"httpie", []string{"/opt/homebrew/bin/http", "POST", "https://echo.free.beeceptor.com", "test=test1"}},
		{"http_get", []string{"/opt/homebrew/bin/http", "GET", "https://www.google.com/"} },
		{"xh_get", []string{"/opt/homebrew/bin/xh", "GET", "https://www.google.com/"}},
		{"cast_get", []string{"../cast", "get", "https://www.google.com/"}},
		//{"cast_post", []string{"../cast", "post", "https://echo.free.beeceptor.com", "test=test1"}},
		//{"hurl", []string{"hurl", "hurl_test.hurl"}},
		{"hurl_get", []string{"hurl", "hurl_test_google.hurl"}},
	}

	// Store results for comparison
	results := make(map[string]time.Duration)

	for _, cmdBenchmark := range commands {
		b.Run(cmdBenchmark.Name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer() // Reset timer to exclude setup time

			start := time.Now()
			for i := 0; i < b.N; i++ {
				runCommand(b, cmdBenchmark.Name, cmdBenchmark.Args)
			}
			duration := time.Since(start)

			results[cmdBenchmark.Name] = duration
		})
	}

	// Print summary of benchmark results
	fmt.Println("\n=== Benchmark Results ===")
	for name, duration := range results {
		fmt.Printf("%-10s: %v per %d runs\n", name, duration, b.N)
	}
}

// runCommand executes a command and reports errors
func runCommand(b *testing.B, name string, args []string) {
	if len(args) == 0 {
		b.Fatalf("[%s] No command provided", name)
	}

	cmd := exec.Command(args[0], args[1:]...) // Create a new command instance
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		b.Logf("[%s] Error: %v", name, err)
		b.Logf("[%s] Stderr: %s", name, stderr.String())
		b.Logf("[%s] Output: %s", name, out.String())
		b.Fail()
	}
}
