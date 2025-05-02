package artifact

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Artifact struct {
	Name      string       `json:"name"      mapstructure:"name"`
	Size      uint64       `json:"size"      mapstructure:"size"`
	Downloads uint64       `json:"downloads" mapstructure:"downloads"`
	User      user.User    `json:"user"      mapstructure:"user"`
	Links     common.Links `json:"links"     mapstructure:"links"`
}

var Command = &cobra.Command{
	Use:   "artifact",
	Short: "Manage artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Artifact requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (artifact Artifact) GetHeader(short bool) []string {
	return []string{"Name", "Size", "Downloads", "Owner"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (artifact Artifact) GetRow(headers []string) []string {
	return []string{
		artifact.Name,
		fmt.Sprintf("%d", artifact.Size),
		fmt.Sprintf("%d", artifact.Downloads),
		artifact.User.Name,
	}
}

// GetArtifactNames gets the names of the artifacts
func GetArtifactNames(context context.Context, cmd *cobra.Command) (names []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("artifact", "getnames")

	artifacts, err := profile.GetAll[Artifact](cmd.Context(), cmd, "downloads")
	if err != nil {
		log.Errorf("Failed to get artifacts: %s", err)
		return
	}
	names = core.Map(artifacts, func(artifact Artifact) string { return artifact.Name })
	core.Sort(names, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return names, nil
}
