package repository

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"bitbucket.org/gildas_cherruel/bb/cmd/workspace"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-flags"
	"github.com/gildas/go-logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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
	Workspace      *flags.EnumFlag
	Destination    string
	Bare           bool
	Protocol       *flags.EnumFlag
	SshKeyFilename string
	VaultKey       string
	Username       string
	Password       string
}

func init() {
	Command.AddCommand(cloneCmd)

	cloneOptions.Workspace = flags.NewEnumFlagWithFunc("", workspace.GetWorkspaceSlugs)
	cloneOptions.Protocol = flags.NewEnumFlag("git", "https", "ssh")
	cloneCmd.Flags().Var(cloneOptions.Workspace, "workspace", "Workspace to clone repositories from. If omitted, it will be extracted from the repository name")
	cloneCmd.Flags().StringVar(&cloneOptions.Destination, "destination", "", "Destination folder. Default is the repository name")
	cloneCmd.Flags().BoolVar(&cloneOptions.Bare, "bare", false, "Clone as a bare repository")
	cloneCmd.Flags().Var(cloneOptions.Protocol, "protocol", "Protocol to use for cloning. Default is set in the profile, can be https, git, or ssh")
	cloneCmd.Flags().StringVar(&cloneOptions.SshKeyFilename, "ssh-key-file", "", "Path to the SSH private key file. Default is ~/.ssh/id_rsa")
	cloneCmd.Flags().StringVar(&cloneOptions.VaultKey, "vault-key", "", "Vault key to use for authentication. On Windows, the Windows Credential Manager will be used, On Linux and macOS, the system keychain will be used. If not set, your git and ssh configuration will take precedence")
	cloneCmd.Flags().StringVar(&cloneOptions.Username, "username", "", "Username for authentication. If not set, it will be retrieved from the bitbucket-cli configuration")
	cloneCmd.Flags().StringVar(&cloneOptions.Password, "password", "", "Password for authentication. If not set, it will be retrieved from the bitbucket-cli configuration")
	_ = cloneCmd.MarkFlagDirname("destination")
	_ = cloneCmd.MarkFlagFilename("ssh-key-file")
	_ = cloneCmd.RegisterFlagCompletionFunc(cloneOptions.Workspace.CompletionFunc("workspace"))
	_ = cloneCmd.RegisterFlagCompletionFunc(cloneOptions.Protocol.CompletionFunc("protocol"))
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

func cloneProcess(cmd *cobra.Command, args []string) (err error) {
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

	options := git.CloneOptions{Progress: os.Stdout}

	switch cloneOptions.Protocol.Value {
	case "git":
		options.URL = fmt.Sprintf("git@bitbucket.org:%s/%s.git", cloneOptions.Workspace.String(), args[0])
	case "ssh":
		if len(cloneOptions.Username) > 0 {
			return errors.New("SSH protocol does not support username.")
		}
		options.URL = fmt.Sprintf("ssh://git@bitbucket.org/%s/%s.git", cloneOptions.Workspace.String(), args[0])
		if len(cloneOptions.SshKeyFilename) == 0 {
			if homeDir, err := os.UserHomeDir(); err != nil {
				return errors.Wrap(err, "failed to get user home directory")
			} else {
				cloneOptions.SshKeyFilename = filepath.Join(homeDir, ".ssh", "id_rsa")
				log.Debugf("SSH key file not specified, using default: %s", cloneOptions.SshKeyFilename)
			}
		}
		auth, err := ssh.NewPublicKeysFromFile("git", cloneOptions.SshKeyFilename, "")
		if err != nil {
			return errors.Wrap(err, "failed to create ssh auth")
		}
		options.Auth = auth
	default:
		if len(cloneOptions.SshKeyFilename) > 0 {
			return errors.New("SSH key file is only applicable for SSH protocol")
		}
		repoURL := url.URL{
			Scheme: "https",
			Host:   "bitbucket.org",
			Path:   fmt.Sprintf("/%s/%s.git", cloneOptions.Workspace.String(), args[0]),
		}
		options.URL = repoURL.String()
		vaultUsername := cloneOptions.Username
		if len(vaultUsername) == 0 {
			vaultUsername = profile.Current.CloneVaultUsername
			if len(vaultUsername) == 0 {
				vaultUsername = profile.Current.User
			}
		}
		if len(vaultUsername) > 0 {
			// go-git does not support username with bitbucket.org authentication, so we need to call git directly
			return GitClone(cmd.Context(), cloneOptions.Workspace.String(), args[0], cloneOptions.Destination, vaultUsername)
		}
	}

	_, err = git.PlainCloneContext(log.ToContext(cmd.Context()), cloneOptions.Destination, cloneOptions.Bare, &options)
	return err
}
