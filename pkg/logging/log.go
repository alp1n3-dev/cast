package logging

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func Init(debug bool) {

	/*
		Logger = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			Formatter:       log.TextFormatter,
		})*/
	if debug {
		Logger = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			Formatter:       log.TextFormatter,
		})
		Logger.SetLevel(log.DebugLevel)
	}
	if !debug {
		Logger = log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: false,
			Formatter:       log.TextFormatter,
		})
		//Logger.SetLevel(log.ErrorLevel)
	}
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().SetString("ERROR").Foreground(lipgloss.Color("204")).Bold(true)

	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Bold(true)

	Logger.SetStyles(styles)
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
