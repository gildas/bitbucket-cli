package plcommon

type PipelineReference struct {
	Type string `json:"type" mapstructure:"type"`
	UUID string `json:"uuid" mapstructure:"uuid"`
}
