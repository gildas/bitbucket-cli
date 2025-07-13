package profile

import (
	"context"
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Profiles is a collection of Profile
type profiles []*Profile

// Profiles is the collection of profiles
var Profiles profiles

// Current gets the current profile
func (profiles profiles) Current(context context.Context) *Profile {
	log := logger.Must(logger.FromContext(context)).Child("profile", "current")

	if gitConfig, err := common.OpenGitConfig(context); err == nil {
		log.Debugf("Found a git config file")
		if section, err := common.GetGitSection(context, gitConfig, `bitbucket "cli"`); err == nil {
			log.Debugf("Found a bitbucket \"cli\" section in git config: %v", section)
			if profileName := section.Key("profile").String(); len(profileName) > 0 {
				log.Debugf("Found a profile in git config: %s", profileName)
				if profile, found := profiles.Find(profileName); found {
					log.Infof("Using profile %s from git config", profileName)
					return profile
				} else {
					log.Warnf("Profile %s not found in %s", profileName, viper.ConfigFileUsed())
					fmt.Fprintf(os.Stderr, "Profile %s from your git config was not found in %s, ignored.\n", profileName, viper.ConfigFileUsed())
				}
			}
		}
	}

	for _, profile := range profiles {
		if profile.Default {
			return profile
		}
	}
	if len(profiles) > 0 {
		return profiles[0]
	}
	return &Profile{
		Name: "default",
	}
}

// Names gets the names of the profiles
func (profiles profiles) Names() []string {
	names := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	return names
}

// GetHeader gets the header for a table
//
// implements common.Tableables
func (profiles profiles) GetHeader() []string {
	return Profile{}.GetHeader(false)
}

// GetRowAt gets the row for a table
//
// implements common.Tableables
func (profiles profiles) GetRowAt(index int, headers []string) []string {
	if index < 0 || index >= len(profiles) {
		return []string{}
	}
	profile := profiles[index]
	return []string{
		profile.Name,
		profile.Description,
		fmt.Sprintf("%v", profile.Default),
		profile.User,
		profile.ClientID,
		profile.AccessToken,
	}
}

// Size gets the number of elements
//
// implements common.Tableables
func (profiles profiles) Size() int {
	return len(profiles)
}

// Find finds a profile by name
func (profiles profiles) Find(name string) (profile *Profile, found bool) {
	for _, profile = range profiles {
		if profile.Name == name {
			return profile, true
		}
	}
	return nil, false
}

// Add adds a profile to the collection
func (profiles *profiles) Add(profile *Profile) {
	*profiles = append(*profiles, profile)
	if profile.Default {
		profiles.SetCurrent(profile.Name)
	}
}

// Delete deletes one or more profiles by their names
func (profiles *profiles) Delete(names ...string) (deleted int) {
	for _, name := range names {
		for i, profile := range *profiles {
			if profile.Name == name {
				deleted++
				*profiles = append((*profiles)[:i], (*profiles)[i+1:]...)
				break
			}
		}
	}
	return deleted
}

// SetCurrent sets the current profile
func (profiles profiles) SetCurrent(name string) {
	if len(name) == 0 {
		return
	}
	if _, found := profiles.Find(name); !found {
		return
	}
	for _, profile := range profiles {
		if profile.Name == name {
			profile.Default = true
		} else {
			profile.Default = false
		}
	}
}

// Load loads the profiles from a viper key
func (profiles *profiles) Load(context context.Context) error {
	log := logger.Must(logger.FromContext(context)).Child("profiles", "load")

	log.Infof("Loading profiles from %s", viper.ConfigFileUsed())
	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		return err
	}
	log.Debugf("Loaded %d profiles", len(*profiles))

	// Get the secret stuff from the Windows credential manager or linux/macOS keychain if not set
	for _, profile := range *profiles {
		if len(profile.ClientID) > 0 {
			if len(profile.ClientSecret) == 0 {
				if credential, err := profile.GetCredentialFromVault("bitbucket-cli", profile.ClientID); err == nil {
					profile.ClientSecret = credential.Password
					log.Infof("Loaded client secret for clientID %s from the vault", profile.ClientID)
				} else {
					log.Errorf("failed to get client secret for profile %s: %v", profile.Name, err)
				}
			}
		} else if len(profile.User) > 0 {
			if len(profile.Password) == 0 {
				if credential, err := profile.GetCredentialFromVault("bitbucket-cli", profile.User); err == nil {
					profile.Password = credential.Password
					log.Infof("Loaded password for user %s from the vault", profile.User)
				} else {
					log.Errorf("failed to get password for profile %s: %v", profile.Name, err)
				}
			}
		} else if len(profile.AccessToken) == 0 {
			if credential, err := profile.GetCredentialFromVault("bitbucket-cli", profile.Name); err == nil {
				profile.AccessToken = credential.Password
				log.Infof("Loaded access token for profile %s from the vault", profile.Name)
			} else {
				log.Errorf("failed to get access token for profile %s: %v", profile.Name, err)
			}
		}
	}
	return nil
}

// ValidProfileNames gets the valid profile names
func ValidProfileNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	names := Profiles.Names()
	return common.FilterValidArgs(names, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}
