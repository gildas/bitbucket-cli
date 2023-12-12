package common

// Tableable is an interface for objects that can be printed as a table
type Tableable interface {
	GetHeader(short bool) []string    // GetHeader retrieves the headers to show, if short is true get the short version
	GetRow(headers []string) []string // GetRow retrieves the row to show for the given headers
}

// Tableables is an interface for array of objects that can be printed as a table
type Tableables interface {
	GetHeader() []string         // GetHeader retrieves the headers to show
	GetRowAt(index int) []string // GetRow retrieves the row to show for the given headers
	Size() int                   // Size gets the number of elements
}
