package common

import (
	"context"
	"fmt"
	"os"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// WhatIf prints what would be done by the command
//
// # If the DryRun flag is set, it prints what would be done by the command and tells the caller to not proceed
//
// otherwise it does nothing
func WhatIf(context context.Context, cmd *cobra.Command, format string, args ...any) (proceed bool) {
	dryRun := false
	for _, name := range []string{"dry-run", "noop", "whatif"} {
		flag := cmd.Flag(name)
		if flag != nil && flag.Value != nil && flag.Value.String() == "true" {
			dryRun = true
			break
		}
	}

	if dryRun {
		logger.Must(logger.FromContext(context)).Infof("Dry run: "+format, args...)
		fmt.Fprintf(os.Stderr, "Dry run: "+format+"\n", args...)
		return false
	}
	return true
}
