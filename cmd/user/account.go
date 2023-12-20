package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
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

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "account",
	Aliases: []string{"user"},
	Short:   "Manage accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Issue requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (account Account) GetHeader(short bool) []string {
	return []string{"ID", "Username", "Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (account Account) GetRow(headers []string) []string {
	return []string{
		account.ID.String(),
		account.Username,
		account.Name,
	}
}

// GetMe gets the current user
func GetMe(context context.Context, cmd *cobra.Command, currentProfile *profile.Profile) (account *Account, err error) {
	if currentProfile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	err = profile.Current.Get(
		context,
		cmd,
		"/user",
		&account,
	)
	return
}

// GetAccount gets a user
func GetAccount(context context.Context, cmd *cobra.Command, currentProfile *profile.Profile, userid string) (account *Account, err error) {
	if currentProfile == nil {
		return nil, errors.ArgumentMissing.With("profile")
	}

	uuid, err := common.ParseUUID(userid)
	if err == nil {
		err = profile.Current.Get(
			context,
			cmd,
			fmt.Sprintf("/users/%s", uuid.String()),
			&account,
		)
	}
	return
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
