package reviewer

import (
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
	Short: "list all reviewers",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Workspace *flags.EnumFlag
	Project   *flags.EnumFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	listOptions.Project = flags.NewEnumFlagWithFunc("", GetProjectKeys)
	listCmd.Flags().Var(listOptions.Workspace, "workspace", "Workspace to list reviewers from")
	listCmd.Flags().Var(listOptions.Project, "project", "Project Key to list reviewers from")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.Workspace.CompletionFunc("workspace"))
	_ = getCmd.RegisterFlagCompletionFunc(listOptions.Project.CompletionFunc("project"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	currentProfile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	workspace, project, err := GetWorkspaceAndProject(cmd, currentProfile)
	if err != nil {
		return err
	}

	log.Infof("Listing all reviewers")
	reviewers, err := GetDefaultReviewers(cmd.Context(), cmd, workspace, project)
	if err != nil {
		return err
	}
	if len(reviewers) == 0 {
		log.Infof("No reviewer found")
		return nil
	}
	core.Sort(reviewers, func(a, b Reviewer) bool {
		return strings.Compare(strings.ToLower(a.User.Username), strings.ToLower(b.User.Username)) == -1
	})
	return profile.Current.Print(cmd.Context(), cmd, Reviewers(reviewers))
}
