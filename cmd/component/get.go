package component

import (
	"fmt"
	"os"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:               "get [flags] <component-id>",
	Aliases:           []string{"show", "info", "display"},
	Short:             "get a component by its <component-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: getValidArgs,
	RunE:              getProcess,
}

var getOptions struct {
	Repository string
}

func init() {
	Command.AddCommand(getCmd)

	getCmd.Flags().StringVar(&getOptions.Repository, "repository", "", "Repository to get a component from. Defaults to the current repository")
}

func getValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return GetComponentIDs(cmd.Context(), cmd), cobra.ShellCompDirectiveNoFileComp
}

func getProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "get")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	log.Infof("Displaying component %s", args[0])
	var component Component

	err = profile.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("components/%s", args[0]),
		&component,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get component %s: %s\n", args[0], err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, component)
}
