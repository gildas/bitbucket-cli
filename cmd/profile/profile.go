package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

// Profile describes the configuration needed to connect to BitBucket
type Profile struct {
	Name         string    `json:"name"                   mapstructure:"name"`
	Description  string    `json:"description,omitempty"  mapstructure:"description,omitempty" yaml:",omitempty"`
	Default      bool      `json:"default"                mapstructure:"default"               yaml:",omitempty"`
	User         string    `json:"user,omitempty"         mapstructure:"user"                  yaml:",omitempty"`
	Password     string    `json:"password,omitempty"     mapstructure:"password"              yaml:",omitempty"`
	ClientID     string    `json:"clientID,omitempty"     mapstructure:"clientID"              yaml:",omitempty"`
	ClientSecret string    `json:"clientSecret,omitempty" mapstructure:"clientSecret"          yaml:",omitempty"`
	AccessToken  string    `json:"accessToken,omitempty"  mapstructure:"accessToken"           yaml:",omitempty"`
	RefreshToken string    `json:"-"                      mapstructure:"refreshToken"          yaml:"-"`
	TokenExpires time.Time `json:"-"                      mapstructure:"tokenExpires"          yaml:"-"`
	TokenScopes  []string  `json:"-"                      mapstructure:"tokenScopes"           yaml:"-"`
}

// Current is the current profile
var Current *Profile

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
	// We must have either an access token or a user/password or a clientID/clientSecret
	if len(profile.AccessToken) == 0 && len(profile.ClientID) == 0 && len(profile.User) == 0 {
		merr.Append(errors.ArgumentMissing.With("accessToken, or user/password, or clientID/clientSecret"))
	} else if len(profile.AccessToken) == 0 {
		if len(profile.User) != 0 {
			if len(profile.Password) == 0 {
				merr.Append(errors.ArgumentMissing.With("password"))
			}
		} else if len(profile.ClientID) != 0 {
			if len(profile.ClientSecret) == 0 {
				merr.Append(errors.ArgumentMissing.With("clientSecret"))
			}
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

// loadAccessToken loads the access token from the cache
func (profile *Profile) loadAccessToken() (err error) {
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		accessTokenFile := filepath.Join(cacheDir, "bitbucket", "access-token-"+profile.Name)
		data, err := os.ReadFile(accessTokenFile)
		if err == nil {
			var token struct {
				TokenType    string         `json:"token_type"`
				AccessToken  string         `json:"access_token"`
				RefreshToken string         `json:"refresh_token"`
				ExpiresOn    core.Timestamp `json:"expires_on"`
				Scope        string         `json:"scope"`
			}
			if err = json.Unmarshal(data, &token); err == nil {
				profile.AccessToken = token.AccessToken
				profile.RefreshToken = token.RefreshToken
				profile.TokenExpires = time.Time(token.ExpiresOn)
				profile.TokenScopes = strings.Split(token.Scope, " ")
				return err
			}
		}
		return err
	}
	return
}

// isTokenExpired tells if the token is expired
func (profile *Profile) isTokenExpired() bool {
	return profile.TokenExpires.Before(time.Now())
}

// saveAccessToken saves the access token to the cache
func (profile *Profile) saveAccessToken(data []byte) {
	var payload []byte = data
	if err := profile.setFromBitbucketTokenData(data); err == nil {
		payload = profile.getTokenData()
	} else {
		profile.AccessToken = string(data)
	}
	if cacheDir, err := os.UserCacheDir(); err == nil {
		cachePath := filepath.Join(cacheDir, "bitbucket")
		if err := os.MkdirAll(cachePath, 0700); err == nil {
			cacheFile := filepath.Join(cachePath, "access-token-"+profile.Name)
			if err := os.WriteFile(cacheFile, payload, 0600); err == nil {
				return
			}
		}
	}
}

// setFromBitbucketTokenData sets the profile token information from the BitBucket token data
//
// The original data carries an expiration duration, that needs to be converted to a time.Time
func (profile *Profile) setFromBitbucketTokenData(data []byte) (err error) {
	var token struct {
		TokenType    string `json:"token_type"`
		State        string `json:"state"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		Scopes       string `json:"scopes"`
	}
	if err = json.Unmarshal(data, &token); err == nil {
		profile.AccessToken = token.AccessToken
		profile.RefreshToken = token.RefreshToken
		profile.TokenExpires = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
		profile.TokenScopes = strings.Split(token.Scopes, " ")
	}
	return
}

// getTokenData gets the token data from the profile
//
// This data carries an expiration date as a timestamp
func (profile *Profile) getTokenData() (data []byte) {
	token := struct {
		TokenType    string         `json:"token_type"`
		AccessToken  string         `json:"access_token"`
		RefreshToken string         `json:"refresh_token"`
		ExpiresOn    core.Timestamp `json:"expires_on"`
		Scopes       string         `json:"scopes"`
	}{
		TokenType:    "bearer",
		AccessToken:  profile.AccessToken,
		RefreshToken: profile.RefreshToken,
		ExpiresOn:    core.Timestamp(profile.TokenExpires),
		Scopes:       strings.Join(profile.TokenScopes, " "),
	}
	data, _ = json.Marshal(token)
	return
}
