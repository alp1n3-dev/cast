package output

import (
	"fmt"

	"github.com/fatih/color"
)

// Outputs status of a single or multi-file run.
func FileRun() error {
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
	fmt.Printf("# of files: %s\n")
	fmt.Printf("# of requests: %s\n")
	fmt.Printf("# of successes: %s\n")
	fmt.Printf("# of failures: %s\n")
	fmt.Printf("Duration: %s ms\n")

	if !success {
		return fmt.Errorf("file run assertion failure: %s", filename)
	}

	return nil
}
