package pipeline

import "fmt"

// PipelineState represents the state of a pipeline
type PipelineState struct {
	Type   string          `json:"type"             mapstructure:"type"`
	Name   string          `json:"name"             mapstructure:"name"`
	Stage  *PipelineStage  `json:"stage,omitempty"  mapstructure:"stage"`
	Result *PipelineResult `json:"result,omitempty" mapstructure:"result"`
}

// String returns the name of the state.
//
// implements fmt.Stringer
func (state PipelineState) String() string {
	if state.Result != nil {
		if state.Stage != nil {
			return fmt.Sprintf("%s - %s (%s)", state.Name, state.Stage, state.Result)
		}
		return fmt.Sprintf("%s (%s)", state.Name, state.Result)
	}
	if state.Stage != nil {
		return fmt.Sprintf("%s - %s", state.Name, state.Stage)
	}
	return state.Name
}
