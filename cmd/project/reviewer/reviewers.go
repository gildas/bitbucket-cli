package reviewer

type Reviewers []Reviewer

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (reviewers Reviewers) GetHeader() []string {
	return Reviewer{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (reviewers Reviewers) GetRowAt(index int) []string {
	if index < 0 || index >= len(reviewers) {
		return []string{}
	}
	return reviewers[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (reviewers Reviewers) Size() int {
	return len(reviewers)
}
