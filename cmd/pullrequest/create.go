package pullrequest

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project/reviewer"
	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type PullRequestCreator struct {
	Title             string      `json:"title"`
	Description       string      `json:"description,omitempty"`
	Source            Endpoint    `json:"source"`
	Destination       *Endpoint   `json:"destination,omitempty"`
	Reviewers         []user.User `json:"reviewers,omitempty"`
	CloseSourceBranch bool        `json:"close_source_branch,omitempty"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create a pullrequest",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	Repository        string
	Title             string
	Description       string
	Source            *flags.EnumFlag
	Destination       *flags.EnumFlag
	Reviewers         *flags.EnumSliceFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.Source = flags.NewEnumFlagWithFunc("", GetBranchNames)
	createOptions.Destination = flags.NewEnumFlagWithFunc("", GetBranchNames)
	createOptions.Reviewers = flags.NewEnumSliceFlagWithAllAllowedAndFunc(GetReviewerNicknames)

	createCmd.Flags().StringVar(&createOptions.Repository, "repository", "", "Repository to create pullrequest in. Defaults to the current repository")
	createCmd.Flags().StringVar(&createOptions.Title, "title", "", "Title of the pullrequest")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the pullrequest")
	createCmd.Flags().Var(createOptions.Source, "source", "Source branch of the pullrequest")
	createCmd.Flags().Var(createOptions.Destination, "destination", "Destination branch of the pullrequest")
	createCmd.Flags().Var(createOptions.Reviewers, "reviewer", "Reviewer(s) of the pullrequest. Can be specified multiple times, or as a comma-separated list. Can be the user Account ID, UUID, name, or nickname. If the first reviewer is `default`, the command will try to find the default reviewers from the repository or project settings.")
	createCmd.Flags().BoolVar(&createOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch of the pullrequest")
	_ = createCmd.MarkFlagRequired("title")
	_ = createCmd.MarkFlagRequired("source")
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Source.CompletionFunc("source"))
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Destination.CompletionFunc("destination"))
	_ = createCmd.RegisterFlagCompletionFunc(createOptions.Reviewers.CompletionFunc("reviewer"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	if len(createOptions.Title) == 0 {
		return errors.ArgumentMissing.With("title")
	}

	payload := PullRequestCreator{
		Title:             createOptions.Title,
		Description:       createOptions.Description,
		Source:            Endpoint{Branch: Branch{Name: createOptions.Source.Value}},
		CloseSourceBranch: createOptions.CloseSourceBranch,
	}
	if len(createOptions.Destination.Value) > 0 {
		payload.Destination = &Endpoint{Branch: Branch{Name: createOptions.Destination.Value}}
	}

	pullrequestRepository, err := repository.GetRepositoryFromGit(cmd.Context(), cmd, profile.Current)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get repository: %s\n", err)
		os.Exit(1)
	}
	log.Record("repository", pullrequestRepository).Infof("Using repository: %s", pullrequestRepository.Slug)

	if len(createOptions.Reviewers.Values) > 0 {
		if createOptions.Reviewers.Values[0] == "default" {
			// Find the default reviewers from the repo or project settings
			var reviewers []reviewer.Reviewer

			log.Debugf("No reviewers in the repository, trying to get default reviewers from project settings")
			reviewers, err = reviewer.GetProjectDefaultReviewers(cmd.Context(), cmd, pullrequestRepository.Workspace.Slug, pullrequestRepository.Project.Key)
			if err != nil {
				log.Errorf("Failed to get default reviewers", err)
				return err
			}
			log.Debugf("Found %d default reviewers", len(reviewers))
			// Replace the first reviewer with the list of default reviewers and appends the rest
			createOptions.Reviewers.Values = append(
				core.Map(reviewers, func(reviewer reviewer.Reviewer) string { return reviewer.User.ID.String() }),
				createOptions.Reviewers.Values[1:]...,
			)
		}

		isMember := func(member workspace.Member, id string) bool {
			if id, err := common.ParseUUID(id); err == nil {
				return member.User.ID == id
			}
			return member.User.AccountID == id || strings.EqualFold(member.User.Nickname, id) || strings.EqualFold(member.User.Name, id)
		}

		members, _ := pullrequestRepository.Workspace.GetMembers(cmd.Context(), cmd)
		payload.Reviewers = make([]user.User, 0, len(createOptions.Reviewers.Values))
		for _, reviewer := range createOptions.Reviewers.Values {
			if matches := core.Filter(members, func(member workspace.Member) bool { return isMember(member, reviewer) }); len(matches) > 0 {
				log.Record("matches", matches).Infof("Adding reviewer: %s", matches[0].User.ID)
				payload.Reviewers = append(payload.Reviewers, matches[0].User)
			} else {
				log.Errorf("Failed to parse reviewer ID: %s", reviewer)
				fmt.Fprintf(os.Stderr, "Failed to parse reviewer ID: %s\n", reviewer)
			}
		}
	}

	log.Record("payload", payload).Infof("Creating pullrequest")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating pullrequest") {
		return nil
	}
	var pullrequest PullRequest

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		"pullrequests",
		payload,
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pullrequest: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, pullrequest)
}
