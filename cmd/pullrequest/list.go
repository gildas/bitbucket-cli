package pullrequest

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all pullrequests",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	State      common.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.State = common.EnumFlag{Allowed: []string{"all", "open", "closed", "merged", "declined"}, Value: "all"}
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequests from. Defaults to the current repository")
	listCmd.Flags().Var(&listOptions.State, "state", "Pull request state to fetch. Defaults to \"all\"")
	_ = listCmd.RegisterFlagCompletionFunc("state", listOptions.State.CompletionFunc())
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Listing all pull requests for repository: %s with profile %s", listOptions.Repository, profile.Current)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"pullrequests?state="+listOptions.State.String(),
	)
	if err != nil {
		return err
	}
	if len(pullrequests) == 0 {
		log.Infof("No pullrequest found")
		return
	}
	payload, _ := json.MarshalIndent(pullrequests, "", "  ")
	fmt.Println(string(payload))
	return nil
}
