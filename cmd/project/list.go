package project

import (
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all projects",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Workspace *flags.EnumFlag
	Columns   *flags.EnumSliceFlag
	SortBy    *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	listOptions.Columns = flags.NewEnumSliceFlagWithAllAllowed(columns...)
	listOptions.SortBy = flags.NewEnumFlag(sortBy...)
	listCmd.Flags().Var(listOptions.Workspace, "workspace", "Workspace to list projects from")
	listCmd.Flags().Var(listOptions.Columns, "columns", "Comma-separated list of columns to display")
	listCmd.Flags().Var(listOptions.SortBy, "sort", "Column to sort by")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Workspace.CompletionFunc("workspace"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Columns.CompletionFunc("columns"))
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.SortBy.CompletionFunc("sort"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	currentProfile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, err := GetWorkspace(cmd, currentProfile)
	if err != nil {
		return err
	}

	log.Infof("Listing all projects from workspace %s with profile %s", workspace, currentProfile)
	projects, err := profile.GetAll[Project](
		cmd.Context(),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects", workspace),
	)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		log.Infof("No project found")
		return nil
	}
	core.Sort(projects, func(a, b Project) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	})
	return currentProfile.Print(cmd.Context(), cmd, Projects(projects))
}
