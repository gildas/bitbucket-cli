package user

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Participant struct {
	Type           string    `json:"type"            mapstructure:"type"`
	User           User      `json:"user"            mapstructure:"user"`
	Role           string    `json:"role"            mapstructure:"role"`
	Approved       bool      `json:"approved"        mapstructure:"approved"`
	State          string    `json:"state"           mapstructure:"state"`
	ParticipatedOn time.Time `json:"participated_on" mapstructure:"participated_on"`
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (participant Participant) GetHeaders(cmd *cobra.Command) []string {
	return []string{"ID", "Name", "participated on", "approved", "state"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (participant Participant) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, participant.User.ID.String())
		case "name":
			row = append(row, participant.User.Name)
		case "participated on":
			row = append(row, participant.ParticipatedOn.Local().String())
		case "approved":
			row = append(row, strconv.FormatBool(participant.Approved))
		case "state":
			row = append(row, participant.State)
		}
	}
	return row
}
