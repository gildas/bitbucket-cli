package pipeline

// PipelineStage represents the current stage of a pipeline
type PipelineStage struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// String returns the name of the stage.
//
// implements fmt.Stringer
func (stage PipelineStage) String() string {
	return stage.Name
}
