package profile

import (
	"github.com/gildas/go-errors"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/zalando/go-keyring"
)

// Credential represents a user credential for authentication.
type Credential struct {
	Username string
	Password string
}

// AsHTTPBasicAuth returns the credential as an HTTP BasicAuth structure.
func (credential *Credential) AsHTTPBasicAuth() *http.BasicAuth {
	return &http.BasicAuth{
		Username: credential.Username,
		Password: credential.Password,
	}
}

// GetCredentialFromVault retrieves the credential for the given key from the Linux keyring.
func (profile Profile) GetCredentialFromVault(service, username string) (credential *Credential, err error) {
	secret, err := keyring.Get(service, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get secret from keyring")
	}
	if secret == "" {
		return nil, errors.NotFound.With("key", service)
	}
	return &Credential{Username: username, Password: secret}, nil
}
