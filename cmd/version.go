package cmd

import "strings"

// commit contains the current git commit and is set in the build.sh script
var app_commit string

// branch contains the git branch this code was built on and should be set via -ldflags
var app_branch string

// stamp contains the build date and should be set via -ldflags
var app_stamp string

// VERSION is the version of this application
var VERSION = "0.0.0"

const APP = "bb"

// Version gets the current version of the application
func Version() string {
	if strings.HasPrefix(strings.ToLower(app_branch), "dev") || strings.HasPrefix(strings.ToLower(app_branch), "feature") {
		return VERSION + "+" + app_stamp + "." + app_commit
	}
	return VERSION
}
