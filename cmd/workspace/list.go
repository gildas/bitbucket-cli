package workspace

import (
	"fmt"
	"net/url"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
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

	uripath := "/user/workspaces"
	if len(listOptions.Query) > 0 {
		uripath = fmt.Sprintf("/user/workspaces?q=%s", url.QueryEscape(listOptions.Query))
	}

	log.Infof("Listing all workspaces")
	workspaceAccesses, err := profile.GetAll[WorkspaceAccess](cmd.Context(), cmd, uripath)
	if err != nil {
		return errors.Join(errors.New("failed to retrieve workspaces"), err)
	}
	if len(workspaceAccesses) == 0 {
		log.Infof("No workspace found")
		return nil
	}
	log.Debugf("Found %d workspace accesses", len(workspaceAccesses))
	workspaces := core.Map(workspaceAccesses, func(access WorkspaceAccess) WorkspaceBase { return access.Workspace })
	core.Sort(workspaces, func(a, b WorkspaceBase) bool { return a.Slug < b.Slug })
	return profile.Current.Print(cmd.Context(), cmd, WorkspaceBases(workspaces))
}
