package reviewer

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a reviewer",
	ValidArgsFunction: getValidArgs,
	Args:              cobra.ExactArgs(1),
	RunE:              getProcess,
}

var getOptions struct {
	Workspace *flags.EnumFlag
	Project   *flags.EnumFlag
	Columns   *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	getOptions.Project = flags.NewEnumFlagWithFunc("", GetProjectKeys)
	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().Var(getOptions.Workspace, "workspace", "Workspace to get reviewers from")
	getCmd.Flags().Var(getOptions.Project, "project", "Project Key to get reviewers from")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Workspace.CompletionFunc("workspace"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Project.CompletionFunc("project"))
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	userIDs, err := GetReviewerUserIDs(cmd.Context(), cmd, deleteOptions.Project.Value)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return userIDs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, project, err := GetWorkspaceAndProject(cmd, profile)
	if err != nil {
		return err
	}

	log.Infof("Displaying reviewer %s", args[0])
	var user user.User

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers/%s", workspace, project, args[0]),
		&user,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get reviewer: %s\n", err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, user)
}
