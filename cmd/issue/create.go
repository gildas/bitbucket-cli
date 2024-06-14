package issue

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type IssueCreator struct {
	Title    string               `json:"title"`
	Kind     string               `json:"kind"`     // bug, enhancement, proposal, task
	Priority string               `json:"priority"` // trivial, minor, major, critical, blocker
	Content  *common.RenderedText `json:"content"`
	Assignee *user.Account        `json:"assignee,omitempty"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create an issue",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	Repository  string
	Title       string
	Kind        *flags.EnumFlag
	Priority    *flags.EnumFlag
	Description string
	Assignee    string
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.Kind = flags.NewEnumFlag("+bug", "enhancement", "proposal", "task")
	createOptions.Priority = flags.NewEnumFlag("+major", "trivial", "minor", "major", "critical", "blocker")
	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create an issue into. Defaults to the current repository")
	createCmd.Flags().StringVar(&createOptions.Title, "title", "", "Title of the issue")
	createCmd.Flags().Var(createOptions.Kind, "kind", "Kind of the issue")
	createCmd.Flags().Var(createOptions.Priority, "priority", "Priority of the issue")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the issue")
	createCmd.Flags().StringVar(&createOptions.Assignee, "assignee", "", "Assignee of the issue. (Optional, \"myself\" or userid)")
	_ = createCmd.MarkFlagRequired("title")
	_ = createCmd.MarkFlagRequired("kind")
	_ = createCmd.MarkFlagRequired("priority")
	_ = createCmd.RegisterFlagCompletionFunc("kind", createOptions.Kind.CompletionFunc("kind"))
	_ = createCmd.RegisterFlagCompletionFunc("priority", createOptions.Priority.CompletionFunc("priority"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	payload := IssueCreator{
		Title:    createOptions.Title,
		Kind:     createOptions.Kind.Value,
		Priority: createOptions.Priority.Value,
	}

	if len(createOptions.Description) > 0 {
		payload.Content = &common.RenderedText{Raw: createOptions.Description, Markup: "markdown"}
	}

	if strings.ToLower(createOptions.Assignee) == "me" || strings.ToLower(createOptions.Assignee) == "myself" {
		me, err := user.GetMe(cmd.Context(), cmd, profile.Current)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current user: %s\n", err)
			os.Exit(1)
		}
		payload.Assignee = me
	} else if len(createOptions.Assignee) > 0 {
		uuid, err := common.ParseUUID(createOptions.Assignee)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("assignee", createOptions.Assignee), err)
		}
		payload.Assignee = &user.Account{ID: uuid}
	}

	log.Record("payload", payload).Infof("Creating issue")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating issue") {
		return nil
	}
	var issue Issue

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		"issues",
		payload,
		&issue,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create project: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, issue)
}
