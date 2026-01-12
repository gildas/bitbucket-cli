package common

type Links struct {
	Self           *Link  `json:"self,omitempty"            mapstructure:"self"`
	HTML           *Link  `json:"html,omitempty"            mapstructure:"html"`
	Avatar         *Link  `json:"avatar,omitempty"          mapstructure:"avatar"`
	Branches       *Link  `json:"branches,omitempty"        mapstructure:"branches"`
	Forks          *Link  `json:"forks,omitempty"           mapstructure:"forks"`
	Commits        *Link  `json:"commits,omitempty"         mapstructure:"commits"`
	PullRequests   *Link  `json:"pullrequests,omitempty"    mapstructure:"pullrequests"`
	Approve        *Link  `json:"approve,omitempty"         mapstructure:"approve"`
	RequestChanges *Link  `json:"request-changes,omitempty" mapstructure:"request-changes"`
	Diff           *Link  `json:"diff,omitempty"            mapstructure:"diff"`
	DiffStat       *Link  `json:"diffstat,omitempty"        mapstructure:"diffstat"`
	Patch          *Link  `json:"patch,omitempty"           mapstructure:"patch"`
	Comments       *Link  `json:"comments,omitempty"        mapstructure:"comments"`
	Activity       *Link  `json:"activity,omitempty"        mapstructure:"activity"`
	Merge          *Link  `json:"merge,omitempty"           mapstructure:"merge"`
	Decline        *Link  `json:"decline,omitempty"         mapstructure:"decline"`
	Statuses       *Link  `json:"statuses,omitempty"        mapstructure:"statuses"`
	Tags           *Link  `json:"tags,omitempty"            mapstructure:"tags"`
	Watchers       *Link  `json:"watchers,omitempty"        mapstructure:"watchers"`
	Downloads      *Link  `json:"downloads,omitempty"       mapstructure:"downloads"`
	Source         *Link  `json:"source,omitempty"          mapstructure:"source"`
	Clone          []Link `json:"clone,omitempty"           mapstructure:"clone"`
	Hooks          *Link  `json:"hooks,omitempty"           mapstructure:"hooks"`
	Steps          *Link  `json:"steps,omitempty"           mapstructure:"steps"`
}

// IsEmpty tells if there is no link defined
func (links Links) IsEmpty() bool {
	return links.Self == nil &&
		links.HTML == nil &&
		links.Avatar == nil &&
		links.Branches == nil &&
		links.Forks == nil &&
		links.Commits == nil &&
		links.PullRequests == nil &&
		links.Approve == nil &&
		links.RequestChanges == nil &&
		links.Diff == nil &&
		links.DiffStat == nil &&
		links.Patch == nil &&
		links.Comments == nil &&
		links.Activity == nil &&
		links.Merge == nil &&
		links.Decline == nil &&
		links.Statuses == nil &&
		links.Tags == nil &&
		links.Watchers == nil &&
		links.Downloads == nil &&
		links.Source == nil &&
		len(links.Clone) == 0 &&
		links.Hooks == nil &&
		links.Steps == nil
}
