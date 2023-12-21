package issue

type Issues []Issue

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (issues Issues) GetHeader() []string {
	return Issue{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (issues Issues) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(issues) {
		return []string{}
	}
	return issues[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (issues Issues) Size() int {
	return len(issues)
}
