package artifact

type Artifacts []Artifact

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (artifacts Artifacts) GetHeader() []string {
	return Artifact{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (artifacts Artifacts) GetRowAt(index int) []string {
	if index < 0 || index >= len(artifacts) {
		return []string{}
	}
	return artifacts[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (artifacts Artifacts) Size() int {
	return len(artifacts)
}
