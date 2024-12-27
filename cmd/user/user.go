package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type User struct {
	Type          string       `json:"type"                     mapstructure:"type"`
	ID            common.UUID  `json:"uuid"                     mapstructure:"uuid"`
	AccountID     string       `json:"account_id"               mapstructure:"account_id"`
	Username      string       `json:"username,omitempty"       mapstructure:"username"`
	Name          string       `json:"display_name"             mapstructure:"display_name"`
	Nickname      string       `json:"nickname"                 mapstructure:"nickname"`
	Raw           string       `json:"raw,omitempty"            mapstructure:"raw"`
	Kind          string       `json:"kind,omitempty"           mapstructure:"kind"`
	Links         common.Links `json:"links"                    mapstructure:"links"`
	CreatedOn     time.Time    `json:"created_on"               mapstructure:"created_on"`
	AccountStatus string       `json:"account_status,omitempty" mapstructure:"account_status"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "user",
	Aliases: []string{"account"},
	Short:   "Manage users",
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
func (user User) GetHeader(short bool) []string {
	return []string{"ID", "Username", "Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (user User) GetRow(headers []string) []string {
	return []string{
		user.ID.String(),
		user.Username,
		user.Name,
	}
}

// String gets the string representation of the user
//
// implements fmt.Stringer
func (user User) String() string {
	if len(user.Name) == 0 {
		return user.ID.String()
	}
	return user.Name
}

// MarshalJSON implements the json.Marshaler interface.
func (user User) MarshalJSON() (data []byte, err error) {
	type surrogate User
	var createdOn string

	if !user.CreatedOn.IsZero() {
		createdOn = user.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}
	data, err = json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on,omitempty"`
	}{
		surrogate: surrogate(user),
		CreatedOn: createdOn,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// GetMe gets the current user
func GetMe(context context.Context, cmd *cobra.Command) (user *User, err error) {
	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}
	err = profile.Get(
		context,
		cmd,
		"/user",
		&user,
	)
	return
}

// GetUser gets a user
func GetUser(context context.Context, cmd *cobra.Command, userid string) (user *User, err error) {
	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return nil, err
	}
	if len(userid) == 0 || strings.ToLower(userid) == "me" || strings.ToLower(userid) == "myself" {
		me, err := GetMe(context, cmd)
		if err != nil {
			return nil, err
		}
		return me, nil
	}
	uuid, err := common.ParseUUID(userid)
	if err == nil {
		err = profile.Get(
			context,
			cmd,
			fmt.Sprintf("/users/%s", uuid.String()),
			&user,
		)
	}
	return
}

// GetUserFromFlags gets the user from the command
func GetUserFromFlags(context context.Context, cmd *cobra.Command) (*User, error) {
	id := cmd.Flag("user").Value.String()
	if len(id) == 0 {
		id = cmd.Flag("account").Value.String()
	}
	return GetUser(context, cmd, id)
}
