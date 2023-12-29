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

// GetProjectKeys gets the keys of the projects in the workspace given in the command
func GetProjectKeys(context context.Context, cmd *cobra.Command, args []string) (keys []string) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")

	workspace := cmd.Flag("workspace").Value.String()
	if len(workspace) == 0 {
		workspace = profile.Current.DefaultWorkspace
		if len(workspace) == 0 {
			log.Warnf("No workspace given")
			return
		}
	}

	type Project struct {
		Key string `json:"key" mapstructure:"key"`
	}

	log.Infof("Getting all projects from workspace %s", workspace)
	projects, err := profile.GetAll[Project](context, cmd, profile.Current, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	return core.Map(projects, func(project Project) string {
		return project.Key
	})
}

// GetReviewerIDs gets the IDs of the reviewers in the given workspace and project
func GetReviewerUserIDs(context context.Context, cmd *cobra.Command, currentProfile *profile.Profile, workspace, project string) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("reviewer", "getids")

	reviewers, err := profile.GetAll[Reviewer](context, cmd, currentProfile, fmt.Sprintf("/workspaces/%s/projects/%s/default-reviewers", workspace, project))
	if err != nil {
		log.Errorf("Failed to get reviewers", err)
		return
	}
	return core.Map(reviewers, func(reviewer Reviewer) string {
		return reviewer.User.ID.String()
	})
}
