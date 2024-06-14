package project

import (
	"encoding/base64"
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
	"github.com/spf13/cobra"
)

type ProjectCreator struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Key         string        `json:"key"`
	Links       *common.Links `json:"links,omitempty"`
	IsPrivate   bool          `json:"is_private"`
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add", "new"},
	Short:   "create a project",
	Args:    cobra.NoArgs,
	RunE:    createProcess,
}

var createOptions struct {
	Workspace   *flags.EnumFlag
	Name        string
	Key         string
	Description string
	AvatarURL   string
	AvatarPath  string
	IsPrivate   bool
}

func init() {
	Command.AddCommand(createCmd)

	createOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	createCmd.Flags().Var(createOptions.Workspace, "workspace", "Workspace to create projects from")
	createCmd.Flags().StringVar(&createOptions.Name, "name", "", "Name of the project")
	createCmd.Flags().StringVar(&createOptions.Key, "key", "", "Key of the project")
	createCmd.Flags().StringVar(&createOptions.Description, "description", "", "Description of the project")
	createCmd.Flags().StringVar(&createOptions.AvatarURL, "avatar-url", "", "Avatar of the project")
	createCmd.Flags().StringVar(&createOptions.AvatarPath, "avatar-file", "", "Avatar of the project")
	createCmd.Flags().BoolVar(&createOptions.IsPrivate, "is-private", false, "Is the project private")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("key")
	_ = createCmd.MarkFlagFilename("avatar-file")
	createCmd.MarkFlagsMutuallyExclusive("avatar-url", "avatar-file")
	_ = createCmd.RegisterFlagCompletionFunc("workspace", createOptions.Workspace.CompletionFunc("workspace"))
}

func createProcess(cmd *cobra.Command, args []string) (err error) {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "create")

	if profile.Current == nil {
		return errors.ArgumentMissing.With("profile")
	}
	if len(createOptions.Workspace.Value) == 0 {
		createOptions.Workspace.Value = profile.Current.DefaultWorkspace
		if len(createOptions.Workspace.Value) == 0 {
			return errors.ArgumentMissing.With("workspace")
		}
	}

	payload := ProjectCreator{
		Name:        createOptions.Name,
		Key:         createOptions.Key,
		Description: createOptions.Description,
		IsPrivate:   createOptions.IsPrivate,
	}

	if len(createOptions.AvatarPath) > 0 {
		log.Debugf("Avatar is a file: %s", createOptions.AvatarPath)
		avatarData, err := os.ReadFile(createOptions.AvatarPath)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("avatar-path", createOptions.AvatarPath), err)
		}
		avatarBlob := base64.StdEncoding.EncodeToString(avatarData)
		payload.Links = &common.Links{
			Avatar: &common.Link{HREF: url.URL{Scheme: "data", Opaque: "image/png;base64," + avatarBlob}},
		}
	} else if strings.HasPrefix(createOptions.AvatarURL, "http") {
		avatarURL, err := url.Parse(createOptions.AvatarURL)
		if err != nil {
			return errors.Join(errors.ArgumentInvalid.With("avatar", createOptions.AvatarURL), err)
		}
		log.Debugf("Avatar is an URL: %s", createOptions.AvatarURL)
		payload.Links = &common.Links{
			Avatar: &common.Link{HREF: *avatarURL},
		}
	} else if len(createOptions.AvatarURL) > 0 {
		log.Errorf("Avatar is not a file nor an URL: %s", createOptions.AvatarURL)
		fmt.Fprintln(os.Stderr, "Avatar is not a file nor an URL")
		os.Exit(1)
	}

	log.Record("payload", payload).Infof("Creating project")
	if !common.WhatIf(log.ToContext(cmd.Context()), cmd, "Creating project") {
		return nil
	}
	var project Project

	err = profile.Current.Post(
		log.ToContext(cmd.Context()),
		cmd,
		fmt.Sprintf("/workspaces/%s/projects", createOptions.Workspace),
		payload,
		&project,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create project: %s\n", err)
		os.Exit(1)
	}
	return profile.Current.Print(cmd.Context(), cmd, project)
}
