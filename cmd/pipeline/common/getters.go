package plcommon

import (
	"context"
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PipelineID struct {
	ID int `json:"build_number" mapstructure:"build_number"`
}

// GetPipelineIDs gets the IDs of the pipelines
func GetPipelineIDs(context context.Context, cmd *cobra.Command, args []string, toComplete string) (ids []string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getpipelines")

	log.Infof("Getting pipelines")
	pipelines, err := profile.GetAll[PipelineID](
		log.ToContext(context),
		cmd,
		"pipelines",
	)
	if err != nil {
		log.Errorf("Failed to get pipelines", err)
		return []string{}, err
	}

	ids = core.Map(pipelines, func(pipeline PipelineID) string { return fmt.Sprintf("%d", pipeline.ID) })
	core.Sort(ids, func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == -1 })
	return ids, nil
}
