package user

import (
	"time"
)

type Participant struct {
	Type           string    `json:"type"            mapstructure:"type"`
	User           User      `json:"user"            mapstructure:"user"`
	Role           string    `json:"role"            mapstructure:"role"`
	Approved       bool      `json:"approved"        mapstructure:"approved"`
	State          string    `json:"state"           mapstructure:"state"`
	ParticipatedOn time.Time `json:"participated_on" mapstructure:"participated_on"`
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (participant Participant) GetHeader(short bool) []string {
	return []string{"ID", "Name", "participated on", "approved", "state"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (participant Participant) GetRow(headers []string) []string {
	return []string{
		participant.User.ID.String(),
		participant.User.Name,
		participant.ParticipatedOn.Local().String(),
		participant.State,
	}
}
