package user

type Author struct {
	Type string `json:"type"          mapstructure:"type"`
	Raw  string `json:"raw,omitempty" mapstructure:"raw"`
	User User   `json:"user"          mapstructure:"user"`
}
