package tests

import (
	"fmt"

	"testing"
	"time"
)

func BenchmarkGETCast(b *testing.B) {
	commands := []CommandBenchmark{
		{"cast_get", []string{"../cast", "get", "https://www.google.com/"}},
		//{"cast_post", []string{"../cast", "post", "https://echo.free.beeceptor.com", "test=test1"}},
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
