package issue

import (
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all issues",
	Args:  cobra.NoArgs,
	RunE:  listProcess,
}

var listOptions struct {
	Repository string
	States     common.EnumSliceFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.States = common.EnumSliceFlag{Allowed: []string{"closed", "duplicate", "invalid", "on hold", "new", "open", "resolved", "submitted", "wontfix"}, AllAllowed: true, Default: []string{"new", "open"}}
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issues from. Defaults to the current repository")
	listCmd.Flags().Var(&listOptions.States, "state", "State of the issues to list. Can be repeated. One of: all, closed, duplicate, invalid, on hold, new, open, resolved, submitted, wontfix. Default: open, new")
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	filter := ""
	if !listOptions.States.Contains("all") {
		if states := listOptions.States.Get(); len(states) > 0 {
			filter = "?q="
			for index, state := range states {
				if index > 0 {
					filter += "+OR+"
				}
				filter += fmt.Sprintf(`state="%s"`, strings.ReplaceAll(state, " ", "+"))
			}
		}
	}

	log.Infof("Listing all issues from repository %s with profile %s", listOptions.Repository, profile.Current)
	issues, err := profile.GetAll[Issue](cmd.Context(), cmd, profile.Current, "issues"+filter)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		log.Infof("No issue found")
		return nil
	}
	return profile.Current.Print(cmd.Context(), cmd, Issues(issues))
}
