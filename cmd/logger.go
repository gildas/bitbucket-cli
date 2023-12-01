package cmd

import (
	"os"

	"github.com/gildas/go-logger"
)

// Log is the logger for this application
var Log *logger.Logger

// createLogger creates the logger for this application
func createLogger() {
	if CmdOptions.Debug {
		os.Setenv("DEBUG", "true")
	}
	if len(CmdOptions.LogDestination) > 0 {
		os.Setenv("LOG_DESTINATION", CmdOptions.LogDestination)
	}
	if len(os.Getenv("LOG_DESTINATION")) == 0 {
		Log = logger.Create(APP, &logger.NilStream{})
	} else {
		Log = logger.Create(APP)
	}
}
