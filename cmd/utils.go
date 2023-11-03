package cmd

import (
	"fmt"
	"os"
)

// Verbose prints verbose messages to stdout
func Verbose(format string, args ...interface{}) {
	if Log != nil {
		Log.Infof(format, args...)
	}
	if CmdOptions.Verbose {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
	}
}

// Error prints error messages to stderr
func Error(format string, args ...interface{}) {
	if Log != nil {
		Log.Errorf(format, args...)
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}

// Die prints error messages to stderr and exit the app with an exitCode
func Die(code int, format string, args ...interface{}) {
	if Log != nil {
		Log.Fatalf(format, args...)
	}
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
	os.Exit(code)
}
