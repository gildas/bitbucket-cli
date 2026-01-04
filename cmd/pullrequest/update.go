package pullrequest

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/project/reviewer"
	"bitbucket.org/gildas_cherruel/bb/cmd/pullrequest/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/user"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:               "update [flags] <pullrequest-id>",
	Aliases:           []string{"edit"},
	Short:             "update a pullrequest by its <pullrequest-id>.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
}

var updateOptions struct {
	Title             string
	Description       string
	Destination       *flags.EnumFlag
	AddReviewers      *flags.EnumSliceFlag
	RemoveReviewers   *flags.EnumSliceFlag
	CloseSourceBranch bool
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.Destination = flags.NewEnumFlagWithFunc("", GetBranchNames)
	updateOptions.AddReviewers = flags.NewEnumSliceFlagWithAllAllowedAndFunc(GetReviewerNicknames)
	updateOptions.RemoveReviewers = flags.NewEnumSliceFlagWithAllAllowedAndFunc(GetReviewerNicknames)

	updateCmd.Flags().StringVar(&updateOptions.Title, "title", "", "Title of the pullrequest")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the pullrequest")
	updateCmd.Flags().Var(updateOptions.Destination, "destination", "Destination branch of the pullrequest")
	updateCmd.Flags().Var(updateOptions.AddReviewers, "add-reviewer", "Reviewer(s) to add to the pullrequest. Can be specified multiple times, or as a comma-separated list. Can be the user Account ID, UUID, name, or nickname.")
	updateCmd.Flags().Var(updateOptions.RemoveReviewers, "remove-reviewer", "Reviewer(s) to remove from the pullrequest. Can be specified multiple times, or as a comma-separated list. Can be the user Account ID, UUID, name, or nickname.")
	updateCmd.Flags().BoolVar(&updateOptions.CloseSourceBranch, "close-source-branch", false, "Close the source branch after merging")

	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.Destination.CompletionFunc("destination"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.AddReviewers.CompletionFunc("add-reviewer"))
	_ = updateCmd.RegisterFlagCompletionFunc(updateOptions.RemoveReviewers.CompletionFunc("remove-reviewer"))
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	ids, err := prcommon.GetPullRequestIDsWithState(cmd.Context(), cmd, "ALL")
	if err != nil {
		cobra.CompErrorln(err.Error())
		return []string{}, cobra.ShellCompDirectiveError
	}
	return common.FilterValidArgs(ids, args, toComplete), cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	var pullrequest PullRequest

	log.Infof("Fetching pullrequest %s", args[0])

	err := profile.Current.Get(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s", args[0]),
		&pullrequest,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}
	log = log.Record("pullrequest", pullrequest.ID)
	log.Record("pullrequest", pullrequest).Debugf("Fetched pullrequest %s", args[0])

	updateWanted := false

	if cmd.Flag("title").Changed {
		pullrequest.Title = updateOptions.Title
		updateWanted = true
	}

	if cmd.Flag("description").Changed {
		pullrequest.Description = updateOptions.Description
		updateWanted = true
	}

	if cmd.Flag("destination").Changed {
		pullrequest.Destination = Endpoint{Branch: Branch{Name: updateOptions.Destination.Value}}
		updateWanted = true
	}

	if cmd.Flag("close-source-branch").Changed {
		pullrequest.CloseSourceBranch = updateOptions.CloseSourceBranch
		updateWanted = true
	}

	pullrequestWorkspace, err := pullrequest.Destination.Repository.FetchWorkspace(cmd.Context(), cmd, profile.Current)
	if err != nil {
		log.Errorf("Failed to get workspace of pullrequest destination repository", err)
		fmt.Fprintf(os.Stderr, "Failed to get workspace of pullrequest destination repository: %s\n", err)
		os.Exit(1)
	}

	isMember := func(member workspace.Member, id string) bool {
		if id, err := common.ParseUUID(id); err == nil {
			return member.User.ID == id
		}
		return member.User.AccountID == id || strings.EqualFold(member.User.Nickname, id) || strings.EqualFold(member.User.Name, id)
	}

	if cmd.Flag("remove-reviewer").Changed {
		if len(updateOptions.RemoveReviewers.Values) > 0 {
			for _, reviewerNameOrID := range updateOptions.RemoveReviewers.Values {
				var found = -1
				for index, reviewer := range pullrequest.Reviewers {
					if isMember(workspace.Member{User: reviewer}, reviewerNameOrID) {
						found = index
						break
					}
				}
				if found != -1 {
					pullrequest.Reviewers = append(pullrequest.Reviewers[:found], pullrequest.Reviewers[found+1:]...)
					updateWanted = true
				}
			}
		}
	}

	if cmd.Flag("add-reviewer").Changed {
		if len(updateOptions.AddReviewers.Values) > 0 {

			if updateOptions.AddReviewers.Values[0] == "default" {
				// Find the default reviewers from the repo or project settings
				var reviewers []reviewer.Reviewer

				log.Debugf("No reviewers in the repository, trying to get default reviewers from project settings")
				reviewers, err = reviewer.GetProjectDefaultReviewers(cmd.Context(), cmd, pullrequestWorkspace.Slug, pullrequest.Source.Repository.Project.Key)
				if err != nil {
					log.Errorf("Failed to get default reviewers", err)
					return err
				}
				log.Debugf("Found %d default reviewers", len(reviewers))
				// Replace the first reviewer with the list of default reviewers and appends the rest
				updateOptions.AddReviewers.Values = append(
					core.Map(reviewers, func(reviewer reviewer.Reviewer) string { return reviewer.User.ID.String() }),
					updateOptions.AddReviewers.Values[1:]...,
				)
			}

			log.Debugf("Getting all members from workspace %s", pullrequestWorkspace)
			members, _ := pullrequestWorkspace.GetMembers(cmd.Context(), cmd)
			log.Infof("Found %d members in workspace %s", len(members), pullrequestWorkspace)
			for _, reviewer := range updateOptions.AddReviewers.Values {
				log.Debugf("Processing reviewer to add: %s", reviewer)
				if matches := core.Filter(members, func(member workspace.Member) bool { return isMember(member, reviewer) }); len(matches) > 0 {
					if !slices.ContainsFunc(pullrequest.Reviewers, func(user user.User) bool { return user.ID == matches[0].User.ID }) {
						log.Record("matches", matches).Infof("Adding reviewer: %s (%s)", matches[0].User.ID, matches[0].User.Nickname)
						pullrequest.Reviewers = append(pullrequest.Reviewers, matches[0].User)
						updateWanted = true
					} else {
						log.Infof("Reviewer %s (%s) is already a reviewer, skipping", matches[0].User.ID, matches[0].User.Nickname)
					}
				} else {
					log.Errorf("reviewer ID %s is not a member of workspace %s", reviewer, pullrequestWorkspace)
					fmt.Fprintf(os.Stderr, "Reviewer %s is not a member of workspace %s\n", reviewer, pullrequestWorkspace)
				}
			}
		}
	}

	if !updateWanted {
		log.Infof("No update options were changed, exiting")
		return nil
	}

	log.Record("update", pullrequest).Infof("Updating pullrequest %s", args[0])
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Updating pullrequest %d", pullrequest.ID) {
		return nil
	}

	var updated PullRequest

	err = profile.Current.Put(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("pullrequests/%s", args[0]),
		pullrequest,
		&updated,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update pullrequest %s: %s\n", args[0], err)
		os.Exit(1)
	}

	return profile.Current.Print(cmd.Context(), cmd, updated)
}
