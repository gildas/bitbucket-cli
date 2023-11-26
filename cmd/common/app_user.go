package common

import (
	"encoding/json"
	"time"

	"github.com/gildas/go-errors"
)

type AppUser struct {
	Type          string    `json:"type"           mapstructure:"type"`
	ID            string    `json:"uuid"           mapstructure:"uuid"`
	Name          string    `json:"display_name"   mapstructure:"display_name"`
	AccountID     string    `json:"account_id"     mapstructure:"account_id"`
	AccountStatus string    `json:"account_status" mapstructure:"account_status"`
	Kind          string    `json:"kind"           mapstructure:"kind"`
	Links         Links     `json:"links"          mapstructure:"links"`
	CreatedOn     time.Time `json:"created_on"     mapstructure:"created_on"`
}

// MarshalJSON implements the json.Marshaler interface.
func (user AppUser) MarshalJSON() (data []byte, err error) {
	type surrogate AppUser

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
	}{
		surrogate: surrogate(user),
		CreatedOn: user.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
