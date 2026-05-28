package profile

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
)

type Token struct {
	TokenType    string         `json:"token_type"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresOn    core.Timestamp `json:"expires_on"`
	Scope        string         `json:"scope"`
}

// loadAccessToken loads the access token from the cache
func (profile *Profile) loadAccessToken(ctx context.Context) (err error) {

	log := logger.Must(logger.FromContext(ctx)).Child("profile", "loadAccessToken")

	if profile.token != nil {
		log.Debugf("Access token already loaded in memory for profile %s", profile.Name)
		return nil
	}

	if len(profile.AccessToken) > 0 {
		log.Debugf("Repository/Project/Workspace Access token for profile %s", profile.Name)
		profile.token = &Token{
			AccessToken: profile.AccessToken,
			ExpiresOn:   core.Timestamp(time.Now().Add(100 * 365 * 24 * time.Hour)), // Loaded Access Tokens never expire
		}
		return nil
	}

	// then load the access token from the file cache
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		accessTokenFile := filepath.Join(cacheDir, "bitbucket", "access-token-"+profile.Name)
		data, err := os.ReadFile(accessTokenFile)
		if err == nil {
			var token Token
			if err = json.Unmarshal(data, &token); err == nil {
				log.Infof("Loaded access token from cache for profile %s", profile.Name)
				log.Record("token", token).Debugf("Access token details for profile %s", profile.Name)
				profile.token = &token
				return nil
			}
		}
		// Load the access token from the vault in case this is an API Token
		log.Debugf("Looking for access token in the vault for profile %s", profile.Name)
		if credential, err := profile.GetCredentialFromVault(profile.VaultKey, profile.Name); err == nil {
			profile.AccessToken = credential.Password
			log.Infof("Loaded Repository/Project/Workspace Access Token for profile %s from the vault", profile.Name)
			profile.token = &Token{
				AccessToken: profile.AccessToken,
				ExpiresOn:   core.Timestamp(time.Now().Add(100 * 365 * 24 * time.Hour)), // Loaded Access Tokens never expire
			}
			return nil
		} else {
			log.Errorf("failed to get access token for profile %s: %v", profile.Name, err)
			return nil // We don't return an error if the token is not found, so the authorization process can continue
		}
	}
	return err
}

// isTokenExpired tells if the token is expired
func (profile *Profile) isTokenExpired() bool {
	return profile.token != nil && profile.token.IsExpired()
}

// saveAccessToken saves the access token to the cache
func (profile *Profile) saveAccessToken(ctx context.Context, data []byte) (accessToken string, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("profile", "saveAccessToken")

	profile.token, err = UnmarshalTokenFromBitbucketData(data)
	if err != nil {
		log.Errorf("Failed to unmarshal access token data for profile %s: %v", profile.Name, err)
		return "", err
	}

	if cacheDir, err := os.UserCacheDir(); err == nil {
		cachePath := filepath.Join(cacheDir, "bitbucket")
		if err = os.MkdirAll(cachePath, 0700); err == nil {
			cacheFile := filepath.Join(cachePath, "access-token-"+profile.Name)
			payload, _ := json.Marshal(profile.token)
			if err = os.WriteFile(cacheFile, payload, 0600); err != nil {
				log.Errorf("Failed to save access token to cache for profile %s", profile.Name, err)
			}
		}
	}
	return profile.token.AccessToken, nil
}

// Redact redacts sensitive information from the token
//
// implements logger.Redactable
func (token Token) Redact() any {
	redacted := token
	if len(redacted.AccessToken) > 0 {
		redacted.AccessToken = logger.RedactWithHash(redacted.AccessToken)
	}
	if len(redacted.RefreshToken) > 0 {
		redacted.RefreshToken = logger.RedactWithHash(redacted.RefreshToken)
	}
	return redacted
}

// IsExpired tells if the token is expired
func (token *Token) IsExpired() bool {
	return time.Time(token.ExpiresOn).Before(time.Now())
}

// GetExpiresOn returns the expiration time of the token
func (token *Token) GetExpiresOn() time.Time {
	return time.Time(token.ExpiresOn)
}

// GetExpiresIn returns the duration until the token expires
func (token *Token) GetExpiresIn() time.Duration {
	return time.Until(time.Time(token.ExpiresOn))
}

// GetExpiredSince returns the duration since the token expired
func (token *Token) GetExpiredSince() time.Duration {
	if !token.IsExpired() {
		return 0
	}
	return time.Since(time.Time(token.ExpiresOn))
}

// GetScopes returns the scopes of the token
func (token *Token) GetScopes() []string {
	return strings.Split(token.Scope, " ")
}

// UnmarshalTokenFromBitbucketData unmarshals the token data from the BitBucket response
func UnmarshalTokenFromBitbucketData(data []byte) (token *Token, err error) {
	var result struct {
		TokenType    string `json:"token_type"`
		State        string `json:"state"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		Scopes       string `json:"scopes"`
	}
	if err = json.Unmarshal(data, &result); err == nil {
		token = &Token{
			TokenType:    result.TokenType,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			ExpiresOn:    core.Timestamp(time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)),
			Scope:        result.Scopes,
		}
	}
	return token, errors.JSONUnmarshalError.WrapIfNotMe(err)
}
