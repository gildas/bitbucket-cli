package pipeline

import (
	"encoding/json"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

// Variable represents a pipeline variable
type Variable struct {
	ID      common.UUID `json:"uuid"    mapstructure:"uuid,omitempty"`
	Key     string      `json:"key"     mapstructure:"key"`
	Value   string      `json:"value"   mapstructure:"value,omitempty"`
	Secured bool        `json:"secured" mapstructure:"secured,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for Variable
//
// Implements json.Marshaler
func (variable Variable) MarshalJSON() ([]byte, error) {
	type surrogate Variable
	var id *common.UUID
	if !variable.ID.IsNil() {
		id = &variable.ID
	}

	data, err := json.Marshal(struct {
		ID *common.UUID `json:"uuid,omitempty"`
		surrogate
	}{
		ID:        id,
		surrogate: surrogate(variable),
	})
	return data, err
}
