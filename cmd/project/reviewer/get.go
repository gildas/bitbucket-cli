package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show", "info", "display"},
	Short:   "get a reviewer",
	Args:    cobra.ExactArgs(1),
	RunE:    getProcess,
}

var getOptions struct {
	Workspace common.RemoteValueFlag
	Project   string
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	getCmd.Flags().Var(&getOptions.Workspace, "workspace", "Workspace to get reviewers from")
	getCmd.Flags().StringVar(&getOptions.Project, "project", "", "Project Key to get reviewers from")
	_ = getCmd.MarkFlagRequired("workspace")
	_ = getCmd.MarkFlagRequired("project")
	_ = getCmd.RegisterFlagCompletionFunc("workspace", getOptions.Workspace.CompletionFunc())
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying reviewer %s", args[0])
	var user user.User

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		"",
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", getOptions.Workspace, getOptions.Project, args[0]),
		&user,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add reviewer: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), user)
}
