package sshkey

type SSHKeys []SSHKey

// GetHeader gets the header for a table
//
// implements common.Tableables
func (keys SSHKeys) GetHeader() []string {
	return SSHKey{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (keys SSHKeys) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(keys) {
		return []string{}
	}
	return keys[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (keys SSHKeys) Size() int {
	return len(keys)
}
