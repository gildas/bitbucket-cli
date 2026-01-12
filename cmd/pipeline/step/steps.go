package step

import "github.com/spf13/cobra"

type Steps []Step

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (steps Steps) GetHeaders(cmd *cobra.Command) []string {
	if len(steps) == 0 {
		return Step{}.GetHeaders(cmd)
	}
	return steps[0].GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (steps Steps) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(steps) {
		return []string{}
	}
	return steps[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (steps Steps) Size() int {
	return len(steps)
}
