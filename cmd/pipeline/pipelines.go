package pipeline

import "github.com/spf13/cobra"

// Pipelines is a collection of Pipeline
type Pipelines []Pipeline

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (pipelines Pipelines) GetHeaders(cmd *cobra.Command) []string {
	return Pipeline{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (pipelines Pipelines) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(pipelines) {
		return []string{}
	}
	return pipelines[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (pipelines Pipelines) Size() int {
	return len(pipelines)
}
