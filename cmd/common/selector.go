package common

// Selector represents a selector for references (pipelines, branches, tags, etc.)
type Selector struct {
	Type    string `json:"type"              mapstructure:"type"`
	Pattern string `json:"pattern,omitempty" mapstructure:"pattern"`
}
