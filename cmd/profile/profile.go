package profile

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Profile describes the configuration needed to connect to BitBucket
type Profile struct {
	Name             string                 `json:"name"                       mapstructure:"name"`
	Description      string                 `json:"description,omitempty"      mapstructure:"description,omitempty"     yaml:",omitempty"`
	Default          bool                   `json:"default"                    mapstructure:"default"                   yaml:",omitempty"`
	APIRoot          *url.URL               `json:"apiRoot,omitempty"         mapstructure:"apiRoot,omitempty"          yaml:",omitempty"`
	DefaultWorkspace string                 `json:"defaultWorkspace,omitempty" mapstructure:"defaultWorkspace"          yaml:",omitempty"`
	DefaultProject   string                 `json:"defaultProject,omitempty"   mapstructure:"defaultProject"            yaml:",omitempty"`
	ErrorProcessing  common.ErrorProcessing `json:"errorProcessing,omitempty"  mapstructure:"errorProcessing,omitempty" yaml:",omitempty"`
	OutputFormat     string                 `json:"outputFormat,omitempty"     mapstructure:"outputFormat,omitempty"    yaml:",omitempty"`
	Progress         bool                   `json:"progress,omitempty"         mapstructure:"progress,omitempty"        yaml:",omitempty"`
	VaultKey         string                 `json:"vaultKey,omitempty"         mapstructure:"vaultKey,omitempty"        yaml:",omitempty"`
	User             string                 `json:"user,omitempty"             mapstructure:"user"                      yaml:",omitempty"`
	Password         string                 `json:"password,omitempty"         mapstructure:"password"                  yaml:",omitempty"`
	ClientID         string                 `json:"clientID,omitempty"         mapstructure:"clientID"                  yaml:",omitempty"`
	ClientSecret     string                 `json:"clientSecret,omitempty"     mapstructure:"clientSecret"              yaml:",omitempty"`
	CallbackPort     uint16                 `json:"callbackPort,omitempty"     mapstructure:"callbackPort"              yaml:",omitempty"`
	AccessToken      string                 `json:"accessToken,omitempty"      mapstructure:"accessToken"               yaml:",omitempty"`
	RefreshToken     string                 `json:"-"                          mapstructure:"refreshToken"              yaml:"-"`
	TokenExpires     time.Time              `json:"-"                          mapstructure:"tokenExpires"              yaml:"-"`
	TokenScopes      []string               `json:"-"                          mapstructure:"tokenScopes"               yaml:"-"`
	CloneProtocol    string                 `json:"cloneProtocol,omitempty"    mapstructure:"cloneProtocol,omitempty"   yaml:",omitempty"`
	CloneUser        string                 `json:"cloneUser,omitempty"        mapstructure:"cloneUser,omitempty"       yaml:",omitempty"`
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

var columns = []string{
	"name",
	"description",
	"default",
	"user",
	"clientid",
	"accesstoken",
	"tokenexpires",
	"apiRoot",
	"defaultworkspace",
	"defaultproject",
	"callbackPort",
}

// GetProfileFromCommand gets the profile from the command line
//
// If the profile is not given, it will use the current profile
func GetProfileFromCommand(context context.Context, cmd *cobra.Command) (profile *Profile, err error) {
	if cmd.Flag("profile").Changed {
		var found bool
		if profile, found = Profiles.Find(cmd.Flag("profile").Value.String()); !found {
			return nil, errors.ArgumentInvalid.With("profile", cmd.Flag("profile").Value.String())
		}
	} else if Current == nil {
		if len(Profiles) == 0 {
			err = Profiles.Load(context)
			if err != nil {
				return nil, err
			}
		}
		Current = Profiles.Current(context)
		profile = Current
	} else {
		profile = Current
	}
	return
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (profile Profile) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"Name", "Description", "Default", "User", "ClientID", "AccessToken"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (profile Profile) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "apiroot":
			row = append(row, profile.APIRoot.String())
		case "name":
			row = append(row, profile.Name)
		case "description":
			row = append(row, profile.Description)
		case "default":
			row = append(row, fmt.Sprintf("%v", profile.Default))
		case "defaultworkspace":
			row = append(row, profile.DefaultWorkspace)
		case "defaultproject":
			row = append(row, profile.DefaultProject)
		case "callbackport":
			row = append(row, fmt.Sprintf("%d", profile.CallbackPort))
		case "user":
			row = append(row, profile.User)
		case "clientid":
			row = append(row, profile.ClientID)
		case "accesstoken":
			if len(profile.AccessToken) > 0 {
				row = append(row, profile.AccessToken)
			} else {
				row = append(row, " ")
			}
		case "tokenexpires":
			if !profile.TokenExpires.IsZero() {
				row = append(row, profile.TokenExpires.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
}

// Redact redacts sensitive information from the profile
//
// implements logger.Redactable
func (profile Profile) Redact() any {
	redacted := profile
	if len(redacted.ClientSecret) > 0 {
		redacted.ClientSecret = logger.RedactWithHash(redacted.ClientSecret)
	}
	if len(redacted.Password) > 0 {
		redacted.Password = logger.RedactWithHash(redacted.Password)
	}
	if len(redacted.AccessToken) > 0 {
		redacted.AccessToken = logger.RedactWithHash(redacted.AccessToken)
	}
	if len(redacted.RefreshToken) > 0 {
		redacted.RefreshToken = logger.RedactWithHash(redacted.RefreshToken)
	}
	redacted.TokenScopes = nil
	return redacted
}

// Update updates this profile with the given one
func (profile *Profile) Update(other Profile) error {
	if len(other.Name) > 0 {
		profile.Name = other.Name
	}
	if len(other.Description) > 0 {
		profile.Description = other.Description
	}
	if other.Default {
		profile.Default = other.Default
	}
	if len(other.OutputFormat) > 0 {
		profile.OutputFormat = other.OutputFormat
	}
	if len(other.User) > 0 {
		profile.User = other.User
	}
	if len(other.Password) > 0 {
		profile.Password = other.Password
	}
	if len(other.ClientID) > 0 {
		profile.ClientID = other.ClientID
		profile.RefreshToken = ""
		profile.TokenExpires = time.Time{}
		profile.TokenScopes = []string{}
	}
	if len(other.ClientSecret) > 0 {
		profile.ClientSecret = other.ClientSecret
	}
	if other.CallbackPort > 0 {
		profile.CallbackPort = other.CallbackPort
	}
	if len(other.AccessToken) > 0 {
		profile.AccessToken = other.AccessToken
		profile.RefreshToken = ""
		profile.TokenExpires = time.Time{}
		profile.TokenScopes = []string{}
	}
	if len(other.DefaultWorkspace) > 0 {
		profile.DefaultWorkspace = other.DefaultWorkspace
	}
	if len(other.DefaultProject) > 0 {
		profile.DefaultProject = other.DefaultProject
	}
	if len(other.CloneProtocol) > 0 {
		profile.CloneProtocol = other.CloneProtocol
	}
	if len(other.CloneUser) > 0 {
		profile.CloneUser = other.CloneUser
	}
	return profile.Validate()
}

// Validate validates a Profile
func (profile *Profile) Validate() error {
	var merr errors.MultiError

	if len(profile.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("name"))
	}
	// We must have either an access token, a user, or a clientID
	// password and clientSecret are now retrieved from the vault
	if len(profile.AccessToken) == 0 && len(profile.ClientID) == 0 && len(profile.User) == 0 {
		merr.Append(errors.ArgumentMissing.With("accessToken, user, or clientID"))
	}
	if len(profile.VaultKey) == 0 {
		profile.VaultKey = "bitbucket-cli"
	}
	if len(profile.CloneProtocol) == 0 {
		profile.CloneProtocol = "git"
	}
	if profile.CloneProtocol != "git" && profile.CloneProtocol != "https" && profile.CloneProtocol != "ssh" {
		merr.Append(errors.ArgumentInvalid.With("cloneProtocol", profile.CloneProtocol))
	}
	return merr.AsError()
}

// ShouldStopOnError tells if the command should stop on error
func (profile Profile) ShouldStopOnError(cmd *cobra.Command) bool {
	if cmd.Flag("stop-on-error").Changed {
		return cmd.Flag("stop-on-error").Value.String() == "true"
	}
	return profile.ErrorProcessing == common.StopOnError
}

// ShouldWarnOnError tells if the command should warn on error
func (profile Profile) ShouldWarnOnError(cmd *cobra.Command) bool {
	if cmd.Flag("warn-on-error").Changed {
		return cmd.Flag("warn-on-error").Value.String() == "true"
	}
	return profile.ErrorProcessing == common.WarnOnError
}

// ShouldIgnoreErrors tells if the command should ignore errors
func (profile Profile) ShouldIgnoreErrors(cmd *cobra.Command) bool {
	if cmd.Flag("ignore-errors").Changed {
		return cmd.Flag("ignore-errors").Value.String() == "true"
	}
	return profile.ErrorProcessing == common.IgnoreErrors
}

// String gets a string representation of this profile
//
// implements fmt.Stringer
func (profile Profile) String() string {
	return profile.Name
}

// Print prints the given payload to the console
func (profile Profile) Print(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context)).Child("profile", "print", "format", profile.OutputFormat)
	outputFormat := profile.OutputFormat

	if cmd.Flag("output").Changed {
		outputFormat = cmd.Flag("output").Value.String()
		log.Debugf("Command output format: %s (was: %s)", outputFormat, profile.OutputFormat)
	}
	switch outputFormat {
	case "json":
		return profile.PrintJSON(context, cmd, payload)
	case "yaml":
		return profile.PrintYAML(context, cmd, payload)
	case "csv":
		return profile.PrintCSV(context, cmd, payload)
	case "tsv":
		return profile.PrintTSV(context, cmd, payload)
	case "table":
		fallthrough
	default:
		return profile.PrintTable(context, cmd, payload)
	}
}

// PrintJSON prints the given payload to the console as JSON
func (profile Profile) PrintJSON(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context))

	log.Debugf("Printing payload as JSON")
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return errors.JSONMarshalError.Wrap(err)
	}
	fmt.Println(string(data))
	return nil
}

// PrintYAML prints the given payload to the console as YAML
func (profile Profile) PrintYAML(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context))

	log.Debugf("Printing payload as YAML")
	data, err := yaml.Marshal(payload)
	if err != nil {
		return errors.JSONMarshalError.Wrap(err)
	}
	fmt.Println(string(data))
	return nil
}

// PrintCSV prints the given payload to the console as CSV
func (profile Profile) PrintCSV(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context))

	log.Debugf("Printing payload as CSV")
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	switch actual := payload.(type) {
	case common.Tableable:
		headers := actual.GetHeaders(cmd)
		_ = writer.Write(headers)
		_ = writer.Write(actual.GetRow(headers))
	case common.Tableables:
		log.Debugf("Payload is a slice of %d elements", actual.Size())
		if actual.Size() > 0 {
			headers := actual.GetHeaders(cmd)
			_ = writer.Write(headers)
			for i := 0; i < actual.Size(); i++ {
				_ = writer.Write(actual.GetRowAt(i, headers))
			}
		}
	default:
		return errors.ArgumentInvalid.With("payload", "not a tableable")
	}
	return nil
}

// PrintTSV prints the given payload to the console as TSV
func (profile Profile) PrintTSV(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context))

	log.Debugf("Printing payload as TSV")
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'
	defer writer.Flush()

	switch actual := payload.(type) {
	case common.Tableable:
		headers := actual.GetHeaders(cmd)
		_ = writer.Write(headers)
		_ = writer.Write(actual.GetRow(headers))
	case common.Tableables:
		log.Debugf("Payload is a slice of %d elements", actual.Size())
		if actual.Size() > 0 {
			headers := actual.GetHeaders(cmd)
			_ = writer.Write(headers)
			for i := 0; i < actual.Size(); i++ {
				_ = writer.Write(actual.GetRowAt(i, headers))
			}
		}
	default:
		return errors.ArgumentInvalid.With("payload", "not a tableable")
	}
	return nil
}

// PrintTable prints the given payload to the console as a table
func (profile Profile) PrintTable(context context.Context, cmd *cobra.Command, payload any) error {
	log := logger.Must(logger.FromContext(context))

	log.Debugf("Printing payload as table")
	table := tablewriter.NewWriter(os.Stdout)

	switch actual := payload.(type) {
	case common.Tableable:
		headers := actual.GetHeaders(cmd)
		table.SetHeader(headers)
		table.Append(actual.GetRow(headers))
	case common.Tableables:
		log.Debugf("Payload is a slice of %d elements", actual.Size())
		if actual.Size() > 0 {
			headers := actual.GetHeaders(cmd)
			table.SetHeader(headers)
			for i := 0; i < actual.Size(); i++ {
				table.Append(actual.GetRowAt(i, headers))
			}
		}
	default:
		return errors.ArgumentInvalid.With("payload", "not a tableable")
	}
	table.Render()
	return nil
}

// MarshalJSON marshals this profile to JSON
//
// implements json.Marshaler
func (profile Profile) MarshalJSON() ([]byte, error) {
	type surrogate Profile
	outputFormat := profile.OutputFormat
	if outputFormat == "table" {
		outputFormat = ""
	}
	errorProcessing := profile.ErrorProcessing.String()
	if errorProcessing == common.StopOnError.String() {
		errorProcessing = ""
	}
	data, err := json.Marshal(struct {
		surrogate
		OutputFormat    string `json:"outputFormat,omitempty"`
		ErrorProcessing string `json:"errorProcessing,omitempty"`
	}{
		surrogate:       surrogate(profile),
		OutputFormat:    outputFormat,
		ErrorProcessing: errorProcessing,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals this profile from JSON
//
// implements json.Unmarshaler
func (profile *Profile) UnmarshalJSON(data []byte) error {
	type surrogate Profile
	var inner surrogate
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*profile = Profile(inner)
	if len(profile.OutputFormat) == 0 {
		profile.OutputFormat = "table"
	}
	return nil
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

// getWorkspaceSlugs gets the slugs of all workspaces
func getWorkspaceSlugs(context context.Context, cmd *cobra.Command, args []string, toComplete string) (slugs []string, err error) {
	// We have to repeat the code here because of the circular dependency with the workspace package
	log := logger.Must(logger.FromContext(context)).Child("workspace", "slugs")
	type Workspace struct {
		Slug string `json:"slug"`
	}

	log.Debugf("Getting all workspaces")
	workspaces, err := GetAll[Workspace](context, cmd, "/workspaces")
	if err != nil {
		log.Errorf("Failed to get workspaces", err)
		return []string{}, err
	}
	slugs = core.Map(workspaces, func(workspace Workspace) string { return workspace.Slug })
	core.Sort(slugs, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return slugs, nil
}

// getProjectKeys gets the keys of all projects
func getProjectKeys(context context.Context, cmd *cobra.Command, args []string, toComplete string) (keys []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("project", "keys")
	type Project struct {
		Key string `json:"key"`
	}

	workspace := cmd.Flag("default-workspace").Value.String()
	if len(workspace) == 0 {
		log.Warnf("No workspace given")
		return
	}

	log.Debugf("Getting all projects in workspace %s", workspace)
	projects, err := GetAll[Project](context, cmd, fmt.Sprintf("/workspaces/%s/projects", workspace))
	if err != nil {
		log.Errorf("Failed to get projects", err)
		return
	}
	keys = core.Map(projects, func(project Project) string { return project.Key })
	core.Sort(keys, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return keys, nil
}
