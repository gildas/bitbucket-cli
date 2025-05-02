package workspace

import (
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
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
	WithMembership bool
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listOptions.WithMembership, "membership", false, "List also the workspace memberships of the current user")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if listOptions.WithMembership {
		log.Infof("Listing all workspace memberships for current user")
		memberships, err := profile.GetAll[Membership](cmd.Context(), cmd, "/user/permissions/workspaces")
		if err != nil {
			return err
		}
		if len(memberships) == 0 {
			log.Infof("No workspace found")
			return nil
		}
		return profile.Current.Print(cmd.Context(), cmd, Memberships(memberships))
	}

	log.Infof("Listing all workspaces")
	workspaces, err := profile.GetAll[Workspace](
		cmd.Context(),
		cmd,
		"/workspaces",
	)
	if err != nil {
		return err
	}
	if len(workspaces) == 0 {
		log.Infof("No workspace found")
		return nil
	}
	core.Sort(workspaces, func(a, b Workspace) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return profile.Current.Print(cmd.Context(), cmd, Workspaces(workspaces))
}
