package workspace

type Members []Member

// GetHeader gets the header for a table
//
// implements common.Tableables
func (members Members) GetHeader() []string {
	return Member{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (members Members) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(members) {
		return []string{}
	}
	return members[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (members Members) Size() int {
	return len(members)
}
