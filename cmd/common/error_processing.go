package common

import (
	"strings"

	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type ErrorProcessing int

const (
	StopOnError ErrorProcessing = iota
	WarnOnError
	IgnoreErrors
)

// Type returns the type of the ErrorProcessing
//
// implements the pflag.Value interface
func (ep ErrorProcessing) Type() string {
	return "string"
}

// Values returns the allowed values of the ErrorProcessing
func (ep ErrorProcessing) Values() []string {
	return []string{StopOnError.String(), WarnOnError.String(), IgnoreErrors.String()}
}

// Set sets the ErrorProcessing value
//
// implements the pflag.Value interface
func (ep *ErrorProcessing) Set(value string) error {
	switch value {
	case "StopOnError":
		*ep = StopOnError
	case "WarnOnError":
		*ep = WarnOnError
	case "IgnoreErrors":
		*ep = IgnoreErrors
	default:
		return errors.ArgumentInvalid.With("value", value, strings.Join(ep.Values(), ", "))
	}
	return nil
}

// CompletionFunc returns the completion function of the ErrorProcessing
//
// implements the pflag.Value interface
func (ep ErrorProcessing) CompletionFunc() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return ep.Values(), cobra.ShellCompDirectiveNoFileComp
	}
}

// String returns the string representation of the ErrorProcessing
//
// implements the fmt.Stringer interface
func (ep ErrorProcessing) String() string {
	return [...]string{"StopOnError", "WarnOnError", "IgnoreErrors"}[ep]
}
