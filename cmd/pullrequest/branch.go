package pullrequest

type Branch struct {
	Name                 string   `json:"name"                             mapstructure:"name"`
	MergeStrategies      []string `json:"merge_strategies,omitempty"       mapstructure:"merge_strategies"`
	DefaultMergeStrategy string   `json:"default_merge_strategy,omitempty" mapstructure:"default_merge_strategy"`
}
