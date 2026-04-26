package workspace

import (
	"net/url"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all workspaces for the current user",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Query      string
	PageLength int
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listOptions.Query, "query", "", "Query string to filter workspaces")
	listCmd.Flags().IntVar(&listOptions.PageLength, "page-length", 0, "Number of items per page to retrieve from Bitbucket. Default is the profile's default page length")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	query := url.Values{}
	if listOptions.Query != "" {
		query.Add("q", listOptions.Query)
	}

	log.Infof("Listing all workspaces")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Showing workspaces") {
		return nil
	}

	workspaces, err := GetWorkspacesWithQuery(cmd.Context(), cmd, query)
	if err != nil {
		return errors.Join(errors.New("failed to retrieve workspaces"), err)
	}
	if len(workspaces) == 0 {
		log.Infof("No workspace found")
		return nil
	}
	log.Debugf("Found %d workspace accesses", len(workspaces))
	return profile.Current.Print(cmd.Context(), cmd, workspaces)
}
