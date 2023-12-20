package common

type Entity struct {
	Type  string `json:"type"  mapstructure:"type"`
	ID    int    `json:"id"    mapstructure:"id"`
	Name  string `json:"name"  mapstructure:"name"`
	Links Links  `json:"links" mapstructure:"links"`
}
