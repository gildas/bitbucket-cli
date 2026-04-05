package commit

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <commit-hash>",
	Aliases:           []string{"show", "describe"},
	Short:             "get a commit",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValiAdrgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
	Columns    *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get a commit from. Defaults to the current repository")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValiAdrgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	names, err := GetCommitHashes(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(names, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	commitName := args[0]

	if commitName == "" {
		log.Infof("No commit hash provided, getting the latest commit")
		latestCommit, err := GetLatestCommit(log.ToContext(cmd.Context()), cmd)
		if err != nil {
			return err
		}
		commitName = latestCommit.Hash
	}

	log.Infof("Displaying commit %s", commitName)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing commit %s", commitName)) {
		return nil
	}
	// Unlike what is mentioned in the official Bitbucket API:
	//   https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-get
	// the commit endpoint returns a list of commits in a struct.
	var commits struct {
		Values []Commit `json:"values" mapstructure:"values"`
	}
	err = profile.Get(log.ToContext(cmd.Context()), cmd, "commits/"+commitName, &commits)
	if err != nil {
		return err
	}
	if len(commits.Values) == 0 {
		log.Infof("No commit found with hash %s", commitName)
		return
	}
	log.Record("commit", commits.Values[0]).Debugf("Commit %s retrieved successfully", commitName)
	return profile.Print(cmd.Context(), cmd, commits.Values[0])
}
