package attachment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:               "download [flags] <path>",
	Aliases:           []string{"get", "fetch"},
	Short:             "download an issue attachment by its <path>.",
	ValidArgsFunction: downloadValidArgs,
	Args:              cobra.ExactArgs(1),
	RunE:              downloadProcess,
}

var downloadOptions struct {
	IssueID     *flags.EnumFlag
	Repository  string
	Destination string
	Progress    bool
}

func init() {
	Command.AddCommand(downloadCmd)

	downloadOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	downloadCmd.Flags().StringVar(&downloadOptions.Repository, "repository", "", "Repository to get an issue attachment from. Defaults to the current repository")
	downloadCmd.Flags().Var(downloadOptions.IssueID, "issue", "Issue to get attachments from")
	downloadCmd.Flags().StringVar(&downloadOptions.Destination, "destination", "", "Destination folder to download the attachment to. Defaults to the current folder")
	downloadCmd.Flags().BoolVar(&downloadOptions.Progress, "progress", false, "Show progress")
	_ = downloadCmd.MarkFlagRequired("issue")
	_ = downloadCmd.RegisterFlagCompletionFunc(downloadOptions.IssueID.CompletionFunc("issue"))
	_ = downloadCmd.MarkFlagDirname("destination")
}

func downloadValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	attachmentNames, err := GetAttachmentNames(cmd.Context(), cmd, deleteOptions.IssueID.Value)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(attachmentNames, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func downloadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "download")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Downloading attachment %s from issue %s to %s", args[0], downloadOptions.IssueID, downloadOptions.Destination) {
		err := profile.Download(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/attachments/%s", downloadOptions.IssueID.Value, args[0]),
			downloadOptions.Destination,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to download attachment %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return nil
}
