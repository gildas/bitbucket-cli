package common

type Links struct {
	Self           *Link `json:"self,omitempty"            mapstructure:"self"`
	HTML           *Link `json:"html,omitempty"            mapstructure:"html"`
	Avatar         *Link `json:"avatar,omitempty"          mapstructure:"avatar"`
	Commits        *Link `json:"commits,omitempty"         mapstructure:"commits"`
	Approve        *Link `json:"approve,omitempty"         mapstructure:"approve"`
	RequestChanges *Link `json:"request-changes,omitempty" mapstructure:"request-changes"`
	Diff           *Link `json:"diff,omitempty"            mapstructure:"diff"`
	DiffStat       *Link `json:"diffstat,omitempty"        mapstructure:"diffstat"`
	Comments       *Link `json:"comments,omitempty"        mapstructure:"comments"`
	Activity       *Link `json:"activity,omitempty"        mapstructure:"activity"`
	Merge          *Link `json:"merge,omitempty"           mapstructure:"merge"`
	Decline        *Link `json:"decline,omitempty"         mapstructure:"decline"`
	Statuses       *Link `json:"statuses,omitempty"        mapstructure:"statuses"`
}
