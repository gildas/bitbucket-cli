package profile

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/gildas/bitbucket-cli/cmd/common"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/kataras/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Profile describes the configuration needed to connect to BitBucket
type Profile struct {
	Name              string                 `json:"name"                        mapstructure:"name"`
	Description       string                 `json:"description,omitempty"       mapstructure:"description,omitempty"       yaml:",omitempty"`
	Default           bool                   `json:"default"                     mapstructure:"default"                     yaml:",omitempty"`
	APIRoot           *url.URL               `json:"apiRoot,omitempty"           mapstructure:"apiRoot,omitempty"           yaml:",omitempty"`
	DefaultWorkspace  string                 `json:"defaultWorkspace,omitempty"  mapstructure:"defaultWorkspace,omitempty"  yaml:",omitempty"`
	DefaultProject    string                 `json:"defaultProject,omitempty"    mapstructure:"defaultProject,omitempty"    yaml:",omitempty"`
	ErrorProcessing   common.ErrorProcessing `json:"errorProcessing,omitempty"   mapstructure:"errorProcessing,omitempty"   yaml:",omitempty"`
	DefaultPageLength int                    `json:"defaultPageLength,omitempty" mapstructure:"defaultPageLength,omitempty" yaml:",omitempty"`
	OutputFormat      string                 `json:"outputFormat,omitempty"      mapstructure:"outputFormat,omitempty"      yaml:",omitempty"`
	Progress          bool                   `json:"progress,omitempty"          mapstructure:"progress,omitempty"          yaml:",omitempty"`
	CloneProtocol     string                 `json:"cloneProtocol,omitempty"     mapstructure:"cloneProtocol,omitempty"     yaml:",omitempty"`
	CloneUser         string                 `json:"cloneUser,omitempty"         mapstructure:"cloneUser,omitempty"         yaml:",omitempty"`
	SshKeyFilename    string                 `json:"sshKeyFilename,omitempty"    mapstructure:"sshKeyFilename,omitempty"    yaml:",omitempty"`
	VaultKey          string                 `json:"vaultKey,omitempty"          mapstructure:"vaultKey,omitempty"          yaml:",omitempty"`
	User              string                 `json:"user,omitempty"              mapstructure:"user"                        yaml:",omitempty"`
	Password          string                 `json:"password,omitempty"          mapstructure:"password"                    yaml:",omitempty"`
	ClientID          string                 `json:"clientID,omitempty"          mapstructure:"clientID"                    yaml:",omitempty"`
	ClientSecret      string                 `json:"clientSecret,omitempty"      mapstructure:"clientSecret"                yaml:",omitempty"`
	CallbackPort      uint16                 `json:"callbackPort,omitempty"      mapstructure:"callbackPort"                yaml:",omitempty"`
	AccessToken       string                 `json:"accessToken,omitempty"       mapstructure:"accessToken,omitempty"       yaml:",omitempty"`
	token             *Token                 `json:"-"                           mapstructure:"-"                           yaml:"-"`
}

// Current is the current profile
var Current *Profile

const (
	DefaultPageLength = 50 // DefaultPageLength is the default number of items per page to retrieve from Bitbucket
)

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

var columns = common.Columns[*Profile]{
	{Name: "name", DefaultSorter: true, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) == -1
	}},
	{Name: "description", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.Description), strings.ToLower(b.Description)) == -1
	}},
	{Name: "default", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return a.Default == b.Default
	}},
	{Name: "user", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.User), strings.ToLower(b.User)) == -1
	}},
	{Name: "clientid", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.ClientID), strings.ToLower(b.ClientID)) == -1
	}},
	{Name: "accesstoken", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.AccessToken), strings.ToLower(b.AccessToken)) == -1
	}},
	{Name: "apiRoot", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return a.APIRoot != nil && b.APIRoot != nil && strings.Compare(strings.ToLower(a.APIRoot.String()), strings.ToLower(b.APIRoot.String())) == -1
	}},
	{Name: "defaultworkspace", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.DefaultWorkspace), strings.ToLower(b.DefaultWorkspace)) == -1
	}},
	{Name: "defaultproject", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.DefaultProject), strings.ToLower(b.DefaultProject)) == -1
	}},
	{Name: "callbackport", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return a.CallbackPort < b.CallbackPort
	}},
	{Name: "outputformat", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.OutputFormat), strings.ToLower(b.OutputFormat)) == -1
	}},
	{Name: "defaultpagelength", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return a.DefaultPageLength < b.DefaultPageLength
	}},
	{Name: "cloneprotocol", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.CloneProtocol), strings.ToLower(b.CloneProtocol)) == -1
	}},
	{Name: "cloneuser", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.CloneUser), strings.ToLower(b.CloneUser)) == -1
	}},
	{Name: "sshkeyfilename", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.SshKeyFilename), strings.ToLower(b.SshKeyFilename)) == -1
	}},
	{Name: "vaultkey", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.VaultKey), strings.ToLower(b.VaultKey)) == -1
	}},
	{Name: "errorprocessing", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return strings.Compare(strings.ToLower(a.ErrorProcessing.String()), strings.ToLower(b.ErrorProcessing.String())) == -1
	}},
	{Name: "progress", DefaultSorter: false, Compare: func(a, b *Profile) bool {
		return a.Progress == b.Progress
	}},
}

// GetProfileFromCommand gets the profile from the command line
//
// If the profile is not given, it will use the current profile
func GetProfileFromCommand(context context.Context, cmd *cobra.Command) (profile *Profile, err error) {
	log := logger.Must(logger.FromContext(context)).Child("profile", "getProfileFromCommand")

	if err = Profiles.Load(context, cmd); err != nil {
		return nil, err
	}

	if cmd.Flag("profile").Changed {
		var found bool
		log.Debugf("Command line has profile flag set to %s", cmd.Flag("profile").Value.String())
		if profile, found = Profiles.Find(cmd.Flag("profile").Value.String()); !found {
			return nil, errors.ArgumentInvalid.With("profile", cmd.Flag("profile").Value.String())
		}
	} else if Current == nil {
		Current = Profiles.Current(context)
		if Current == nil {
			return nil, errors.ArgumentMissing.With("profile")
		}
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
			if profile.APIRoot != nil {
				row = append(row, profile.APIRoot.String())
			} else {
				row = append(row, " ")
			}
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
		case "outputformat":
			row = append(row, profile.OutputFormat)
		case "defaultpagelength":
			row = append(row, fmt.Sprintf("%d", profile.DefaultPageLength))
		case "cloneprotocol":
			row = append(row, profile.CloneProtocol)
		case "cloneuser":
			row = append(row, profile.CloneUser)
		case "sshkeyfilename":
			row = append(row, profile.SshKeyFilename)
		case "vaultkey":
			row = append(row, profile.VaultKey)
		case "errorprocessing":
			row = append(row, profile.ErrorProcessing.String())
		case "progress":
			row = append(row, fmt.Sprintf("%t", profile.Progress))
		default:
			row = append(row, " ")
		}
	}
	return row
}

// Redact redacts sensitive information from the profile
//
// implements logger.Redactable
func (profile Profile) Redact() any {
	redacted := profile
	if len(redacted.ClientID) > 0 {
		redacted.ClientID = logger.RedactWithHash(redacted.ClientID)
	}
	if len(redacted.ClientSecret) > 0 {
		redacted.ClientSecret = logger.RedactWithHash(redacted.ClientSecret)
	}
	if len(redacted.User) > 0 {
		redacted.User = logger.RedactWithHash(redacted.User)
	}
	if len(redacted.Password) > 0 {
		redacted.Password = logger.RedactWithHash(redacted.Password)
	}
	if len(redacted.AccessToken) > 0 {
		redacted.AccessToken = logger.RedactWithHash(redacted.AccessToken)
	}
	if len(redacted.CloneUser) > 0 {
		redacted.CloneUser = logger.RedactWithHash(redacted.CloneUser)
	}
	return redacted
}

// GetClientSecret gets the client secret from the profile, either from the vault or from the profile
func (profile *Profile) GetClientSecret(ctx context.Context) (clientSecret string, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("profile", "getClientSecret")
	if len(profile.ClientSecret) > 0 {
		log.Debugf("Client secret for profile %s is set in the profile", profile.Name)
		return profile.ClientSecret, nil
	}
	if credential, err := profile.GetCredentialFromVault(profile.VaultKey, profile.ClientID); err == nil {
		log.Debugf("Loaded client secret for clientID %s from the vault", profile.ClientID)
		return credential.Password, nil
	}
	return "", errors.Join(errors.Errorf("Profile %s does not have a client secret", profile.Name), err)
}

// GetPassword gets the password from the profile, either from the vault or from the profile
func (profile *Profile) GetPassword(ctx context.Context) (password string, err error) {
	log := logger.Must(logger.FromContext(ctx)).Child("profile", "getPassword")
	if len(profile.Password) > 0 {
		log.Debugf("Password for profile %s is set in the profile", profile.Name)
		return profile.Password, nil
	}
	if credential, err := profile.GetCredentialFromVault(profile.VaultKey, profile.User); err == nil {
		log.Debugf("Loaded password for user %s from the vault", profile.User)
		return credential.Password, nil
	}
	return "", errors.Join(errors.Errorf("Profile %s does not have a password", profile.Name), err)
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
	if len(other.AccessToken) > 0 && other.AccessToken != profile.AccessToken {
		profile.AccessToken = other.AccessToken
	}
	if len(other.User) > 0 && other.User != profile.User {
		profile.User = other.User
	}
	if len(other.Password) > 0 && other.Password != profile.Password {
		profile.Password = other.Password
	}
	if len(other.ClientID) > 0 && other.ClientID != profile.ClientID {
		profile.ClientID = other.ClientID
		profile.token = nil
	}
	if len(other.ClientSecret) > 0 && other.ClientSecret != profile.ClientSecret {
		profile.ClientSecret = other.ClientSecret
		profile.token = nil
	}
	if other.CallbackPort > 0 {
		profile.CallbackPort = other.CallbackPort
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
	if len(other.SshKeyFilename) > 0 {
		profile.SshKeyFilename = other.SshKeyFilename
	}
	return profile.Validate()
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
		table.SetAutoWrapText(false)
		table.Append(actual.GetRow(headers))
	case common.Tableables:
		log.Debugf("Payload is a slice of %d elements", actual.Size())
		if actual.Size() > 0 {
			headers := actual.GetHeaders(cmd)
			table.SetHeader(headers)
			table.SetAutoWrapText(false)
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

// Validate validates a Profile
func (profile *Profile) Validate() error {
	var merr errors.MultiError

	if len(profile.Name) == 0 {
		merr.Append(errors.ArgumentMissing.With("name"))
	}

	if len(profile.VaultKey) == 0 && runtime.GOOS != "windows" {
		profile.VaultKey = "bitbucket-cli"
	}

	if len(profile.CloneProtocol) == 0 {
		profile.CloneProtocol = "git"
	}
	if profile.CloneProtocol != "git" && profile.CloneProtocol != "https" && profile.CloneProtocol != "ssh" {
		merr.Append(errors.ArgumentInvalid.With("cloneProtocol", profile.CloneProtocol))
	}
	if len(profile.OutputFormat) == 0 {
		profile.OutputFormat = "table"
	}
	if profile.DefaultPageLength == 0 {
		profile.DefaultPageLength = DefaultPageLength
	} else if profile.DefaultPageLength < 0 || profile.DefaultPageLength > 100 {
		merr.Append(errors.Errorf("Default Page Length must be between 0 and 100 (value: %d)", profile.DefaultPageLength))
	}
	return merr.AsError()
}

// MarshalJSON marshals this profile to JSON
//
// implements json.Marshaler
func (profile Profile) MarshalJSON() ([]byte, error) {
	type surrogate Profile

	if profile.OutputFormat == "table" {
		profile.OutputFormat = ""
	}
	if profile.DefaultPageLength == DefaultPageLength {
		profile.DefaultPageLength = 0
	}
	errorProcessing := profile.ErrorProcessing.String()
	if errorProcessing == common.StopOnError.String() {
		errorProcessing = ""
	}
	data, err := json.Marshal(struct {
		surrogate
		APIRoot         *core.URL `json:"apiRoot,omitempty"`
		ErrorProcessing string    `json:"errorProcessing,omitempty"`
	}{
		surrogate:       surrogate(profile),
		APIRoot:         (*core.URL)(profile.APIRoot),
		ErrorProcessing: errorProcessing,
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals this profile from JSON
//
// implements json.Unmarshaler
func (profile *Profile) UnmarshalJSON(data []byte) error {
	type surrogate Profile
	var inner struct {
		surrogate
		APIRoot *core.URL `json:"apiRoot,omitempty"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*profile = Profile(inner.surrogate)
	profile.APIRoot = (*url.URL)(inner.APIRoot)
	return errors.JSONUnmarshalError.Wrap(profile.Validate())
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

// disableUnsupportedFlags disables the flags that are not supported by the profile command
func disableUnsupportedFlags(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("repository") {
		return fmt.Errorf("the --repository flag is not supported by the profile command")
	}
	if cmd.Flags().Changed("workspace") {
		return fmt.Errorf("the --workspace flag is not supported by the profile command")
	}
	return nil
}

// hideUnsupportedFlags hides the flags that are not supported by the profile command
func hideUnsupportedFlags(cmd *cobra.Command, args []string) {
	cmd.Flags().MarkHidden("repository")
	cmd.Flags().MarkHidden("workspace")
	cmd.Parent().HelpFunc()(cmd, args)
}
