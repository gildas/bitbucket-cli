package issue

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:               "watch",
	Short:             "watch an issue",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: watchValidArgs,
	RunE:              watchProcess,
}

var watchOptions struct {
	Repository string
	Check      bool
}

func init() {
	Command.AddCommand(watchCmd)

	watchCmd.Flags().StringVar(&watchOptions.Repository, "repository", "", "Repository to watch an issue from. Defaults to the current repository")
	watchCmd.Flags().BoolVar(&watchOptions.Check, "check", false, "Check if the issue is watched")
}

func watchValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func watchProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "watch")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if watchOptions.Check {
		err = profile.Current.Get(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/watch", args[0]),
			nil,
		)
		return
	}

	log.Infof("watch for issue %s", args[0])
	err = profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/watch", args[0]),
		nil,
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to watch issue %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return
}
