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

// GetCredentialFromVault retrieves the credential for the given key from the Windows Credential Manager or Linux/macOS keychain.
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

// SetCredentialInVault stores the credential in the Windows Credential Manager or Linux/macOS keychain.
func (profile Profile) SetCredentialInVault(service, username, password string) error {
	if err := keyring.Set(service, username, password); err != nil {
		return errors.Wrap(err, "failed to set secret in keyring")
	}
	return nil
}

// DeleteCredentialFromVault removes the credential from the Windows Credential Manager or Linux/macOS keychain.
func (profile Profile) DeleteCredentialFromVault(service, username string) error {
	if err := keyring.Delete(service, username); err != nil {
		return errors.Wrap(err, "failed to delete secret from keyring")
	}
	return nil
}
