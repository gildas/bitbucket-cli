package common

type RenderedText struct {
	Type   string `json:"type,omitempty"   mapstructure:"type"`
	Raw    string `json:"raw"              mapstructure:"raw"`
	Markup string `json:"markup,omitempty" mapstructure:"markup"` // markdown, creaole, plaintext
	HTML   string `json:"html,omitempty"   mapstructure:"html"`
}
