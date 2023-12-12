package main

import "strings"

// commit contains the current git commit this code was built on and should be set via -ldflags
var commit string

// branch contains the git branch this code was built on and should be set via -ldflags
var branch string

// stamp contains the build date and should be set via -ldflags
var stamp string

// VERSION is the version of this application
var VERSION = "0.3.0"

// APP is the name of the application
const APP = "bb"

// Version gets the current version of the application
func Version() string {
	if strings.HasPrefix(strings.ToLower(branch), "dev") || strings.HasPrefix(strings.ToLower(branch), "feature") {
		return VERSION + "+" + stamp + "." + commit
	}
	return VERSION
}
