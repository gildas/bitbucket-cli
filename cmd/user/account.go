package user

import (
	"encoding/json"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
)

type Account struct {
	Type          string       `json:"type"           mapstructure:"type"`
	ID            common.UUID  `json:"uuid"           mapstructure:"uuid"`
	Username      string       `json:"username"       mapstructure:"username"`
	Name          string       `json:"display_name"   mapstructure:"display_name"`
	AccountID     string       `json:"account_id"     mapstructure:"account_id"`
	AccountStatus string       `json:"account_status" mapstructure:"account_status"`
	Kind          string       `json:"kind"           mapstructure:"kind"`
	Links         common.Links `json:"links"          mapstructure:"links"`
	CreatedOn     time.Time    `json:"created_on"     mapstructure:"created_on"`
}

// MarshalJSON implements the json.Marshaler interface.
func (account Account) MarshalJSON() (data []byte, err error) {
	type surrogate Account

	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
	}{
		surrogate: surrogate(account),
		CreatedOn: account.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
