package pullrequest

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a profile",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get pullrequest from. Defaults to the current repository")
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Infof("Getting open pullrequests for repository %s", approveOptions.Repository)
	pullrequests, err := profile.GetAll[PullRequest](
		log.ToContext(cmd.Context()),
		profile.Current,
		listOptions.Repository,
		"pullrequests?state=OPEN",
	)
	if err != nil {
		log.Errorf("Failed to get pullrequests for repository %s", unapproveOptions.Repository, err)
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	var result []string
	for _, pullrequest := range pullrequests {
		result = append(result, fmt.Sprintf("%d", pullrequest.ID))
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying pull request %s", args[0])
	var pullrequest PullRequest

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		getOptions.Repository,
		fmt.Sprintf("pullrequests/%s", args[0]),
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}

	payload, _ := json.MarshalIndent(pullrequest, "", "  ")
	fmt.Println(string(payload))
	return nil
}
