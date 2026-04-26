package attachment

import (
	"fmt"
	"net/url"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all issue attachments",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	Query      string
	IssueID    *flags.EnumFlag
	Columns    *flags.EnumSliceFlag
	SortBy     *flags.EnumFlag
	PageLength int
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.IssueID = flags.NewEnumFlagWithFunc("", GetIssueIDs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns.Columns()...)
	listOptions.SortBy = flags.NewEnumFlag(columns.Sorters()...)
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issue attachments from. Defaults to the current repository.\nExpected format: <workspace>/<repository> or <repository>.\nIf only <repository> is given, the profile's default workspace is used.")
	listCmd.Flags().Var(listOptions.IssueID, "issue", "Issue to list attachments from")
	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter attachments")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	listCmd.Flags().IntVar(&listOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
	_ = listCmd.MarkFlagRequired("issue")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.IssueID.CompletionFunc("issue"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	var uripath string

	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("issues/%s/attachments?q=%s", listOptions.IssueID.Value, url.QueryEscape(listOptions.Query))
	} else {
		uripath = fmt.Sprintf("issues/%s/attachments", listOptions.IssueID.Value)
	}

	log.Infof("Listing all attachments from repository %s", listOptions.Repository)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing attachments for issue %s in repository %s", listOptions.IssueID.Value, listOptions.Repository)) {
		return nil
	}

	attachments, err := profile.GetAll[Attachment](cmd.Context(), cmd, uripath)
	if err != nil {
		return err
	}
	if len(attachments) == 0 {
		log.Infof("No attachment found")
		return nil
	}
	core.Sort(attachments, columns.SortBy(listOptions.SortBy.Value))
	return profile.Current.Print(cmd.Context(), cmd, Attachments(attachments))
}
