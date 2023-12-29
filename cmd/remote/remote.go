package remote

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"gopkg.in/ini.v1"
)

type Remote struct {
	URL   string
	Fetch string
}

func OpenGitConfig(context context.Context) (io.ReadCloser, error) {
	log := logger.Must(logger.FromContext(context)).Child("remote", "opengitconfig")
	folder := "."

	for {
		filename := filepath.Join(folder, ".git/config")
		log.Debugf("opening %s", filename)
		file, err := os.Open(filename)
		if err == nil {
			return file, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, errors.RuntimeError.Wrap(err)
		}
		if folder == "/" {
			return nil, errors.New("not a git repository")
		}
		folder += "/.."
	}
}

func GetFromGitConfig(context context.Context, name string) (remote *Remote, err error) {
	file, err := OpenGitConfig(context)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Get(context, file, name)
}

func Get(context context.Context, reader io.Reader, name string) (remote *Remote, err error) {
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
