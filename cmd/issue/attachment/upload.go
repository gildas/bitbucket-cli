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

var uploadCmd = &cobra.Command{
	Use:     "upload [flags] <filename>",
	Aliases: []string{"add", "create"},
	Short:   "upload an artifact.",
	Args:    cobra.ExactArgs(1),
	RunE:    uploadProcess,
}

var uploadOptions struct {
	IssueID    *flags.EnumFlag
	Repository string
	Progress   bool
}

func init() {
	Command.AddCommand(uploadCmd)

	uploadOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	uploadCmd.Flags().StringVar(&uploadOptions.Repository, "repository", "", "Repository to upload issue attachments to. Defaults to the current repository")
	uploadCmd.Flags().Var(uploadOptions.IssueID, "issue", "Issue to upload attachments to")
	uploadCmd.Flags().BoolVar(&uploadOptions.Progress, "progress", false, "Show progress")
	_ = uploadCmd.MarkFlagRequired("issue")
	_ = uploadCmd.RegisterFlagCompletionFunc(uploadOptions.IssueID.CompletionFunc("issue"))
}

func uploadProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "upload")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	if common.WhatIf(log.ToContext(cmd.Context()), cmd, "Uploading attachment %s to issue %s", args[0], uploadOptions.IssueID) {
		err = profile.Upload(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s/attachments", uploadOptions.IssueID.Value),
			args[0],
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to upload attachment %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return nil
}
