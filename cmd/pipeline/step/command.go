package step

type StepCommand struct {
	Type    string `json:"commandType" mapstructure:"commandType"`
	Name    string `json:"name"        mapstructure:"name"`
	Command string `json:"command"     mapstructure:"command"`
}
