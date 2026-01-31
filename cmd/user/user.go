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
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
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

var UserCache = common.NewCache[User]()

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

var columns = common.Columns[User]{
	{Name: "id", DefaultSorter: true, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.ID.String()), strings.ToLower(b.ID.String())) == -1
	}},
	{Name: "username", DefaultSorter: false, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.Username), strings.ToLower(b.Username)) == -1
	}},
	{Name: "name", DefaultSorter: false, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "nickname", DefaultSorter: false, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.Nickname), strings.ToLower(b.Nickname)) == -1
	}},
	{Name: "account", DefaultSorter: false, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.AccountID), strings.ToLower(b.AccountID)) == -1
	}},
	{Name: "created_on", DefaultSorter: false, Compare: func(a, b User) bool {
		return a.CreatedOn.Before(b.CreatedOn)
	}},
	{Name: "account_status", DefaultSorter: false, Compare: func(a, b User) bool {
		return strings.Compare(strings.ToLower(a.AccountStatus), strings.ToLower(b.AccountStatus)) == -1
	}},
}

// GetID gets the ID of the user
//
// implements core.Identifiable
func (user User) GetID() uuid.UUID {
	return uuid.UUID(user.ID)
}

// GetName gets the name of the user
//
// implements core.Named
func (user User) GetName() string {
	return user.Username
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (user User) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return columns
		}
	}
	return []string{"ID", "Username", "Name"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (user User) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, user.ID.String())
		case "username":
			if len(user.Username) > 0 {
				row = append(row, user.Username)
			} else {
				row = append(row, user.Nickname)
			}
		case "name":
			row = append(row, user.Name)
		case "nickname":
			row = append(row, user.Nickname)
		case "account":
			row = append(row, user.AccountID)
		case "created on", "created_on":
			if user.CreatedOn.IsZero() {
				row = append(row, " ")
			} else {
				row = append(row, user.CreatedOn.Format("2006-01-02 15:04:05"))
			}
		case "account status":
			if user.AccountStatus == "" {
				row = append(row, " ")
			} else {
				row = append(row, user.AccountStatus)
			}
		}
	}
	return row
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
	log := logger.Must(logger.FromContext(context)).Child("user", "me")
	if user, err = UserCache.Get("me"); err == nil {
		log.Debugf("User found in cache")
		return
	}
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
	if err == nil {
		_ = UserCache.Set(*user, "me")
	}
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
	userUUID, err := common.ParseUUID(userid)
	if err == nil {
		if user, err = UserCache.Get(userUUID.String()); err != nil {
			err = profile.Get(
				context,
				cmd,
				fmt.Sprintf("/users/%s", userUUID.String()),
				&user,
			)
			if err == nil {
				_ = UserCache.Set(*user)
			}
		}
	}
	return
}

// GetUserFromFlags gets the user from the command
func GetUserFromFlags(context context.Context, cmd *cobra.Command) (*User, error) {
	if cmd.Flag("user") == nil {
		return nil, errors.Errorf("The command %s does not have a --user flag", cmd.Name())
	}
	return GetUser(context, cmd, cmd.Flag("user").Value.String())
}
