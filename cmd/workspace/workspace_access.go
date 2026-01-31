package workspace

import (
	"encoding/json"

	"github.com/gildas/go-errors"
)

// WorkspaceAccess represents a workspace access entry
type WorkspaceAccess struct {
	Workspace     WorkspaceBase `json:"workspace"     mapstructure:"workspace"`
	Administrator bool          `json:"administrator" mapstructure:"administrator"`
}

// GetType gets the type of the workspace access
//
// implements core.TypeCarrier
func (access WorkspaceAccess) GetType() string {
	return "workspace_access"
}

// UnmarshalJSON implements the json.Unmarshaler interface
//
// implements json.Unmarshaler
func (access *WorkspaceAccess) UnmarshalJSON(data []byte) error {
	type surrogate WorkspaceAccess

	var inner struct {
		Type string `json:"type"`
		surrogate
	}

	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != access.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, access.GetType()))
	}

	*access = WorkspaceAccess(inner.surrogate)
	return nil
}
