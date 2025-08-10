package common

import "github.com/spf13/cobra"

// Tableable is an interface for objects that can be printed as a table
type Tableable interface {
	GetHeaders(cmd *cobra.Command) []string // GetHeaders retrieves the headers to show, use --columns flag or gives a default list
	GetRow(headers []string) []string       // GetRow retrieves the row to show for the given headers
}

// Tableables is an interface for array of objects that can be printed as a table
type Tableables interface {
	GetHeaders(cmd *cobra.Command) []string        // GetHeaders retrieves the headers to show
	GetRowAt(index int, headers []string) []string // GetRow retrieves the row to show for the given headers
	Size() int                                     // Size gets the number of elements
}
