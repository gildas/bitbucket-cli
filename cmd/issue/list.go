package issue

import (
	"fmt"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-core"
	"github.com/gildas/go-flags"
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
	States     *flags.EnumSliceFlag
}

func init() {
	Command.AddCommand(listCmd)

	listOptions.States = flags.NewEnumSliceFlagWithAllAllowed("closed", "duplicate", "invalid", "on hold", "+new", "+open", "resolved", "submitted", "wontfix")
	listCmd.Flags().StringVar(&listOptions.Repository, "repository", "", "Repository to list issues from. Defaults to the current repository")
	listCmd.Flags().Var(listOptions.States, "state", "State of the issues to list. Can be repeated. One of: all, closed, duplicate, invalid, on hold, new, open, resolved, submitted, wontfix. Default: open, new")
	_ = listCmd.RegisterFlagCompletionFunc(listOptions.States.CompletionFunc("state"))
}

func listProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "list")

	filter := ""
	if !core.Contains(listOptions.States.GetSlice(), "all") {
		if states := listOptions.States.GetSlice(); len(states) > 0 {
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
	issues, err := profile.GetAll[Issue](cmd.Context(), cmd, "issues"+filter)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		log.Infof("No issue found")
		return nil
	}
	core.Sort(issues, func(a, b Issue) bool { return a.ID < b.ID })
	return profile.Current.Print(cmd.Context(), cmd, Issues(issues))
}
