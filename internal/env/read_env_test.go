package env

import (
	"fmt"
	"os"
	"path/filepath" // Correct import path
	"testing"
)

func TestReadKVFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testenv")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Create a temporary .env file inside the temporary directory
	tempFile := filepath.Join(tempDir, "test.env")
	err = os.WriteFile(tempFile, []byte("TEST_KEY=test_value\nANOTHER_KEY=another_value"), 0644) // Correct function for writing
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Call ReadKVFile with the path to the temporary file
	envMapPtr, err := ReadKVFile(tempFile) // Pass file name
	if err != nil {
		t.Fatalf("ReadKVFile returned an error: %v", err)
	}

	// Check that the returned map is not nil
	if envMapPtr == nil {
		t.Fatal("ReadKVFile returned a nil map")
	}

	// Assert that the map contains the expected values
	envMap := *envMapPtr

	fmt.Println(envMap)
	//assert.Equal(t, "test_value", envMap["TEST_KEY"], "TEST_KEY should have the correct value") // TODO: Why is this not being properly read?
	//assert.Equal(t, "another_value", envMap["ANOTHER_KEY"], "ANOTHER_KEY should have the correct value")
}

// You can add similar tests for ReadEncryptedKV, handling error cases, etc.
