package issue

type IssueChanges []IssueChange

// GetHeader gets the headers for the list command
//
// implements common.Tableables
func (issueChanges IssueChanges) GetHeader() []string {
	return IssueChange{}.GetHeader(false)
}

// GetRowAt gets the row for the list command
//
// implements common.Tableables
func (issueChanges IssueChanges) GetRowAt(index int) []string {
	if index < 0 || index >= len(issueChanges) {
		return []string{}
	}
	return issueChanges[index].GetRow(nil)
}

// Size gets the number of elements
//
// implements common.Tableables
func (issueChanges IssueChanges) Size() int {
	return len(issueChanges)
}
