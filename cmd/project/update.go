package project

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/link"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

type ProjectUpdator struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Key         string      `json:"key,omitempty"`
	Links       *link.Links `json:"links,omitempty"`
	IsPrivate   bool        `json:"is_private"`
}

var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "update a project",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: updateValidArgs,
	RunE:              updateProcess,
}

var updateOptions struct {
	Workspace   common.RemoteValueFlag
	Name        string
	Key         string
	Description string
	AvatarURL   string
	AvatarPath  string
	IsPrivate   bool
}

func init() {
	Command.AddCommand(updateCmd)

	updateOptions.Workspace = common.RemoteValueFlag{AllowedFunc: workspace.GetWorkspaceSlugs}
	updateCmd.Flags().Var(&updateOptions.Workspace, "workspace", "Workspace to update projects from")
	updateCmd.Flags().StringVar(&updateOptions.Name, "name", "", "Name of the project")
	updateCmd.Flags().StringVar(&updateOptions.Key, "key", "", "Key of the project")
	updateCmd.Flags().StringVar(&updateOptions.Description, "description", "", "Description of the project")
	updateCmd.Flags().StringVar(&updateOptions.AvatarURL, "avatar-url", "", "Avatar of the project")
	updateCmd.Flags().StringVar(&updateOptions.AvatarPath, "avatar-file", "", "Avatar of the project")
	updateCmd.Flags().BoolVar(&updateOptions.IsPrivate, "is-private", false, "Is the project private")
	_ = updateCmd.MarkFlagRequired("workspace")
	updateCmd.MarkFlagsMutuallyExclusive("avatar-url", "avatar-file")
	_ = updateCmd.RegisterFlagCompletionFunc("workspace", updateOptions.Workspace.CompletionFunc())
}

func updateValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if profile.Current == nil {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}
	return GetProjectKeys(cmd.Context(), profile.Current, deleteOptions.Workspace.String()), cobra.ShellCompDirectiveNoFileComp
}

func updateProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "update")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}

	payload := ProjectUpdator{
		Name:        updateOptions.Name,
		Key:         updateOptions.Key,
		Description: updateOptions.Description,
		IsPrivate:   updateOptions.IsPrivate,
	}

	if len(updateOptions.AvatarPath) > 0 {
		log.Debugf("Avatar is a file: %s", updateOptions.AvatarPath)
		avatarData, err := os.ReadFile(updateOptions.AvatarPath)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("avatar-path", updateOptions.AvatarPath), err)
		}
		avatarBlob := base64.StdEncoding.EncodeToString(avatarData)
		payload.Links = &link.Links{
			Avatar: &link.Link{HREF: url.URL{Scheme: "data", Opaque: "image/png;base64," + avatarBlob}},
		}
	} else if strings.HasPrefix(updateOptions.AvatarURL, "http") {
		avatarURL, err := url.Parse(updateOptions.AvatarURL)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("avatar", updateOptions.AvatarURL), err)
		}
		log.Debugf("Avatar is an URL: %s", updateOptions.AvatarURL)
		payload.Links = &link.Links{
			Avatar: &link.Link{HREF: *avatarURL},
		}
	} else if len(updateOptions.AvatarURL) > 0 {
		log.Errorf("Avatar is not a file nor an URL: %s", updateOptions.AvatarURL)
		fmt.Fprintln(os.Stderr, "Avatar is not a file nor an URL")
		os.Exit(1)
	}

	log.Record("payload", payload).Infof("Updating project %s", args[0])
	var project Project

	err := profile.Current.Put(
		log.ToContext(cmd.Context()),
		"",
		fmt.Sprintf("/workspaces/%s/projects/%s", updateOptions.Workspace, args[0]),
		payload,
		&project,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update project: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), project)
}
