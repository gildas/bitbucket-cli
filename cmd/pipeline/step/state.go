package step

import "fmt"

type StepState struct {
	Type   string      `json:"type" mapstructure:"type"`
	Name   string      `json:"name" mapstructure:"name"`
	Result *StepResult `json:"result,omitempty" mapstructure:"result"`
}

// String returns the name of the state.
//
// implements fmt.Stringer
func (state StepState) String() string {
	if state.Result != nil {
		return fmt.Sprintf("%s (%s)", state.Name, state.Result)
	}
	return state.Name
}
