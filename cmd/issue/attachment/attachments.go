package attachment

type Attachments []Attachment

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (comments Attachments) GetHeader() []string {
	return Attachment{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (comments Attachments) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(comments) {
		return []string{}
	}
	return comments[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (comments Attachments) Size() int {
	return len(comments)
}
