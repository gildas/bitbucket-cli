package workspace

import (
	"encoding/json"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type WorkspaceBase struct {
	ID    common.UUID  `json:"uuid"  mapstructure:"uuid"`
	Slug  string       `json:"slug"  mapstructure:"slug"`
	Links common.Links `json:"links" mapstructure:"links"`
}

// GetType gets the type of the workspace
//
// implements core.TypeCarrier
func (workspace WorkspaceBase) GetType() string {
	return "workspace_base"
}

// GetID gets the ID of the workspace
//
// implements core.Identifiable
func (workspace WorkspaceBase) GetID() uuid.UUID {
	return uuid.UUID(workspace.ID)
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (workspace WorkspaceBase) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"ID", "Slug"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (workspace WorkspaceBase) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, workspace.ID.String())
		case "slug":
			row = append(row, workspace.Slug)
		}
	}
	return row
}

// String returns the string representation of the workspace
//
// implements fmt.Stringer
func (workspace WorkspaceBase) String() string {
	return workspace.Slug
}

// MarshalJSON marshals the workspace to JSON
//
// implements json.Marshaler
func (workspace WorkspaceBase) MarshalJSON() ([]byte, error) {
	type surrogate WorkspaceBase
	data, err := json.Marshal(struct {
		Type string `json:"type"`
		surrogate
	}{
		Type:      workspace.GetType(),
		surrogate: surrogate(workspace),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the workspace from JSON
//
// implements json.Unmarshaler
func (workspace *WorkspaceBase) UnmarshalJSON(data []byte) error {
	type surrogate WorkspaceBase

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != workspace.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, workspace.GetType()))
	}

	*workspace = WorkspaceBase(inner.surrogate)
	return nil
}
