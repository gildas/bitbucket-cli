package pullrequest

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:               "approve [flags] <pullrequest-id>",
	Short:             "approve a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: approveValidArgs,
	RunE:              approveProcess,
}

var approveOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(approveCmd)

	approveCmd.Flags().StringVar(&approveOptions.Repository, "repository", "", "Repository to approve pullrequest from. Defaults to the current repository")
}

func approveValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "validargs")
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
	log.Debugf("Fetched %d pullrequest ids", len(ids))
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func approveProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "approve")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Approving pullrequest %s", args[0]) {
		return nil
	}
	var participant user.Participant

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s/approve", args[0]),
		nil,
		&participant,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to approve pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, participant)
}
