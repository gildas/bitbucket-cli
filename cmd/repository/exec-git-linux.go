//go:build linux
// +build linux

package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gildas/go-logger"
)

func GitClone(context context.Context, workspace, repository, destination, username string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child("repository", "clone")

	repoURL := url.URL{
		Scheme: "https",
		Host:   "bitbucket.org",
		Path:   fmt.Sprintf("/%s/%s.git", workspace, repository),
		User:   url.User(username),
	}
	user, err := user.Current()
	if err != nil {
		return err
	}
	out, err := exec.Command("getent", "passwd", user.Username).Output()
	if err != nil {
		return err
	}
	shell := strings.TrimSpace(strings.Split(string(out), ": ")[6])
	cmd := exec.Command(shell, "-c", fmt.Sprintf("git clone %s %s", repoURL.String(), cloneOptions.Destination))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Executing command: %s", cmd.String())
	return cmd.Run()
}
