package profile

import (
	"github.com/spf13/viper"
)

// Profiles is a collection of Profile
type profiles []*Profile

// Profiles is the collection of profiles
var Profiles profiles

// Current gets the current profile
func (profiles profiles) Current() *Profile {
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
func (profiles *profiles) Load() error {

	if err := viper.UnmarshalKey("profiles", &profiles); err != nil {
		return err
	}

	if len(*profiles) == 0 {
		Log.Warnf("No profiles found in config file %s", viper.GetViper().ConfigFileUsed())
	}
	return nil
}
