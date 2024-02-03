package common

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"gopkg.in/ini.v1"
)

// OpenGitConfig opens the .git/config file in the current folder or one of its parents
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

// GetGitSection returns the INI section from the git config file
func GetGitSection(context context.Context, reader io.Reader, name string) (section *ini.Section, err error) {
	payload, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	data, err := ini.Load(payload)
	if err != nil {
		return nil, err
	}
	section = data.Section(name)
	if section == nil {
		return nil, errors.NotFound.With("section", name)
	}
	return section, nil
}
