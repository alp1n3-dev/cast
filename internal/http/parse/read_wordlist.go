package parse

import (
	"os"
	"bufio"

	"github.com/alp1n3-eth/cast/pkg/models"
)

// Step 1: Read the wordlist / make sure the file exists.
func ReadWordlistFile (wordlistStruct *models.Wordlist) error {
	file, err := os.Open(wordlistStruct.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordlistStruct.Words = append(wordlistStruct.Words, scanner.Text())
	}

	return nil
}

// Step 2: Replace the vals in the HTTP request that was built right before sending it.
// Replaces the designated target value(s) of a single request.
// Looks for multiple occurrences of the value, but won't do every word in the list.
// AKA: Meant to be called in a "for" loop
func SwapWordlistVals (wordlistStruct *models.Wordlist) error {

	return nil
}
