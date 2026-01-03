package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

// Pipeline represents a Bitbucket Pipeline
type Pipeline struct {
	ID                   common.UUID           `json:"uuid"                            mapstructure:"uuid"`
	BuildNumber          uint64                `json:"build_number"                    mapstructure:"build_number"`
	State                PipelineState         `json:"state"                           mapstructure:"state"`
	Creator              user.User             `json:"creator"                         mapstructure:"creator"`
	Repository           Repository            `json:"repository"                      mapstructure:"repository"`
	Target               Target                `json:"target"                          mapstructure:"target"`
	Variables            []Variable            `json:"variables,omitempty"             mapstructure:"variables"`
	ConfigurationSources []ConfigurationSource `json:"configuration_sources,omitempty" mapstructure:"configuration_sources"`
	Duration             time.Duration         `json:"duration_in_seconds"             mapstructure:"duration_in_seconds"`
	CreatedOn            time.Time             `json:"created_on"                      mapstructure:"created_on"`
	CompletedOn          time.Time             `json:"completed_on"                    mapstructure:"completed_on"`
	Links                common.Links          `json:"links"                           mapstructure:"links"`
}

// PipelineState represents the state of a pipeline
type PipelineState struct {
	Type   string          `json:"type"             mapstructure:"type"`
	Name   string          `json:"name"             mapstructure:"name"`
	Stage  *PipelineStage  `json:"stage,omitempty"  mapstructure:"stage"`
	Result *PipelineResult `json:"result,omitempty" mapstructure:"result"`
}

// PipelineStage represents the current stage of a pipeline
type PipelineStage struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// PipelineResult represents the result of a completed pipeline
type PipelineResult struct {
	Type string `json:"type" mapstructure:"type"`
	Name string `json:"name" mapstructure:"name"`
}

// Repository represents a repository reference in a pipeline
type Repository struct {
	Type     string       `json:"type"      mapstructure:"type"`
	UUID     string       `json:"uuid"      mapstructure:"uuid"`
	Name     string       `json:"name"      mapstructure:"name"`
	FullName string       `json:"full_name" mapstructure:"full_name"`
	Links    common.Links `json:"links"     mapstructure:"links"`
}

// Variable represents a pipeline variable
type Variable struct {
	ID      common.UUID `json:"uuid"         mapstructure:"uuid"`
	Key     string      `json:"key"              mapstructure:"key"`
	Value   string      `json:"value"            mapstructure:"value"`
	Secured bool        `json:"secured"          mapstructure:"secured"`
}

// ConfigurationSource represents a pipeline configuration source
type ConfigurationSource struct {
	Source string `json:"source" mapstructure:"source"`
	URI    string `json:"uri"    mapstructure:"uri"`
}

// TriggerBody represents the body for triggering a pipeline
type TriggerBody struct {
	Target    Target     `json:"target"`
	Variables []Variable `json:"variables,omitempty"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:     "pipeline",
	Aliases: []string{"pipelines", "pipe", "pp"},
	Short:   "Manage pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pipeline requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

// Columns defines the available columns for the Pipeline type
var columns = common.Columns[Pipeline]{
	{Name: "build_number", DefaultSorter: true, Compare: func(a, b Pipeline) bool {
		return a.BuildNumber < b.BuildNumber
	}},
	{Name: "uuid", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return strings.Compare(strings.ToLower(a.ID.String()), strings.ToLower(b.ID.String())) == -1
	}},
	{Name: "state", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return strings.Compare(strings.ToLower(a.State.Name), strings.ToLower(b.State.Name)) == -1
	}},
	{Name: "result", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		aResult := ""
		bResult := ""
		if a.State.Result != nil {
			aResult = a.State.Result.Name
		}
		if b.State.Result != nil {
			bResult = b.State.Result.Name
		}
		return strings.Compare(strings.ToLower(aResult), strings.ToLower(bResult)) == -1
	}},
	{Name: "branch", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return strings.Compare(strings.ToLower(a.Target.RefName), strings.ToLower(b.Target.RefName)) == -1
	}},
	{Name: "commit", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		aHash := ""
		bHash := ""
		if a.Target.Commit != nil {
			aHash = a.Target.Commit.Hash
		}
		if b.Target.Commit != nil {
			bHash = b.Target.Commit.Hash
		}
		return strings.Compare(strings.ToLower(aHash), strings.ToLower(bHash)) == -1
	}},
	{Name: "creator", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return strings.Compare(strings.ToLower(a.Creator.Name), strings.ToLower(b.Creator.Name)) == -1
	}},
	{Name: "duration", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return a.Duration < b.Duration
	}},
	{Name: "creator", Compare: func(a, b Pipeline) bool {
		return a.Creator.Username < b.Creator.Username
	}},
	{Name: "created_on", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		return a.CreatedOn.Before(b.CreatedOn)
	}},
	{Name: "completed_on", DefaultSorter: false, Compare: func(a, b Pipeline) bool {
		if a.CompletedOn.IsZero() && b.CompletedOn.IsZero() {
			return false
		}
		if a.CompletedOn.IsZero() {
			return true
		}
		if b.CompletedOn.IsZero() {
			return false
		}
		return a.CompletedOn.Before(b.CompletedOn)
	}},
}

// GetType gets the type name
//
// implements core.TypeCarrier
func (pipeline Pipeline) GetType() string {
	return "pipeline"
}

// GetHeaders gets the header for a table
//
// implements common.Tableable
func (pipeline Pipeline) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	return []string{"build number", "state", "result", "branch", "creator", "duration", "created on"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (pipeline Pipeline) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "build number", "build_number":
			row = append(row, fmt.Sprintf("%d", pipeline.BuildNumber))
		case "uuid", "id":
			row = append(row, pipeline.ID.String())
		case "state":
			row = append(row, pipeline.State.Name)
		case "result":
			if pipeline.State.Result != nil {
				row = append(row, pipeline.State.Result.Name)
			} else {
				row = append(row, " ")
			}
		case "branch":
			row = append(row, pipeline.Target.RefName)
		case "commit":
			if pipeline.Target.Commit != nil {
				hash := pipeline.Target.Commit.Hash
				if len(hash) > 7 {
					hash = hash[:7]
				}
				row = append(row, hash)
			} else {
				row = append(row, " ")
			}
		case "creator":
			row = append(row, pipeline.Creator.Name)
		case "duration":
			row = append(row, pipeline.Duration.String())
		case "created on", "created_on":
			row = append(row, pipeline.CreatedOn.Format("2006-01-02 15:04:05"))
		case "completed on", "completed_on":
			if !pipeline.CompletedOn.IsZero() {
				row = append(row, pipeline.CompletedOn.Format("2006-01-02 15:04:05"))
			} else {
				row = append(row, " ")
			}
		}
	}
	return row
}

// Validate validates a Pipeline
func (pipeline *Pipeline) Validate() error {
	var merr errors.MultiError

	return merr.AsError()
}

// String gets a string representation of this pipeline
//
// implements fmt.Stringer
func (pipeline Pipeline) String() string {
	return fmt.Sprintf("#%d", pipeline.BuildNumber)
}

// MarshalJSON implements the json.Marshaler interface.
//
// implements json.Marshaler
func (pipeline Pipeline) MarshalJSON() (data []byte, err error) {
	type surrogate Pipeline

	var completedOn string
	if !pipeline.CompletedOn.IsZero() {
		completedOn = pipeline.CompletedOn.Format("2006-01-02T15:04:05.999999999-07:00")
	}

	data, err = json.Marshal(struct {
		Type string `json:"type"`
		surrogate
		CreatedOn         string `json:"created_on"`
		CompletedOn       string `json:"completed_on,omitempty"`
		DurationInSeconds uint64 `json:"duration_in_seconds"`
	}{
		Type:              pipeline.GetType(),
		surrogate:         surrogate(pipeline),
		CreatedOn:         pipeline.CreatedOn.Format("2006-01-02T15:04:05.999999999-07:00"),
		CompletedOn:       completedOn,
		DurationInSeconds: uint64(pipeline.Duration.Seconds()),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//
// implements json.Unmarshaler
func (pipeline *Pipeline) UnmarshalJSON(data []byte) error {
	type surrogate Pipeline
	var inner struct {
		Type string `json:"type"`
		surrogate
		CreatedOn         core.Time `json:"created_on"`
		CompletedOn       core.Time `json:"completed_on,omitempty"`
		DurationInSeconds uint64    `json:"duration_in_seconds"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != pipeline.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, pipeline.GetType()))
	}
	*pipeline = Pipeline(inner.surrogate)
	pipeline.CreatedOn = time.Time(inner.CreatedOn)
	pipeline.CompletedOn = time.Time(inner.CompletedOn)
	pipeline.Duration = time.Duration(inner.DurationInSeconds) * time.Second
	return errors.JSONUnmarshalError.Wrap(pipeline.Validate())
}
