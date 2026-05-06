package common

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"gopkg.in/ini.v1"
)

// OpenGitConfig opens the .git/config file in the current folder or one of its parents
func OpenGitConfig(ctx context.Context) (io.ReadCloser, error) {
	log := logger.Must(logger.FromContext(ctx)).Child("remote", "opengitconfig")
	folder, err := filepath.Abs(".")
	if err != nil {
		folder = "."
	}
	last := folder + "dummy"

	for {
		// If .git is a filem (e.g. worktree), read the actual git dir from there (field gitdir)
		gitPath := filepath.Join(folder, ".git")
		info, err := os.Stat(gitPath)
		if err == nil && !info.IsDir() {
			log.Debugf(".git is a file, reading gitdir from there")
			content, err := os.ReadFile(gitPath)
			if err == nil {
				lines := string(content)
				const prefix = "gitdir: "
				for line := range strings.SplitSeq(lines, "\n") {
					if len(line) > len(prefix) && line[:len(prefix)] == prefix {
						gitDir := line[len(prefix):]
						log.Debugf("found gitdir: %s", gitDir)
						if !filepath.IsAbs(gitDir) {
							folder = filepath.Join(folder, gitDir)
						} else {
							folder = gitDir
						}
						break
					}
				}
			}
		}
		filename := filepath.Join(folder, ".git/config")
		if folder == last {
			return nil, errors.NotFound.With("file", filename)
		}
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
		last = folder
		folder = filepath.Dir(folder)
	}
}

// GetGitSection returns the INI section from the git config file
func GetGitSection(ctx context.Context, reader io.Reader, name string) (section *ini.Section, err error) {
	data, err := getINIContent(ctx, reader)
	if err != nil {
		return nil, err
	}
	section = data.Section(name)
	if section == nil {
		return nil, errors.NotFound.With("section", name)
	}
	return section, nil
}

// GetGitSectionsMatching returns the INI sections from the git config file matching the given regex
func GetGitSectionsMatching(ctx context.Context, reader io.Reader, rex *regexp.Regexp) (sections []*ini.Section, err error) {
	data, err := getINIContent(ctx, reader)
	if err != nil {
		return nil, err
	}
	for _, section := range data.Sections() {
		if rex.MatchString(section.Name()) {
			sections = append(sections, section)
		}
	}
	return
}

// getINIContent reads the INI content from a reader and returns it as an ini.File
func getINIContent(_ context.Context, reader io.Reader) (data *ini.File, err error) {
	payload, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	data, err = ini.Load(payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}
