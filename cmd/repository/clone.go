package repository

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:               "clone [flags] <slug>",
	Short:             "clone a repository by its <slug>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: cloneValidArgs,
	RunE:              cloneProcess,
}

var cloneOptions struct {
	Workspace   *flags.EnumFlag
	Destination string
	Bare        bool
	Username    string
}

func init() {
	Command.AddCommand(cloneCmd)

	cloneOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	cloneCmd.Flags().Var(cloneOptions.Workspace, "workspace", "Workspace to clone repositories from. If omitted, it will be extracted from the repository name")
	cloneCmd.Flags().StringVar(&cloneOptions.Destination, "destination", "", "Destination folder. Default is the repository name")
	cloneCmd.Flags().BoolVar(&cloneOptions.Bare, "bare", false, "Clone as a bare repository")
	cloneCmd.Flags().StringVar(&cloneOptions.Username, "username", "", "Username to use for authentication. If not set, your git and ssh configuration will take precedence")
	_ = cloneCmd.MarkFlagDirname("destination")
	_ = cloneCmd.RegisterFlagCompletionFunc(cloneOptions.Workspace.CompletionFunc("workspace"))
}

func cloneValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	slugs, err := GetRepositorySlugs(cmd.Context(), cmd, cloneOptions.Workspace.String())
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(slugs, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func cloneProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "clone")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(cloneOptions.Workspace.Value) == 0 {
		cloneOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(cloneOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	if len(cloneOptions.Workspace.Value) == 0 {
		components := strings.Split(args[0], "/")
		if len(components) != 2 {
			return errors.ArgumentInvalid.With("repository", args[0])
		}
		cloneOptions.Workspace.Value = components[0]
		args[0] = components[1]
	}

	if len(cloneOptions.Destination) == 0 {
		cloneOptions.Destination = strings.TrimSuffix(args[0], ".git")
		log.Debugf("Destination not specified, using repository slug as destination: %s", cloneOptions.Destination)
	}

	var user *url.Userinfo

	if len(cloneOptions.Username) > 0 {
		user = url.User(cloneOptions.Username)
	}

	repoURL := url.URL{
		Scheme: "https",
		Host:   "bitbucket.org",
		Path:   fmt.Sprintf("/%s/%s.git", cloneOptions.Workspace.String(), args[0]),
		User:   user,
	}
	log.Debugf("Cloning repository from %s", repoURL.String())
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Cloning repository %s/%s into %s", cloneOptions.Workspace, args[0], cloneOptions.Destination) {
		return nil
	}
	_, err := git.PlainCloneContext(log.ToContext(cmd.Context()), cloneOptions.Destination, cloneOptions.Bare, &git.CloneOptions{
		URL:      repoURL.String(),
		Progress: os.Stdout,
	})
	return err
}
