package remote

import (
	"context"
	"io"
	"strings"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
)

type Remote struct {
	URL   string
	Fetch string
}

func GetFromGitConfig(context context.Context, name string) (remote *Remote, err error) {
	file, err := common.OpenGitConfig(context)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Get(context, file, name)
}

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

func (remote Remote) RepositoryName() string {
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
	}
	return remote.URL
}
