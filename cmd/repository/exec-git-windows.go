//go:build windows
// +build windows

package repository

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

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
	var cmd *exec.Cmd

	shell := os.Getenv("COMSPEC")
	if len(shell) == 0 {
		shell = "cmd.exe"
	}
	log.Debugf("Using shell: %s", shell)
	switch filepath.Base(shell) {
	case "powershell.exe":
		cmd = exec.Command(shell, "Start-Process", "git", "-ArgumentList", fmt.Sprintf(`"clone", "%s", "%s"`, repoURL.String(), destination), "-NoNewWindow")
	case "cmd.exe":
		fallthrough // Use cmd.exe for cloning
	default:
		cmd = exec.Command(shell, "/C", fmt.Sprintf("git clone %s %s", repoURL.String(), destination))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("Executing command: %s", cmd.String())
	return cmd.Run()
}
