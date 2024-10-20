package component

import (
	"context"
	"fmt"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
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

// GetHeader gets the header for a table
//
// implements common.Tableable
func (issue Component) GetHeader(short bool) []string {
	return []string{"ID", "Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (component Component) GetRow(headers []string) []string {
	return []string{
		fmt.Sprintf("%d", component.ID),
		component.Name,
	}
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
