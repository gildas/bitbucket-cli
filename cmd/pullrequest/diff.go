package pullrequest

import (
	"fmt"
	"io"
	"os"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/bitbucket-cli/cmd/profile"
	prcommon "github.com/gildas/bitbucket-cli/cmd/pullrequest/common"
	"github.com/gildas/bitbucket-cli/cmd/repository"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:               "diff [flags] <pullrequest-id>",
	Short:             "show the diff of a pull request by its <pullrequest-id>. If not provided, it will try to show the diff of the only open pullrequest.",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: validDiffArgs,
	RunE:              diffProcess,
}

var diffOptions struct {
	Stat bool
}

func init() {
	Command.AddCommand(diffCmd)

	diffCmd.Flags().BoolVar(&diffOptions.Stat, "stat", false, "show only the diffstat")
}

func validDiffArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	ids, err := prcommon.GetPullRequestIDsWithState(cmd.Context(), cmd, "OPEN")
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func diffProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "diff")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot show diff of Pull Request"), err)
	}

	pullRequestID, err := GetPullRequestIDFromArgs(cmd.Context(), cmd, repository, args)
	if err != nil {
		return errors.Join(errors.Errorf("Cannot show diff of Pull Request"), err)
	}

	log.Debugf("Displaying diff for Pull Request ID: %s", pullRequestID)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing diff for Pull Request ID %s", pullRequestID)) {
		return nil
	}

	uripath := repository.GetPath("pullrequests", pullRequestID, "diff")
	if diffOptions.Stat {
		uripath = repository.GetPath("pullrequests", pullRequestID, "diffstat")
	}

	diff, err := profile.GetRaw(log.ToContext(cmd.Context()), cmd, uripath)
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, diff)
	return err
}
