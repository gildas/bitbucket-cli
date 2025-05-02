package prcommon

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PullRequestID struct {
	ID int `json:"id" mapstructure:"id"`
}

// GetPullRequestIDsWithState gets the pullrequest Ids for completion for a given state
func GetPullRequestIDsWithState(context context.Context, cmd *cobra.Command, state string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getpullrequests")

	log.Infof("Getting %s pullrequests", state)
	pullrequests, err := profile.GetAll[PullRequestID](
		log.ToContext(context),
		cmd,
		fmt.Sprintf("pullrequests?state=%s", state),
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
	ids, err = GetPullRequestIDsWithState(context, cmd, "OPEN")
	if err != nil {
		return []string{}, err
	}
	if len(ids) > 0 {
		return ids, nil
	}
	return GetPullRequestIDsWithState(context, cmd, "ALL")
}
