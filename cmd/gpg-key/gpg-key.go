package gpgkey

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

type GPGKey struct {
	Type        string       `json:"type"               mapstructure:"type"`
	Parent      string       `json:"parent_fingerprint" mapstructure:"parent_fingerprint"`
	Fingerprint string       `json:"fingerprint"        mapstructure:"fingerprint"`
	KeyID       string       `json:"key_id"             mapstructure:"key_id"`
	Name        string       `json:"name"               mapstructure:"name"`
	AddedOn     time.Time    `json:"added_on"           mapstructure:"added_on"`
	CreatedOn   time.Time    `json:"created_on"         mapstructure:"created_on"`
	Links       common.Links `json:"links"              mapstructure:"links"`
	Owner       user.User    `json:"owner"              mapstructure:"owner"`
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

var columns = []string{
	"fingerprint",
	"name",
	"owner",
	"added_on",
	"created_on",
	"type",
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (key GPGKey) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Fingerprint", "Name", "Owner"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (key GPGKey) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "added_on", "added on":
			row = append(row, key.AddedOn.Format("2006-01-02 15:04:05"))
		case "created_on", "created on":
			row = append(row, key.CreatedOn.Format("2006-01-02 15:04:05"))
		case "key_id":
			row = append(row, key.KeyID)
		case "parent", "parent_fingerprint":
			row = append(row, key.Parent)
		case "fingerprint":
			row = append(row, key.Fingerprint)
		case "name":
			row = append(row, key.Name)
		case "owner":
			if key.Owner.Name == "" {
				row = append(row, " ")
			} else {
				row = append(row, key.Owner.Name)
			}
		case "type":
			row = append(row, key.Type)
		}
	}
	return row
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
func GetGPGKeyFingerprints(context context.Context, cmd *cobra.Command) (fingerprints []string, err error) {
	keys, err := GetGPGKeys(context, cmd)
	if err != nil {
		return
	}
	fingerprints = core.Map(keys, func(key GPGKey) string { return key.Fingerprint })
	core.Sort(fingerprints, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return fingerprints, nil
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
