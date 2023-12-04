package user

import "time"

type Participant struct {
	Type           string    `json:"type"            mapstructure:"type"`
	User           User      `json:"user"            mapstructure:"user"`
	Role           string    `json:"role"            mapstructure:"role"`
	Approved       bool      `json:"approved"        mapstructure:"approved"`
	State          string    `json:"state"           mapstructure:"state"`
	ParticipatedOn time.Time `json:"participated_on" mapstructure:"participated_on"`
}
