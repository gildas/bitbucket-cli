package repository

import (
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all public repositories",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Role      common.EnumFlag
	Workspace common.RemoteValueFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Role = common.EnumFlag{Allowed: []string{"owner", "admin", "contributor", "member", "all"}, Value: "owner"}
	listOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	listCmd.Flags().Var(&listOptions.Role, "role", "Role of the user in the repository")
	listCmd.Flags().Var(&listOptions.Workspace, "workspace", "Workspace to list repositories from")
	_ = listCmd.RegisterFlagCompletionFunc("workspace", listOptions.Workspace.CompletionFunc())
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	filter := ""
	if listOptions.Role.Value != "all" {
		filter = "?role=" + listOptions.Role.Value
	}

	workspace := ""
	if len(listOptions.Workspace.Value) > 0 {
		workspace = "/" + listOptions.Workspace.Value
	}

	log.Infof("Listing all repositories, workspace %s, role %s", listOptions.Workspace, listOptions.Role)
	repositories, err := profile.GetAll[Repository](
		cmd.Context(),
		cmd,
		profile.Current,
		"/repositories"+workspace+filter,
	)
	if err != nil {
		return err
	}
	if len(repositories) == 0 {
		log.Infof("No repository found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), Repositories(repositories))
}
