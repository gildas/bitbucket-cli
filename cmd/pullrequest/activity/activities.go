package activity

type Activities []Activity

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (activities Activities) GetHeader() []string {
	return Activity{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (activities Activities) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(activities) {
		return []string{}
	}
	return activities[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (activities Activities) Size() int {
	return len(activities)
}
