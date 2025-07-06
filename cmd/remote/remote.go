package remote

import (
	"context"
	"io"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

// Remote represents a remote repository in a git configuration
type Remote struct {
	URL   string
	Fetch string
}

// GetFromGitConfig gets a remote from the git configuration
func GetFromGitConfig(context context.Context, name string) (remote *Remote, err error) {
	file, err := common.OpenGitConfig(context)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Get(context, file, name)
}

// Get gets a remote from a reader
func Get(context context.Context, reader io.Reader, name string) (remote *Remote, err error) {
	section, err := common.GetGitSection(context, reader, "remote \""+name+"\"")
	if err != nil {
		return nil, err
	}
	return &Remote{
		URL:   section.Key("url").String(),
		Fetch: section.Key("fetch").String(),
	}, nil
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
