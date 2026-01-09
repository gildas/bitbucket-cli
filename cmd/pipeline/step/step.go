package step

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	plcommon "bitbucket.org/gildas_cherruel/bb/cmd/pipeline/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

// Stp represents a pipeline step
type Step struct {
	ID               common.UUID                `json:"uuid"                mapstructure:"uuid"`
	Name             string                     `json:"name"                mapstructure:"name"`
	BuildNumber      uint64                     `json:"-"                   mapstructure:"-"`
	RunNumber        uint64                     `json:"run_number"          mapstructure:"run_number"`
	Pipeline         plcommon.PipelineReference `json:"pipeline"            mapstructure:"pipeline"`
	State            StepState                  `json:"state"               mapstructure:"state"`
	Image            StepImage                  `json:"image"               mapstructure:"image"`
	SetupCommands    []StepCommand              `json:"setup_commands"      mapstructure:"setup_commands"`
	ScriptCommands   []StepCommand              `json:"script_commands"     mapstructure:"script_commands"`
	TeardownCommands []StepCommand              `json:"teardown_commands"   mapstructure:"teardown_commands"`
	MaxTime          time.Duration              `json:"maxTime"             mapstructure:"maxTime"`
	StartedOn        time.Time                  `json:"started_on"          mapstructure:"started_on"`
	CompletedOn      time.Time                  `json:"completed_on"        mapstructure:"completed_on"`
	Duration         time.Duration              `json:"duration_in_seconds" mapstructure:"duration_in_seconds"`
	ShowLogsCommand  bool                       `json:"-"                   mapstructure:"-"`
}

// Command represents this folder's command
var Command = &cobra.Command{
	Use:   "step",
	Short: "Manage pipeline steps",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Step requires a subcommand:")
		for _, command := range cmd.Commands() {
			fmt.Println(command.Name())
		}
	},
}

var columns = common.Columns[Step]{
	{Name: "id", DefaultSorter: true, Compare: func(a, b Step) bool {
		return a.ID.String() < b.ID.String()
	}},
	{Name: "state", DefaultSorter: false, Compare: func(a, b Step) bool {
		return a.State.Name < b.State.Name
	}},
	{Name: "image", DefaultSorter: false, Compare: func(a, b Step) bool {
		return a.Image.Name < b.Image.Name
	}},
	{Name: "duration", DefaultSorter: false, Compare: func(a, b Step) bool {
		return a.Duration < b.Duration
	}},
	{Name: "logs-command", DefaultSorter: false, Compare: func(a, b Step) bool {
		return a.ID.String() < b.ID.String()
	}},
}

// GetType gets the type of the struct
//
// implements core.TypeCarrier
func (step Step) GetType() string {
	return "pipeline_step"
}

// GetHeaders gets the headers for a table
//
// implements common.Tableable
func (step Step) GetHeaders(cmd *cobra.Command) []string {
	if cmd != nil && cmd.Flag("columns") != nil && cmd.Flag("columns").Changed {
		if columns, err := cmd.Flags().GetStringSlice("columns"); err == nil {
			return core.Map(columns, func(column string) string { return strings.ReplaceAll(column, "_", " ") })
		}
	}
	if step.ShowLogsCommand {
		return []string{"ID", "Name", "State", "Duration", "Image", "Logs Command"}
	}
	return []string{"ID", "Name", "State", "Duration", "Image"}
}

// GetRow gets the row for a table
//
// implements common.Tableable
func (step Step) GetRow(headers []string) []string {
	var row []string

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "id":
			row = append(row, step.ID.String())
		case "state":
			row = append(row, step.State.String())
		case "image":
			row = append(row, step.Image.Name)
		case "duration":
			row = append(row, step.Duration.String())
		case "started on":
			row = append(row, step.StartedOn.Format(time.RFC3339))
		case "completed on":
			row = append(row, step.CompletedOn.Format(time.RFC3339))
		case "name":
			row = append(row, step.Name)
		case "run number":
			row = append(row, fmt.Sprintf("%d", step.RunNumber))
		case "max time":
			row = append(row, step.MaxTime.String())
		case "logs command", "logs-command":
			row = append(row, fmt.Sprintf("bb pipeline step logs --pipeline %d --step %s", step.BuildNumber, step.ID.String()))
		default:
			row = append(row, "")
		}
	}

	return row
}

// MarshalJSON marshals the Step struct into a JSON string
//
// implements json.Marshaler
func (step Step) MarshalJSON() ([]byte, error) {
	type surrogate Step
	data, err := json.Marshal(struct {
		Type string `json:"type" mapstructure:"type"`
		surrogate
		CreatedOn         core.Time `json:"created_on"`
		CompletedOn       core.Time `json:"completed_on,omitempty"`
		MaxTime           uint64    `json:"maxTime"`
		DurationInSeconds uint64    `json:"duration_in_seconds"`
	}{
		Type:              step.GetType(),
		surrogate:         surrogate(step),
		CreatedOn:         core.Time(step.StartedOn),
		CompletedOn:       core.Time(step.CompletedOn),
		MaxTime:           uint64(step.MaxTime.Seconds()),
		DurationInSeconds: uint64(step.Duration.Seconds()),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals a JSON string into the Step struct
//
// implements json.Unmarshaler
func (step *Step) UnmarshalJSON(data []byte) error {
	type surrogate Step
	var inner struct {
		Type string `json:"type" mapstructure:"type"`
		surrogate
		CreatedOn         core.Time `json:"created_on"`
		CompletedOn       core.Time `json:"completed_on,omitempty"`
		MaxTime           uint64    `json:"maxTime"`
		DurationInSeconds uint64    `json:"duration_in_seconds"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.WrapIfNotMe(err)
	}
	if inner.Type != step.GetType() {
		return errors.JSONUnmarshalError.Wrap(errors.InvalidType.With(inner.Type, step.GetType()))
	}
	*step = Step(inner.surrogate)
	step.StartedOn = time.Time(inner.StartedOn)
	step.CompletedOn = time.Time(inner.CompletedOn)
	step.Duration = time.Duration(inner.DurationInSeconds) * time.Second
	step.MaxTime = time.Duration(inner.MaxTime) * time.Second
	return nil
}

// GetPipelineStepIDs gets the IDs of the steps for a pipeline
func GetPipelineStepIDs(context context.Context, cmd *cobra.Command, PipelineID string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child("pipeline", "getids")

	steps, err := profile.GetAll[Step](context, cmd, fmt.Sprintf("pipelines/%s/steps", PipelineID))
	if err != nil {
		log.Errorf("Failed to get pipelines", err)
		return []string{}, err
	}
	return core.Map(steps, func(step Step) string {
		return step.ID.String()
	}), nil
}
