package project

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <project-key>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a project by its <project-key>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Workspace *flags.EnumFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	getCmd.Flags().Var(getOptions.Workspace, "workspace", "Workspace to get projects from")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Workspace.CompletionFunc("workspace"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	keys, err := GetProjectKeys(cmd.Context(), cmd, args, toComplete)
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return common.FilterValidArgs(keys, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, err := GetWorkspace(cmd, profile)
	if err != nil {
		return err
	}

	log.Infof("Displaying project %s", args[0])
	var project Project

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects/%s", workspace, args[0]),
		&project,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, project)
}
