package remote

import (
	"context"
	"io"
	"regexp"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/gildas/go-errors"
	"github.com/spf13/cobra"
)

// Remote represents a remote repository in a git configuration
type Remote struct {
	URL   string
	Fetch string
}

// GetRemoteFromGitConfig gets a remote from the git configuration
func GetRemoteFromGitConfig(context context.Context, name string) (remote *Remote, err error) {
	file, err := common.OpenGitConfig(context)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return GetRemoteFromReader(context, file, name)
}

// GetRemoteFromReader gets a remote from a reader
//
// If the name is empty, it gets the first remote in the reader
func GetRemoteFromReader(context context.Context, reader io.Reader, name string) (remote *Remote, err error) {
	if len(name) == 0 {
		sections, err := common.GetGitSectionsMatching(context, reader, regexp.MustCompile("remote \".*\""))
		if err != nil {
			return nil, err
		}
		if len(sections) == 0 {
			return nil, errors.NotFound.With("remote", "any")
		}
		section := sections[0]
		return &Remote{
			URL:   section.Key("url").String(),
			Fetch: section.Key("fetch").String(),
		}, nil
	}
	section, err := common.GetGitSection(context, reader, "remote \""+name+"\"")
	if err != nil {
		return nil, err
	}
	return &Remote{
		URL:   section.Key("url").String(),
		Fetch: section.Key("fetch").String(),
	}, nil
}

// GetRemote gets the remote from the command flags or the git configuration
//
// Checks the --git-remote flag first,
//
// Then checks the "origin" remote in the git configuration,
//
// Falls back to the first remote in the git configuration
func GetRemote(context context.Context, cmd *cobra.Command) (remote *Remote, err error) {
	if cmd.Flag("git-remote") != nil {
		remoteName := cmd.Flag("git-remote").Value.String()
		if len(remoteName) > 0 {
			return GetRemoteFromGitConfig(context, remoteName)
		}
	}
	if remote, err = GetRemoteFromGitConfig(context, "origin"); err == nil {
		return
	}
	return GetRemoteFromGitConfig(context, "")
}

// RepositoryName gets the full repository name from the remote URL (without the .git extension)
func (remote Remote) RepositoryName() string {
	if strings.HasPrefix(remote.URL, "bitbucket.org:") {
		if strings.HasSuffix(remote.URL, ".git") {
			return remote.URL[strings.Index(remote.URL, ":")+1 : len(remote.URL)-4]
		}
		return remote.URL[strings.Index(remote.URL, ":")+1:]
	}
	if strings.HasPrefix(remote.URL, "git@") {
		if strings.HasSuffix(remote.URL, ".git") {
			return remote.URL[strings.LastIndex(remote.URL, ":")+1 : len(remote.URL)-4]
		}
		return remote.URL[strings.LastIndex(remote.URL, ":")+1:]
	} else if strings.HasPrefix(remote.URL, "https://") {
		previousToLastSlash := strings.LastIndex(remote.URL[:strings.LastIndex(remote.URL, "/")], "/")
		if strings.HasSuffix(remote.URL, ".git") {
			return remote.URL[previousToLastSlash+1 : len(remote.URL)-4]
		}
		return remote.URL[previousToLastSlash+1:]
	} else if strings.HasPrefix(remote.URL, "ssh://") {
		previousToLastSlash := strings.LastIndex(remote.URL[:strings.LastIndex(remote.URL, "/")], "/")
		if strings.HasSuffix(remote.URL, ".git") {
			return remote.URL[previousToLastSlash+1 : len(remote.URL)-4]
		}
		return remote.URL[previousToLastSlash+1:]
	}
	return remote.URL
}

// WorkspaceName gets the workspace name from the remote URL
func (remote Remote) WorkspaceName() string {
	repositoryName := remote.RepositoryName()
	return repositoryName[:strings.Index(repositoryName, "/")]
}
