package logging

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func Init(debug bool) {
    Logger = log.NewWithOptions(os.Stderr, log.Options{
        ReportCaller:   true,
        ReportTimestamp: true,
        Formatter: log.TextFormatter,
    })
    if debug {
        Logger.SetLevel(log.DebugLevel)
    }
}

//var Logger = log.NewWithOptions(os.Stderr, log.Options{
    //ReportCaller:   true,
    //ReportTimestamp: true,
//})

/*
Usage in commands:

func Execute() {
    logging.Init(debugMode) // debugMode needs to equal true or false.
    logging.Logger.Info("Starting request", "url", targetURL) // Will show every time.
    logging.Logger.Debug("Raw headers", headers) // Will only show with debug mode turned on.
}
*/
