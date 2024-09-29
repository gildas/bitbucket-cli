package workspace

import (
	"context"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <workspace-slug-or-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a workspace by its <workspace-slug-or-id>. With the --members flag, it will display the members of the workspace. With the --member flag, it will display workspaces for the given user.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getVAlidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Member      string
	WithMembers bool
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Member, "member", "", "Get a workspace member")
	getCmd.Flags().BoolVar(&getOptions.WithMembers, "members", false, "List the workspace members")
	getCmd.MarkFlagsMutuallyExclusive("member", "members")
}

func getVAlidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	slugs, err := GetWorkspaceSlugs(cmd.Context(), cmd, args)
	if err != nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	return slugs, cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if getOptions.WithMembers {
		log.Infof("Displaying workspace %s members", args[0])
		members, err := Workspace{Slug: args[0]}.GetMembers(cmd.Context(), cmd)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			log.Infof("No member found")
			return nil
		}
		return profile.Current.Print(cmd.Context(), cmd, Members(members))
	}

	if len(getOptions.Member) != 0 {
		log.Infof("Displaying workspace %s member %s", args[0], getOptions.Member)
		member, err := getWorkspaceMember(cmd.Context(), cmd, profile.Current, args[0], getOptions.Member)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get workspace member %s: %s\n", getOptions.Member, err)
			os.Exit(1)
		}
		return profile.Current.Print(cmd.Context(), cmd, member)
	}

	log.Infof("Displaying workspace %s", args[0])
	workspace, err := GetWorkspace(cmd.Context(), cmd, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get workspace %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, workspace)
}

func getWorkspaceMember(context context.Context, cmd *cobra.Command, profile *profile.Profile, workspace string, member string) (*Member, error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "get")

	if profile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying workspace %s member %s", workspace, member)
	var result Member

	err := profile.Get(
		log.ToContext(context),
		cmd,
		fmt.Sprintf("/workspaces/%s/members/%s", workspace, member),
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
