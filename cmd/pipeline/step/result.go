package step

type StepResult struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// String returns the name of the result.
//
// implements fmt.Stringer
func (result StepResult) String() string {
	return result.Name
}
