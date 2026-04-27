package commit

import (
	"fmt"
	"io"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var patchCmd = &cobra.Command{
	Use:               "patch <commit-hash> <commit-hash>",
	Short:             "show the patch between two commits",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: validPatchArgs,
	RunE:              patchProcess,
}

func init() {
	Command.AddCommand(patchCmd)
}

func validPatchArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "patch")
	if len(args) == 0 && toComplete == "" || len(args) > 2 {
		log.Debugf("No args or too many args for completion: %v, toComplete: %s", args, toComplete)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	log.Debugf("Getting commit hashes for completion with args: %v and toComplete: %s", args, toComplete)
	names, err := GetCommitHashes(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(names, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func patchProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "patch")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Debugf("Displaying patch between commit %s and %s", args[0], args[1])
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing patch between commit %s and %s", args[0], args[1])) {
		return nil
	}
	repository, err := repository.GetRepository(cmd.Context(), cmd)
	if err != nil {
		return err
	}
	patch, err := profile.GetRaw(log.ToContext(cmd.Context()), cmd, repository.GetPath("patch", fmt.Sprintf("%s..%s", args[0], args[1])))
	if err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, patch)
	return err
}
