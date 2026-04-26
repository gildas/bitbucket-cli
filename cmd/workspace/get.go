package workspace

import (
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <workspace-slug-or-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a workspace by its <workspace-slug-or-id> or the current workspace by default. With the --members flag, it ill display the members of the workspace. With the --member flag, it will display workspaces for the given user.",
	Args:              cobra.RangeArgs(0, 1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Member      string
	WithMembers bool
	Columns     *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(getCmd)

	getOptions.Columns = flags.NewEnumSliceFlag(columns.Columns()...)
	getCmd.Flags().StringVar(&getOptions.Member, "member", "", "Get a workspace member")
	getCmd.Flags().BoolVar(&getOptions.WithMembers, "members", false, "List the workspace members")
	getCmd.Flags().Var(getOptions.Columns, "columns", "Comma-separated list of columns to display")
	getCmd.MarkFlagsMutuallyExclusive("member", "members")
	_ = getCmd.RegisterFlagCompletionFunc(getOptions.Columns.CompletionFunc("columns"))
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	slugs, err := GetWorkspaceSlugs(cmd.Context(), cmd)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return common.FilterValidArgs(slugs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	var workspace *Workspace

	if len(args) == 0 {
		if workspace, err = GetWorkspace(cmd.Context(), cmd); err != nil {
			return errors.Join(
				errors.Errorf("Failed to get current workspace"),
				err,
			)
		}
	} else {
		if workspace, err = GetWorkspaceBySlugOrID(cmd.Context(), cmd, args[0]); err != nil {
			return errors.Join(
				errors.Errorf("Failed to get workspace %s", args[0]),
				err,
			)
		}
	}

	if getOptions.WithMembers {
		log.Infof("Displaying workspace %s members", workspace.Slug)
		if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing workspace %s members", workspace.Slug)) {
			return nil
		}
		members, err := workspace.GetMembers(cmd.Context(), cmd)
		if err != nil {
			return errors.Join(
				errors.Errorf("Failed to get members of workspace %s", workspace.Slug),
				err,
			)
		}
		if len(members) == 0 {
			log.Infof("No member found")
			return nil
		}
		return profile.Print(cmd.Context(), cmd, Members(members))
	}

	if len(getOptions.Member) != 0 {
		log.Infof("Displaying workspace %s member %s", workspace.Slug, getOptions.Member)
		if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing workspace %s member %s", workspace.Slug, getOptions.Member)) {
			return nil
		}
		member, err := workspace.GetMember(cmd.Context(), cmd, profile, getOptions.Member)
		if err != nil {
			return errors.Join(
				errors.Errorf("Failed to get workspace member %s of workspace %s", getOptions.Member, workspace.Slug),
				err,
			)
		}
		return profile.Print(cmd.Context(), cmd, member)
	}

	log.Infof("Displaying workspace %s", workspace.Slug)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, fmt.Sprintf("Showing workspace %s", workspace.Slug)) {
		return nil
	}
	return profile.Print(cmd.Context(), cmd, workspace)
}
