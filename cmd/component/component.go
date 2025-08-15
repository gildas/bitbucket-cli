package component

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type Component struct {
	Type  string       `json:"type"  mapstructure:"type"`
	ID    int          `json:"id"    mapstructure:"id"`
	Name  string       `json:"name"  mapstructure:"name"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "component",
	Short: "Manage components",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Issue requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Component]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b Component) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "type", DefaultSorter: false, Compare: func(a, b Component) bool {
		return strings.Compare(strings.ToLower(a.Type), strings.ToLower(b.Type)) == -1
	}},
	{Name: "id", DefaultSorter: false, Compare: func(a, b Component) bool {
		return a.ID < b.ID
	}},
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (issue Component) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (component Component) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, fmt.Sprintf("%d", component.ID))
		case "name":
			row = append(row, component.Name)
		case "type":
			row = append(row, component.Type)
		}
	}
	return row
}

// String gets a string representation
//
// implements fmt.Stringer
func (component Component) String() string {
	return component.Name
}

// GetComponentIDs gets the IDs of the components
func GetComponentIDs(context context.Context, cmd *cobra.Command) (ids []string) {
	log := logger.Must(logger.FromContext(context)).Child("component", "getids")

	components, err := profile.GetAll[Component](context, cmd, "components")
	if err != nil {
		log.Errorf("Failed to get components", err)
		return []string{}
	}
	ids = make([]string, 0, len(components))
	for _, component := range components {
		ids = append(ids, fmt.Sprintf("%d", component.ID))
	}
	return
}
