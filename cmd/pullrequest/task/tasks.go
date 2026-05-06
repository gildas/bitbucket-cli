package task

import "github.com/spf13/cobra"

type Tasks []Task

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (tasks Tasks) GetHeaders(cmd *cobra.Command) []string {
	return Task{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (tasks Tasks) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(tasks) {
		return []string{}
	}
	return tasks[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (tasks Tasks) Size() int {
	return len(tasks)
}
