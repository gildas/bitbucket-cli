package branch

import (
	"context"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/spf13/cobra"
)

type Branches []Branch

// GetHeader gets the header for a table
//
// implements common.Tableables
func (branches Branches) GetHeader() []string {
	return Branch{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (branches Branches) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(branches) {
		return []string{}
	}
	return branches[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (branches Branches) Size() int {
	return len(branches)
}

// GetBranches gets the branches of a repository
func GetBranches(context context.Context, cmd *cobra.Command) (branches []Branch, err error) {
	return profile.GetAll[Branch](context, cmd, "refs/branches")
}

// GetBranchNames gets the branch names of a repository
func GetBranchNames(context context.Context, cmd *cobra.Command) (names []string, err error) {
	branches, err := GetBranches(context, cmd)
	if err != nil {
		return []string{}, err
	}
	names = core.Map(branches, func(branch Branch) string { return branch.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return names, nil
}
