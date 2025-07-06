package reviewer

import (
	"context"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
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

// GetHeader gets the header for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetHeader(short bool) []string {
	return []string{"Type", "Reviewer Type", "User"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (reviewer Reviewer) GetRow(headers []string) []string {
	return []string{reviewer.Type, reviewer.ReviewerType, reviewer.User.Name}
}

// Validate validates a Reviewer
func (reviewer *Reviewer) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// GetWorkspaceAndProject gets the workspace and project from the command
func GetWorkspaceAndProject(cmd *cobra.Command, profile *profile.Profile) (workspace, project string, err error) {
	workspace = cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = profile.DefaultWorkspace
		if len(workspace) == 0 {
			return "", "", errors.ArgumentMissing.With("workspace")
		}
	}

	project = cmd.Flag("project").Value.String()
	if len(project) == 0 {
		project = profile.DefaultProject
		if len(project) == 0 {
			return "", "", errors.ArgumentMissing.With("project")
		}
	}
	return
}

// GetProjectKeys gets the keys of the projects in the workspace given in the command
func GetProjectKeys(context context.Context, cmd *cobra.Command, args []string, toComplete string) (keys []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")

	currentProfile, err := profile.GetProfileFromCommand(context, cmd)
	if err != nil {
		log.Errorf("Failed to get profile.", err)
		return nil, err
	}

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = currentProfile.DefaultWorkspace
		if len(workspace) == 0 {
			log.Warnf("No workspace given")
			return
		}
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

	currentProfile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return []string{}, err
	}

	workspace := deleteOptions.Workspace.Value
	if len(workspace) == 0 {
		workspace = currentProfile.DefaultWorkspace
		if len(workspace) == 0 {
			return []string{}, errors.ArgumentMissing.With("workspace")
		}
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

// GetDefaultReviewers gets the reviewers in the given workspace and project
func GetDefaultReviewers(context context.Context, cmd *cobra.Command, workspace, project string) (reviewers []Reviewer, err error) {
	return profile.GetAll[Reviewer](context, cmd, fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers", workspace, project))
}
