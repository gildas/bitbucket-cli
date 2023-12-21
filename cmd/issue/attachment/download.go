package attachment

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:               "download",
	Aliases:           []string{"get", "fetch"},
	Short:             "download an issue attachment",
	ValidArgsFunction: downloadValidArgs,
	Args:              cobra.ExactArgs(1),
	RunE:              downloadProcess,
}

var downloadOptions struct {
	IssueID     common.RemoteValueFlag
	Repository  string
	Destination string
}

func init() {
	Command.AddCommand(downloadCmd)

	downloadOptions.IssueID = common.RemoteValueFlag{AllowedFunc: GetIssueIDs}
	downloadCmd.Flags().StringVar(&downloadOptions.Repository, "repository", "", "Repository to get an issue attachment from. Defaults to the current repository")
	downloadCmd.Flags().Var(&downloadOptions.IssueID, "issue", "Issue to get attachments from")
	downloadCmd.Flags().StringVar(&downloadOptions.Destination, "destination", "", "Destination folder to download the attachment to. Defaults to the current folder")
	_ = downloadCmd.MarkFlagRequired("issue")
	_ = downloadCmd.RegisterFlagCompletionFunc("issue", downloadOptions.IssueID.CompletionFunc())
	_ = downloadCmd.MarkFlagDirname("destination")
}

func downloadValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetAttachmentNames(cmd.Context(), cmd, profile.Current, downloadOptions.IssueID.Value), cobra.ShellCompDirectiveNoFileComp
}

func downloadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "download")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Downloading artifact %s to %s", args[0], downloadOptions.Destination)

	err := profile.Current.Download(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("issues/%s/attachments/%s", downloadOptions.IssueID.Value, args[0]),
		downloadOptions.Destination,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to download attachment %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return nil
}
