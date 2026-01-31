package user

import "github.com/spf13/cobra"

type Emails []Email

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (emails Emails) GetHeaders(cmd *cobra.Command) []string {
	return Email{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (emails Emails) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(emails) {
		return []string{}
	}
	return emails[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (emails Emails) Size() int {
	return len(emails)
}
