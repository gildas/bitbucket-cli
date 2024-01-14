package comment

type Comments []Comment

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (comments Comments) GetHeader() []string {
	short := true
	for _, comment := range comments {
		if !comment.UpdatedOn.IsZero() {
			short = false
			break
		}
	}
	return Comment{}.GetHeader(short)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Comments) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Comments) Size() int {
	return len(comments)
}
