package common

type RenderedText struct {
	Type   string `json:"type" mapstructure:"type"`
	Raw    string `json:"raw"  mapstructure:"raw"`
	Markup string `json:"markup"  mapstructure:"markup"`
	HTML   string `json:"html"  mapstructure:"html"`
}
