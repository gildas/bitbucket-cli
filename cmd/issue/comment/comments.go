package comment

type Comments []Comment

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (comments Comments) GetHeader() []string {
	return Comment{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Comments) GetRowAt(index int) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Comments) Size() int {
	return len(comments)
}
