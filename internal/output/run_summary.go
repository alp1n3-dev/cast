package output

import (
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/models"
	"github.com/fatih/color"
)

// Outputs status of a single or multi-file run.
func FileRun(results *models.ResultOut) error {
	filename := "test.http"
	success := true

	d := color.New(color.Bold)
	d.Printf("%s: ", filename)
	if success {
		d.Add(color.FgGreen)
		d.Printf("Success\n")
	} else {
		d.Add(color.FgRed)
		d.Printf("Failed\n")
	}

	fmt.Println("────────────────────")
	//fmt.Printf("# of files: %s\n")
	fmt.Printf("# of requests: %d\n", results.RequestTotal)
	//fmt.Printf("# of successes: %s\n")
	fmt.Printf("# of failures: %d\n", results.FailureTotal)
	fmt.Printf("Duration: %d ms\n", results.Duration)

	if !success {
		return fmt.Errorf("file run assertion failure: %s", filename)
	}

	return nil
}
