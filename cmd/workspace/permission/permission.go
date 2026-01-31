package permission

import (
	"encoding/json"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	wkcommon "bitbucket.org/gildas_cherruel/bb/cmd/workspace/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type Permission struct {
	User       user.User          `json:"user"             mapstructure:"user"`
	Permission string             `json:"permission"       mapstructure:"permission"`
	Workspace  wkcommon.Workspace `json:"workspace"        mapstructure:"workspace"`
	Links      common.Links       `json:"links"            mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "permission",
	Short: "Manage permissions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Permission requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Permission]{
	{Name: "user", DefaultSorter: false, Compare: func(a, b Permission) bool {
		return strings.Compare(strings.ToLower(a.User.Name), strings.ToLower(b.User.Name)) == -1
	}},
	{Name: "permission", DefaultSorter: false, Compare: func(a, b Permission) bool {
		return strings.Compare(strings.ToLower(a.Permission), strings.ToLower(b.Permission)) == -1
	}},
	{Name: "workspace", DefaultSorter: false, Compare: func(a, b Permission) bool {
		return strings.Compare(strings.ToLower(a.Workspace.Name), strings.ToLower(b.Workspace.Name)) == -1
	}},
}

// GetHeaders get the headers for a table
//
// implements common.Tableable
func (permission Permission) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Workspace", "User", "Permission"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (permission Permission) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "user":
			row = append(row, permission.User.Name)
		case "permission":
			row = append(row, permission.Permission)
		case "workspace":
			row = append(row, permission.Workspace.Slug)
		}
	}
	return row
}

// GetType gets the type of the permission
//
// implements core.TypeCarrier
func (permission Permission) GetType() string {
	return "workspace_membership"
}

// MarshalJSON marshals the permission to JSON
//
// implements json.Marshaler
func (permission Permission) MarshalJSON() ([]byte, error) {
	type surrogate Permission

	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      permission.GetType(),
		surrogate: surrogate(permission),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the permission from JSON
//
// implements json.Unmarshaler
func (permission *Permission) UnmarshalJSON(data []byte) error {
	type surrogate Permission

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != permission.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, permission.GetType()))
	}

	*permission = Permission(inner.surrogate)
	return nil
}
