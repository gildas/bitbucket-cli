package repository

type Repositories []Repository

// GetHeader gets the header for a table
//
// implements common.Tableables
func (repositories Repositories) GetHeader() []string {
	return Repository{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (repositories Repositories) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(repositories) {
		return []string{}
	}
	return repositories[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (repositories Repositories) Size() int {
	return len(repositories)
}
