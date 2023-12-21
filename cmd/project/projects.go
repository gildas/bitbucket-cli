package project

type Projects []Project

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (projects Projects) GetHeader() []string {
	return Project{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (projects Projects) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(projects) {
		return []string{}
	}
	return projects[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (projects Projects) Size() int {
	return len(projects)
}
