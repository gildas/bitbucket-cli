package attachment

import "github.com/spf13/cobra"

type Attachments []Attachment

// GetHeaders gets the headers for the list command
//
// implements common.Tableables
func (comments Attachments) GetHeaders(cmd *cobra.Command) []string {
	return Attachment{}.GetHeaders(cmd)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Attachments) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Attachments) Size() int {
	return len(comments)
}
