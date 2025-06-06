package sshkey

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

type SSHKey struct {
	Type        string       `json:"type"        mapstructure:"type"`
	ID          common.UUID  `json:"uuid"        mapstructure:"uuid"`
	Label       string       `json:"label"       mapstructure:"label"`
	Comment     string       `json:"comment"     mapstructure:"comment"`
	Key         string       `json:"key"         mapstructure:"key"`
	Fingerprint string       `json:"fingerprint" mapstructure:"fingerprint"`
	CreatedOn   time.Time    `json:"created_on"  mapstructure:"created_on"`
	ExpiresOn   time.Time    `json:"expires_on"  mapstructure:"expires_on"`
	LastUsed    time.Time    `json:"last_used"   mapstructure:"last_used"`
	Owner       user.User    `json:"owner"       mapstructure:"owner"`
	Links       common.Links `json:"links"       mapstructure:"links"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "ssh-key",
	Short: "Manage SSH keys",
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
func (key SSHKey) GetHeader(short bool) []string {
	return []string{"Key ID", "Name", "Owner", "Fingerprint", "Comment"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (key SSHKey) GetRow(headers []string) []string {
	return []string{key.ID.String(), key.Label, key.Owner.Name, key.Fingerprint, key.Comment}
}

// GetSSHKeys gets the SSHKeys
func GetSSHKeys(context context.Context, cmd *cobra.Command) (keys []SSHKey, err error) {
	owner, err := user.GetUserFromFlags(context, cmd)
	if err != nil {
		return
	}
	return profile.GetAll[SSHKey](
		cmd.Context(),
		cmd,
		fmt.Sprintf("/users/%s/ssh-keys", owner.ID.String()),
	)
}

// GetSSHKeyFingerprints gets the fingerprints of the SSHKeys
func GetSSHKeyFingerprints(context context.Context, cmd *cobra.Command) (fingerprints []string, err error) {
	keys, err := GetSSHKeys(context, cmd)
	if err != nil {
		return []string{}, err
	}
	fingerprints = core.Map(keys, func(key SSHKey) string { return key.Fingerprint })
	core.Sort(fingerprints, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return fingerprints, nil
}

// String gets a string representation of the SSHKey
//
// implements fmt.Stringer
func (key SSHKey) String() string {
	return key.Fingerprint
}

// Validate validates the SSHKey
func (key SSHKey) Validate() error {
	var merr errors.MultiError

	if len(key.Fingerprint) == 0 {
		merr.Append(errors.ArgumentMissing.With("fingerprint"))
	}
	if len(key.Key) == 0 {
		merr.Append(errors.ArgumentMissing.With("key"))
	}
	return merr.AsError()
}

// MarshalJSON marshals the SSHKey
//
// implements json.Marshaler
func (key SSHKey) MarshalJSON() ([]byte, error) {
	type surrogate SSHKey

	data, err := json.Marshal(struct {
		surrogate
		CreatedOn string `json:"created_on"`
		ExpiresOn string `json:"expires_on"`
		LastUsed  string `json:"last_used"`
	}{
		surrogate: surrogate(key),
		CreatedOn: key.CreatedOn.Format(time.RFC3339),
		ExpiresOn: key.ExpiresOn.Format(time.RFC3339),
		LastUsed:  key.LastUsed.Format(time.RFC3339),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}
