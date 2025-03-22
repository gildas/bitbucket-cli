package gpgkey

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type GPGKey struct {
	Type        string       `json:"type" mapstructure:"type"`
	Parent      string       `json:"parent_fingerprint" mapstructure:"parent_fingerprint"`
	Fingerprint string       `json:"fingerprint" mapstructure:"fingerprint"`
	KeyID       string       `json:"key_id" mapstructure:"key_id"`
	Name        string       `json:"name" mapstructure:"name"`
	AddedOn     time.Time    `json:"added_on" mapstructure:"added_on"`
	CreatedOn   time.Time    `json:"created_on" mapstructure:"created_on"`
	Links       common.Links `json:"links" mapstructure:"links"`
	Owner       user.User    `json:"owner" mapstructure:"owner"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "gpg-key",
	Aliases: []string{"key"}, // backward compatibility
	Short:   "Manage GPG keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Key requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// GetHeader gets the header for a table
//
// implements common.Tableable
func (key GPGKey) GetHeader(short bool) []string {
	return []string{"Fingerprint", "Name", "Owner"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (key GPGKey) GetRow(headers []string) []string {
	return []string{key.Fingerprint, key.Name, key.Owner.Name}
}

// GetGPGKeys gets the GPGKeys
func GetGPGKeys(context context.Context, cmd *cobra.Command) (keys []GPGKey, err error) {
	owner, err := user.GetUserFromFlags(context, cmd)
	if err != nil {
		return
	}
	return profile.GetAll[GPGKey](
		cmd.Context(),
		cmd,
		fmt.Sprintf("/users/%s/gpg-keys", owner.ID.String()),
	)
}

// GetGPGKeyFingerprints gets the fingerprints of the GPGKeys
func GetGPGKeyFingerprints(context context.Context, cmd *cobra.Command) []string {
	keys, err := GetGPGKeys(context, cmd)
	if err != nil {
		return []string{}
	}
	return core.Map(keys, func(key GPGKey) string {
		return key.Fingerprint
	})
}

// String gets a string representation of the GPGKey
//
// implements fmt.Stringer
func (key GPGKey) String() string {
	return key.Fingerprint
}

// Validate validates the GPGKey
func (key GPGKey) Validate() error {
	var merr errors.MultiError

	if len(key.Fingerprint) == 0 {
		merr.Append(errors.ArgumentMissing.With("fingerprint"))
	}
	if len(key.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("name"))
	}
	return merr.AsError()
}

// MarshalJSON marshals the GPGKey
//
// implements json.Marshaler
func (key GPGKey) MarshalJSON() ([]byte, error) {
	type surrogate GPGKey

	data, err := json.Marshal(struct {
		surrogate
		AddedOn   string `json:"added_on"`
		CreatedOn string `json:"created_on"`
	}{
		surrogate: surrogate(key),
		AddedOn:   key.AddedOn.Format(time.RFC3339),
		CreatedOn: key.CreatedOn.Format(time.RFC3339),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
