package pullrequest

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
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
	State      *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.State = flags.NewEnumFlag("all", "declined", "merged", "+open", "superseded")
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list pullrequests from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.State, "state", "Pull request state to fetch. Defaults to \"open\"")
	_ = listCmd.RegisterFlagCompletionFunc("state", listOptions.State.CompletionFunc("state"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	log.Infof("Listing %s pull requests for repository: %s", listOptions.State, listOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		cmd,
		"pullrequests/?state="+strings.ToUpper(listOptions.State.String()),
	)
	if err != nil {
		return err
	}
	if len(pullrequests) == 0 {
		log.Infof("No pullrequest found")
		return
	}
	return profile.Current.Print(cmd.Context(), cmd, PullRequests(pullrequests))
}
