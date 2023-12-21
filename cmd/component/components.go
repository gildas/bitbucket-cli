package component

type Components []Component

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (components Components) GetHeader() []string {
	return Component{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (components Components) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(components) {
		return []string{}
	}
	return components[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (components Components) Size() int {
	return len(components)
}
