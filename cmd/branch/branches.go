package branch

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Branches []Branch

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (branches Branches) GetHeaders(cmd *cobra.Command) []string {
	return Branch{}.GetHeaders(cmd)
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
	uripath := "refs/branches"
	if cmd != nil && cmd.Flag("query") != nil && cmd.Flag("query").Changed {
		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return []Branch{}, err
		}
		uripath = fmt.Sprintf("%s?q=%s", uripath, url.QueryEscape(query))
	}
	return profile.GetAll[Branch](context, cmd, uripath)
}

// GetBranchNames gets the branch names of a repository
func GetBranchNames(context context.Context, cmd *cobra.Command, args []string, toComplete string) (names []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getbranches")
	log.Infof("Getting branches for profile %v", profile.Current)
	branches, err := GetBranches(context, cmd)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, err
	}
	names = core.Map(branches, func(branch Branch) string { return branch.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return common.FilterValidArgs(names, args, toComplete), nil
}
