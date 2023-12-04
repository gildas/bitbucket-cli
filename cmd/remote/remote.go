package remote

import (
	"io"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Remote struct {
	URL   string
	Fetch string
}

func OpenGitConfig() (io.ReadCloser, error) {
	return os.Open(".git/config")
}

func GetFromGitConfig(name string) (remote *Remote, err error) {
	file, err := OpenGitConfig()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Get(file, name)
}

func Get(reader io.Reader, name string) (remote *Remote, err error) {
	payload, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	data, err := ini.Load(payload)
	if err != nil {
		return nil, err
	}
	section := data.Section("remote \"" + name + "\"")
	return &Remote{
		URL:   section.Key("url").String(),
		Fetch: section.Key("fetch").String(),
	}, nil
}

func (remote Remote) Repository() string {
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
