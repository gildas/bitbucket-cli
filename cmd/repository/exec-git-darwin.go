//go:build darwin
// +build darwin

package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
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
	out, err := exec.Command("dscl", "localhost", "-read", path.Join("Local", "Default", "Users", os.Getenv("USER")), "UserShell").Output()
	if err != nil {
		return err
	}
	shell := strings.TrimSpace(strings.Split(string(out), ": ")[1])
	cmd := exec.Command(shell, "-c", fmt.Sprintf("git clone %s %s", repoURL.String(), cloneOptions.Destination))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Executing command: %s", cmd.String())
	return cmd.Run()
}
