package commit

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Commits []Commit

// GetHeaders gets the header for a table
//
// implements common.Tableables
func (commits Commits) GetHeaders(cmd *cobra.Command) []string {
	return Commit{}.GetHeaders(cmd)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (commits Commits) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(commits) {
		return []string{}
	}
	return commits[index].GetRow(headers)
}

// Size gets the number of elements
//
// implements common.Tableables
func (commits Commits) Size() int {
	return len(commits)
}

// GetCommits gets the commits of a repository
func GetCommits(context context.Context, cmd *cobra.Command) (commits []Commit, err error) {
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return []Commit{}, err
	}
	uripath := repository.GetPath("commits")
	if cmd != nil && cmd.Flag("query") != nil && cmd.Flag("query").Changed {
		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return []Commit{}, err
		}
		uripath = fmt.Sprintf("%s?q=%s", uripath, url.QueryEscape(query))
	}
	if cmd != nil && cmd.Flag("include") != nil && cmd.Flag("include").Changed {
		include, err := cmd.Flags().GetStringSlice("include")
		if err != nil {
			return []Commit{}, err
		}
		if !strings.Contains(uripath, "?") {
			uripath = fmt.Sprintf("%s?include=%s", uripath, url.QueryEscape(include[0]))
			include = include[1:]
		}
		for _, hash := range include {
			uripath = fmt.Sprintf("%s&include=%s", uripath, url.QueryEscape(hash))
		}
	}
	if cmd != nil && cmd.Flag("exclude") != nil && cmd.Flag("exclude").Changed {
		exclude, err := cmd.Flags().GetStringSlice("exclude")
		if err != nil {
			return []Commit{}, err
		}
		if !strings.Contains(uripath, "?") {
			uripath = fmt.Sprintf("%s?exclude=%s", uripath, url.QueryEscape(exclude[0]))
			exclude = exclude[1:]
		}
		for _, hash := range exclude {
			uripath = fmt.Sprintf("%s&exclude=%s", uripath, url.QueryEscape(hash))
		}
	}
	return profile.GetAll[Commit](context, cmd, uripath)
}

// GetCommitsWithPrefix gets the commits of a repository with a prefix
func GetCommitsWithPrefix(context context.Context, cmd *cobra.Command, prefix string) (commits []Commit, err error) {
	if len(prefix) == 0 {
		return GetCommits(context, cmd)
	}
	repository, err := repository.GetRepository(context, cmd)
	if err != nil {
		return []Commit{}, err
	}
	uripath := repository.GetPath("commits", fmt.Sprintf("?q=hash~\"%s\"", url.QueryEscape(prefix)))
	return profile.GetAll[Commit](context, cmd, uripath)
}

// GetCommitHashes gets the commit hashes of a repository
func GetCommitHashes(context context.Context, cmd *cobra.Command, args []string, toComplete string) (hashes []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getcommits")
	log.Infof("Getting commits for profile %v", profile.Current)
	commits, err := GetCommitsWithPrefix(context, cmd, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, err
	}
	hashes = core.Map(commits, func(commit Commit) string { return commit.Hash })
	core.Sort(hashes, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return common.FilterValidArgs(hashes, args, toComplete), nil
}
