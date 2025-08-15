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

var columns = common.Columns[Artifact]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b Artifact) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "size", DefaultSorter: false, Compare: func(a, b Artifact) bool {
		return a.Size < b.Size
	}},
	{Name: "downloads", DefaultSorter: false, Compare: func(a, b Artifact) bool {
		return a.Downloads < b.Downloads
	}},
	{Name: "owner", DefaultSorter: false, Compare: func(a, b Artifact) bool {
		return strings.Compare(strings.ToLower(a.User.Username), strings.ToLower(b.User.Username)) == -1
	}},
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (artifact Artifact) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Name", "Size", "Downloads", "Owner"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (artifact Artifact) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "name":
			row = append(row, artifact.Name)
		case "size":
			row = append(row, fmt.Sprintf("%d", artifact.Size))
		case "downloads":
			row = append(row, fmt.Sprintf("%d", artifact.Downloads))
		case "owner":
			row = append(row, artifact.User.Name)
		}
	}
	return row
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
