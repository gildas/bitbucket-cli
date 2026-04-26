package tag

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/commit"
	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create [flags]",
	Aliases: []string{"add", "new"},
	Short:   "create a tag",
	RunE:    createProcess,
}

var createOptions struct {
	Repository string
	Name       string
	Message    string
	Commit     *flags.EnumFlag
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.Commit = flags.NewEnumFlagWithFunc("latest", commit.GetCommitHashes)
	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create a tag into. Defaults to the current repository.\nExpected format: <workspace>/<repository> or <repository>.\nIf only <repository> is given, the profile's default workspace is used.")
	createCmd.Flags().StringVar(&createOptions.Name, "name", "", "Name of the tag")
	createCmd.Flags().StringVar(&createOptions.Message, "message", "", "Message of the tag")
	createCmd.Flags().Var(createOptions.Commit, "commit", "Target commit hash for the tag. Defaults to the latest commit")
	_ = createCmd.MarkFlagRequired("name")

	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Commit.CompletionFunc("commit"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	profile, err := profile.GetProfileFromCommand(cmd.Context(), cmd)
	if err != nil {
		return err
	}

	payload := Tag{
		Name:    createOptions.Name,
		Message: createOptions.Message,
	}

	if len(createOptions.Commit.Value) > 0 && strings.ToLower(createOptions.Commit.Value) != "latest" && strings.ToLower(createOptions.Commit.Value) != "head" {
		payload.Commit.Hash = commit.Commit{Hash: createOptions.Commit.Value}.GetShortHash()
	} else {
		commit, err := commit.GetLatestCommit(log.ToContext(cmd.Context()), cmd)
		if err != nil {
			return err
		}
		payload.Commit.Hash = commit.GetShortHash()
	}

	log.Record("payload", payload).Infof("Creating tag %s", createOptions.Name)
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating tag %s", createOptions.Name) {
		return nil
	}

	var tag Tag

	err = profile.Post(
		log.ToContext(cmd.Context()),
		cmd,
		"refs/tags",
		payload,
		&tag,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create tag: %s\n", err)
		os.Exit(1)
	}
	return profile.Print(cmd.Context(), cmd, tag)
}
