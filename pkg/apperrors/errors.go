package apperrors

import (
	"errors"
	"os"
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/logging"
	"github.com/alp1n3-eth/cast/pkg/models"
)

var (
    ErrInvalidHeaderFormat = errors.New("Invalid Header Format")
    ErrBodyParseFailed     = errors.New("Body Parsing Failed")
    ErrRequestCreation     = errors.New("Request Creation Failed")
)

func Wrap(err error, message string) error {
    return fmt.Errorf("%s: %w", message, err)
}

func HandleExecutionError(e error) {
	var execErr *models.ExecutionError

	logging.Logger.Error("reached 2")

	fmt.Println(execErr.Stage)
	fmt.Println(e)
	fmt.Println(execErr.Message)

    if errors.As(e, &execErr) {
        logging.Logger.Error("Execution failed",
            "stage", execErr.Stage,
            "error", e,
            "message", execErr.Message)
        os.Exit(1)
    }
}
