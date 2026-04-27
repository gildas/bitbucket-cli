package reviewer

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Reviewer struct {
	Type         string    `json:"type" mapstructure:"type"`
	ReviewerType string    `json:"reviewer_type" mapstructure:"reviewer_type"`
	User         user.User `json:"user" mapstructure:"user"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "reviewer",
	Short: "Manage reviewers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reviewer requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Reviewer]{
	{Name: "user", DefaultSorter: true, Compare: func(a, b Reviewer) bool {
		return strings.Compare(strings.ToLower(a.User.Name), strings.ToLower(b.User.Name)) == -1
	}},
	{Name: "type", DefaultSorter: false, Compare: func(a, b Reviewer) bool {
		return strings.Compare(strings.ToLower(a.Type), strings.ToLower(b.Type)) == -1
	}},
	{Name: "reviewer_type", DefaultSorter: false, Compare: func(a, b Reviewer) bool {
		return strings.Compare(strings.ToLower(a.ReviewerType), strings.ToLower(b.ReviewerType)) == -1
	}},
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Type", "Reviewer Type", "User"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "type":
			row = append(row, reviewer.Type)
		case "reviewer type", "reviewer_type":
			row = append(row, reviewer.ReviewerType)
		case "user":
			row = append(row, reviewer.User.Name)
		}
	}
	return row
}

// Validate validates a Reviewer
func (reviewer *Reviewer) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// GetWorkspaceAndProject gets the workspace and project from the command
func GetWorkspaceAndProject(cmd *cobra.Command, profile *profile.Profile) (workspaceName, projectName string, err error) {
	workspaceName, err = workspace.GetWorkspaceName(cmd.Context(), cmd)
	if err != nil {
		return "", "", err
	}

	projectName = cmd.Flag("project").Value.String()
	if len(projectName) == 0 {
		projectName = profile.DefaultProject
		if len(projectName) == 0 {
			return "", "", errors.ArgumentMissing.With("project")
		}
	}
	return
}

// GetProjectName gets the project name from the command or profile
func GetProjectName(cmd *cobra.Command, profile *profile.Profile) (projectName string, err error) {
	projectName = cmd.Flag("project").Value.String()
	if len(projectName) == 0 {
		projectName = profile.DefaultProject
		if len(projectName) == 0 {
			return "", errors.ArgumentMissing.With("project")
		}
	}
	return
}

// GetProjectKeys gets the keys of the projects in the workspace given in the command
func GetProjectKeys(context context.Context, cmd *cobra.Command, args []string, toComplete string) (keys []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")

	workspace, err := workspace.GetWorkspace(cmd.Context(), cmd)
	if err != nil {
		log.Warnf("No workspace given")
		return
	}

	type Project struct {
		Key string `json:"key" mapstructure:"key"`
	}

	log.Infof("Getting all projects from workspace %s", workspace)
	projects, err := profile.GetAll[Project](context, cmd, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	return core.Map(projects, func(project Project) string {
		return project.Key
	}), nil
}

// GetReviewerIDs gets the IDs of the reviewers in the given workspace and project
func GetReviewerUserIDs(context context.Context, cmd *cobra.Command, project string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("reviewer", "getids")

	workspace, err := workspace.GetWorkspace(cmd.Context(), cmd)
	if err != nil {
		return []string{}, err
	}
	reviewers, err := profile.GetAll[Reviewer](context, cmd, fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers", workspace, project))
	if err != nil {
		log.Errorf("Failed to get reviewers", err)
		return
	}
	return core.Map(reviewers, func(reviewer Reviewer) string {
		return reviewer.User.ID.String()
	}), nil
}

// GetProjectDefaultReviewers gets the reviewers in the given workspace and project
func GetProjectDefaultReviewers(context context.Context, cmd *cobra.Command, project string) (reviewers []Reviewer, err error) {
	workspace, err := workspace.GetWorkspace(cmd.Context(), cmd)
	if err != nil {
		return []Reviewer{}, err
	}
	return profile.GetAll[Reviewer](context, cmd, fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers", workspace, project))
}

// disableUnsupportedFlags disables the flags that are not supported by the project reviewer command
func disableUnsupportedFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("repository") {
		return fmt.Errorf("the --repository flag is not supported by the project reviewer command")
	}
	return nil
}

// hideUnsupportedFlags hides the flags that are not supported by the repository command
func hideUnsupportedFlags(cmd *cobra.Command, args []string) {
	cmd.Flags().MarkHidden("repository")
	cmd.Parent().HelpFunc()(cmd, args)
}
