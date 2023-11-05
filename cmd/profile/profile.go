package profile

import (
	"fmt"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// Profile describes the configuration needed to connect to BitBucket
type Profile struct {
	Name        string `json:"name"                  mapstructure:"name"`
	Description string `json:"description,omitempty" mapstructure:"description,omitempty" yaml:",omitempty"`
	Default     bool   `json:"default"               mapstructure:"default"               yaml:",omitempty"`
	User        string `json:"user,omitempty"        mapstructure:"user"                  yaml:",omitempty"`
	Password    string `json:"-"                     mapstructure:"password"              yaml:",omitempty"`
	AccessToken string `json:"accessToken,omitempty" mapstructure:"accessToken"           yaml:",omitempty"`
}

// Log is the logger for this application
var Log *logger.Logger

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Profile requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// Validate validates a Profile
func (profile *Profile) Validate() error {
	var merr errors.MultiError

	if len(profile.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("name"))
	}
	if _, found := Profiles.Find(profile.Name); found {
		merr.Append(errors.DuplicateFound.With("name", profile.Name))
	}
	if len(profile.AccessToken) == 0 {
		if len(profile.User) == 0 {
			merr.Append(errors.ArgumentMissing.With("user"))
		}
		if len(profile.Password) == 0 {
			merr.Append(errors.ArgumentMissing.With("password"))
		}
	}
	return merr.AsError()
}

// String gets a string representation of this profile
//
// implements fmt.Stringer
func (profile Profile) String() string {
	return profile.Name
}
