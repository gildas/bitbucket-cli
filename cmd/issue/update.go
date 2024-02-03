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

type IssueUpdator struct {
	Title    string               `json:"title,omitempty"`
	Kind     string               `json:"kind,omitempty"`     // bug, enhancement, proposal, task
	Priority string               `json:"priority,omitempty"` // trivial, minor, major, critical, blocker
	Content  *common.RenderedText `json:"content,omitempty"`
	Assignee *user.Account        `json:"assignee,omitempty"`
	Version  *common.Entity       `json:"version,omitempty"`
}

var updateCmd = &cobra.Command{
	Use:               "update",
	Aliases:           []string{"edit"},
	Short:             "update an issue",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
}

var updateOptions struct {
	Repository  string
	Title       string
	Kind        *flags.EnumFlag
	Priority    *flags.EnumFlag
	Description string
	Assignee    string
	Version     string
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.Kind = flags.NewEnumFlag("bug", "enhancement", "proposal", "task")
	updateOptions.Priority = flags.NewEnumFlag("major", "trivial", "minor", "major", "critical", "blocker")
	updateCmd.Flags().StringVar(&updateOptions.Repository, "repository", "", "Repository to update an issue from. Defaults to the current repository")
	updateCmd.Flags().StringVar(&updateOptions.Title, "title", "", "Title of the issue")
	updateCmd.Flags().Var(updateOptions.Kind, "kind", "Kind of the issue")
	updateCmd.Flags().Var(updateOptions.Priority, "priority", "Priority of the issue")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the issue")
	updateCmd.Flags().StringVar(&updateOptions.Assignee, "assignee", "", "Assignee of the issue. (Optional, \"myself\" or userid)")
	updateCmd.Flags().StringVar(&updateOptions.Version, "version", "", "Version of the issue")
	_ = updateCmd.RegisterFlagCompletionFunc("kind", updateOptions.Kind.CompletionFunc("kind"))
	_ = updateCmd.RegisterFlagCompletionFunc("priority", updateOptions.Priority.CompletionFunc("priority"))
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetIssueIDs(cmd.Context(), cmd, profile.Current), cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	payload := IssueUpdator{
		Title:    updateOptions.Title,
		Kind:     updateOptions.Kind.Value,
		Priority: updateOptions.Priority.Value,
	}

	if len(updateOptions.Description) > 0 {
		payload.Content = &common.RenderedText{Raw: updateOptions.Description, Markup: "markdown"}
	}

	if strings.ToLower(updateOptions.Assignee) == "me" || strings.ToLower(updateOptions.Assignee) == "myself" {
		me, err := user.GetMe(cmd.Context(), cmd, profile.Current)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current user: %s\n", err)
			os.Exit(1)
		}
		payload.Assignee = me
	} else if len(updateOptions.Assignee) > 0 {
		uuid, err := common.ParseUUID(updateOptions.Assignee)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("assignee", updateOptions.Assignee), err)
		}
		payload.Assignee = &user.Account{ID: uuid}
	}

	if len(updateOptions.Version) > 0 {
		payload.Version = &common.Entity{Name: updateOptions.Version}
	}

	log.Record("payload", payload).Infof("Updating issue %s", args[0])
	var issue Issue

	if profile.Current.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating issue %s", args[0]) {
		err = profile.Current.Put(
			log.ToContext(cmd.Context()),
			cmd,
			fmt.Sprintf("issues/%s", args[0]),
			payload,
			&issue,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update issue %s: %s\n", args[0], err)
			os.Exit(1)
		}
	}
	return profile.Current.Print(cmd.Context(), cmd, issue)
}
