package key

type GPGKeys []GPGKey

// GetHeader gets the header for a table
//
// implements common.Tableables
func (keys GPGKeys) GetHeader() []string {
	return GPGKey{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (keys GPGKeys) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(keys) {
		return []string{}
	}
	return keys[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (keys GPGKeys) Size() int {
	return len(keys)
}
