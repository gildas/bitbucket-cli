package project

import "github.com/spf13/cobra"

type Projects []Project

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (projects Projects) GetHeaders(cmd *cobra.Command) []string {
	return Project{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (projects Projects) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(projects) {
		return []string{}
	}
	return projects[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (projects Projects) Size() int {
	return len(projects)
}
