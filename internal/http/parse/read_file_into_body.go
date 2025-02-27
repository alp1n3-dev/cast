package parse

import (
	"io"
	"os"

	"github.com/alp1n3-eth/cast/pkg/logging"
)

func ReadFileIntoBody(uploadFilePath *string) []byte {
	var fileContents []byte

	// Check if file exists
	if _, err := os.Stat(*uploadFilePath); os.IsNotExist(err) {
		logging.Logger.Fatalf("File does not exist: %s", *uploadFilePath)
	}

	file, err := os.Open(*uploadFilePath)
	if err != nil {
		logging.Logger.Fatal("Error opening file")
	}
	defer file.Close()

	fileContents, err = io.ReadAll(file)
	if err != nil {
		logging.Logger.Fatal("Error reading file")
	}

	logging.Logger.Debugf("fileContents: %x", fileContents)
	logging.Logger.Debug("Successfully read file contents for file upload.")

	return fileContents
}
