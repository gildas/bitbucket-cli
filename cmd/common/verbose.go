package common

import (
	"context"
	"fmt"

	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// Verbose prints a message if the verbose flag is set
func Verbose(context context.Context, cmd *cobra.Command, format string, args ...any) {
	logger.Must(logger.FromContext(context)).Infof(format, args...)
	if cmd.Flag("verbose").Changed {
		fmt.Printf(format, args...)
		fmt.Println()
	}
}
