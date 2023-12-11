package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show", "info", "display"},
	Short:   "get a workspace",
	Args:    cobra.ExactArgs(1),
	RunE:    getProcess,
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

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if getOptions.WithMembers {
		log.Infof("Displaying workspace %s members", args[0])
		members, err := profile.GetAll[Member](
			cmd.Context(),
			profile.Current,
			"",
			fmt.Sprintf("/workspaces/%s/members", args[0]),
		)
		if err != nil {
			return err
		}
		if len(members) == 0 {
			log.Infof("No member found")
			return nil
		}
		payload, _ := json.MarshalIndent(members, "", "  ")
		fmt.Println(string(payload))
		return nil
	}

	if len(getOptions.Member) != 0 {
		log.Infof("Displaying workspace %s member %s", args[0], getOptions.Member)
		member, err := getWorkspaceMember(cmd.Context(), profile.Current, args[0], getOptions.Member)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get workspace member %s: %s\n", getOptions.Member, err)
			os.Exit(1)
		}

		payload, _ := json.MarshalIndent(member, "", "  ")
		fmt.Println(string(payload))
		return nil
	}

	log.Infof("Displaying workspace %s", args[0])
	workspace, err := getWorkspace(cmd.Context(), profile.Current, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get workspace %s: %s\n", args[0], err)
		os.Exit(1)
	}

	payload, _ := json.MarshalIndent(workspace, "", "  ")
	fmt.Println(string(payload))
	return nil
}

func getWorkspaceMember(context context.Context, profile *profile.Profile, workspace string, member string) (*Member, error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "get")

	if profile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying workspace %s member %s", workspace, member)
	var result Member

	err := profile.Get(
		log.ToContext(context),
		"",
		fmt.Sprintf("/workspaces/%s/members/%s", workspace, member),
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getWorkspace(context context.Context, profile *profile.Profile, workspace string) (*Workspace, error) {
	log := logger.Must(logger.FromContext(context)).Child("workspace", "get")

	if profile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	log.Infof("Displaying workspace %s", workspace)
	var result Workspace

	err := profile.Get(
		log.ToContext(context),
		"",
		fmt.Sprintf("/workspaces/%s", workspace),
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
