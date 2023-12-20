package common

import (
	"context"

	"github.com/spf13/cobra"
)

type RemoteValueFlag struct {
	AllowedFunc func(context.Context, *cobra.Command) []string
	Value       string
}

// Type returns the type of the flag
func (flag RemoteValueFlag) Type() string {
	return "string"
}

// String returns the string representation of the flag
func (flag RemoteValueFlag) String() string {
	return flag.Value
}

// Set sets the flag value
func (flag *RemoteValueFlag) Set(value string) error {
	flag.Value = value
	return nil
}

// CompletionFunc returns the completion function of the flag
func (flag *RemoteValueFlag) CompletionFunc() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if flag.AllowedFunc != nil {
			return flag.AllowedFunc(cmd.Context(), cmd), cobra.ShellCompDirectiveNoFileComp
		}
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
}
