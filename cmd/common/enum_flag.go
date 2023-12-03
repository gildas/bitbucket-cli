package common

import (
	"strings"

	"github.com/gildas/go-errors"
)

type EnumFlag struct {
	Allowed []string
	Value   string
}

// Type returns the type of the flag
func (flag EnumFlag) Type() string {
	return "string"
}

// String returns the string representation of the flag
func (flag EnumFlag) String() string {
	return flag.Value
}

// Set sets the flag value
func (flag *EnumFlag) Set(value string) error {
	for _, allowed := range flag.Allowed {
		if value == allowed {
			flag.Value = value
			return nil
		}
	}
	return errors.ArgumentInvalid.With("value", value, strings.Join(flag.Allowed, ", "))
}
