package artifact

import "github.com/spf13/cobra"

type Artifacts []Artifact

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (artifacts Artifacts) GetHeaders(cmd *cobra.Command) []string {
	return Artifact{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (artifacts Artifacts) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(artifacts) {
		return []string{}
	}
	return artifacts[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (artifacts Artifacts) Size() int {
	return len(artifacts)
}
