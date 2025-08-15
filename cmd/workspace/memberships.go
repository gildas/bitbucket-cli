package workspace

import "github.com/spf13/cobra"

type Memberships []Membership

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (memberships Memberships) GetHeaders(cmd *cobra.Command) []string {
	return Membership{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (memberships Memberships) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(memberships) {
		return []string{}
	}
	return memberships[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (memberships Memberships) Size() int {
	return len(memberships)
}
