package prcommon

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gildas/bitbucket-cli/cmd/profile"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PullRequestID struct {
	ID int `json:"id" mapstructure:"id"`
}

// GetPullRequestIDsWithState gets the pullrequest Ids for completion for a given state
func GetPullRequestIDsWithState(context context.Context, cmd *cobra.Command, state string) (ids []string, err error) {
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return []string{}, err
	}
	return GetPullRequestIDsFromRepositoryWithState(context, cmd, repository, state)
}

// GetPullRequestIDsFromRepositoryWithState gets the pullrequest Ids for completion for a given state and repository
func GetPullRequestIDsFromRepositoryWithState(context context.Context, cmd *cobra.Command, repository *repository.Repository, state string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getpullrequests")

	log.Infof("Getting %s pullrequests", state)
	pullrequests, err := profile.GetAll[PullRequestID](
		log.ToContext(context),
		cmd,
		repository.GetPath(fmt.Sprintf("pullrequests?state=%s", state)),
	)
	if err != nil {
		log.Errorf("Failed to get %s pullrequests", state, err)
		return []string{}, err
	}

	ids = core.Map(pullrequests, func(pullrequest PullRequestID) string { return fmt.Sprintf("%d", pullrequest.ID) })
	core.Sort(ids, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return ids, nil
}

// GetPullRequestIDs gets the IDs of the pullrequests
//
// First only the open pullrequests are fetched, if none are found, all pullrequests are fetched
func GetPullRequestIDs(context context.Context, cmd *cobra.Command, args []string, toComplete string) (ids []string, err error) {
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return []string{}, err
	}

	pullrequests, err := profile.GetAllWithLimit[PullRequestID](
		context,
		cmd,
		repository.GetPath("pullrequests?sort=-created_on"),
		1,
	)
	if err != nil {
		return []string{}, err
	}

	if len(pullrequests) > 0 {
		// return all the IDs from 1 to the latest pull request ID
		latestID := pullrequests[0].ID
		ids = make([]string, 0, latestID)
		for i := 1; i <= latestID; i++ {
			ids = append(ids, strconv.Itoa(i))
		}
		return ids, nil
	}
	return []string{}, nil
}
