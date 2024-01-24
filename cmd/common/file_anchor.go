package common

import (
	"strconv"
	"strings"
)

type FileAnchor struct {
	From uint64 `json:"from,omitempty" mapstructure:"from"`
	To   uint64 `json:"to,omitempty"   mapstructure:"to"`
	Path string `json:"path"           mapstructure:"path"`
}

// String gets a string representation of this FileAnchor
//
// implements fmt.Stringer
func (anchor FileAnchor) String() string {
	var value strings.Builder

	value.WriteString(anchor.Path)
	if anchor.From > 0 {
		value.WriteString(":")
		value.WriteString(strconv.FormatUint(anchor.From, 10))
		if anchor.To > 0 {
			value.WriteString("-")
			value.WriteString(strconv.FormatUint(anchor.To, 10))
		}
	} else if anchor.To > 0 {
		value.WriteString(":")
		value.WriteString(strconv.FormatUint(anchor.To, 10))
	}
	return value.String()
}
