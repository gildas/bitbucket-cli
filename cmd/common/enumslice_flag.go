package common

import (
	"strings"

	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type EnumSliceFlag struct {
	Allowed    []string
	Values     []string
	Default    []string
	AllAllowed bool
	all        bool
}

// Type returns the type of the flag
func (flag EnumSliceFlag) Type() string {
	return "stringSlice"
}

// String returns the string representation of the flag
func (flag EnumSliceFlag) String() string {
	return strings.Join(flag.Values, ",")
}

// Set sets the flag value
func (flag *EnumSliceFlag) Set(value string) error {
	if value == "all" && flag.AllAllowed {
		flag.Values = flag.Allowed
		flag.all = true
		return nil
	}
	for _, allowed := range flag.Allowed {
		if value == allowed {
			for _, existing := range flag.Values {
				if existing == value {
					return nil
				}
			}
			flag.Values = append(flag.Values, value)
			return nil
		}
	}
	return errors.ArgumentInvalid.With("value", value, strings.Join(flag.Allowed, ", "))
}

// Get returns the flag value
func (flag EnumSliceFlag) Get() []string {
	if len(flag.Values) == 0 {
		return flag.Default
	}
	return flag.Values
}

// Contains returns true if the flag contains the given value
func (flag EnumSliceFlag) Contains(value string) bool {
	for _, allowed := range flag.Allowed {
		if value != allowed {
			return false
		}
	}
	if flag.all {
		return true
	}
	for _, existing := range flag.Values {
		if existing == value {
			return true
		}
	}
	return false
}

// CompletionFunc returns the completion function of the flag
func (flag EnumSliceFlag) CompletionFunc() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return flag.Allowed, cobra.ShellCompDirectiveNoFileComp
	}
}
