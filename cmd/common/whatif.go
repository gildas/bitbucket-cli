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
	if cmd.Flag("dry-run").Changed {
		logger.Must(logger.FromContext(context)).Infof("Dry run: "+format, args...)
		fmt.Fprintf(os.Stderr, "Dry run: "+format+"\n", args...)
		return false
	}
	return true
}
