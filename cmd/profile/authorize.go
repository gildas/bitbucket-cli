package profile

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/common"
	"github.com/briandowns/spinner"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
)

var authorizeCmd = &cobra.Command{
	Use:               "authorize",
	Short:             "authorize an Authorization Code Grant profile",
	ValidArgsFunction: ValidProfileNames,
	RunE:              authorizeProcess,
}

func init() {
	Command.AddCommand(authorizeCmd)
}

func authorizeProcess(cmd *cobra.Command, args []string) error {
	log := logger.Must(logger.FromContext(cmd.Context())).Child(cmd.Parent().Name(), "authorize")

	if len(args) == 0 {
		return errors.ArgumentMissing.With("profile")
	}

	log.Infof("Authorizing profile %s (Valid names: %v)", args[0], Profiles.Names())
	profile, found := Profiles.Find(args[0])
	if !found {
		return errors.NotFound.With("profile", args[0])
	}
	if profile.CallbackPort == 0 {
		return errors.Join(errors.Errorf("Profile %s does not support Authorization Code Grant", profile.Name), errors.ArgumentInvalid.With("profile", profile.Name))
	}

	// Start a web server to listen for the Authorization Code Grant
	resultchan := make(chan error)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", profile.CallbackPort),
		Handler: log.HttpHandler()(profile.CodeGrantCallback(resultchan)),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Open the browser to the Authorization Code Grant URL
	common.Verbose(cmd.Context(), cmd, "Opening browser to authorize profile %s...", profile.Name)
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	bitbucketAuthURL := url.URL{
		Scheme: "https",
		Host:   "bitbucket.org",
		Path:   "/site/oauth2/authorize",
		RawQuery: url.Values{
			"response_type": {"code"},
			"client_id":     {profile.ClientID},
		}.Encode(),
	}
	common.Verbose(cmd.Context(), cmd, "\nIf you are not redirected automatically, please open the following URL in your browser:\n%s\n", bitbucketAuthURL.String())

	if cmd.Flag("verbose").Changed {
		spinner.Reverse()
		_ = spinner.Color("blue", "bold")
		spinner.Start()
	}

	err := openBrowser(bitbucketAuthURL)
	if err != nil {
		log.Errorf("Failed to open browser: %v", err)
		spinner.Stop()
		return err
	}

	// Wait until the user stops the server by pressing Ctrl+C
	results := <-resultchan

	spinner.Stop()
	log.Infof("Received results, shutting down server...")
	if err := server.Shutdown(cmd.Context()); err != nil {
		log.Errorf("Failed to shut down server: %v", err)
	}

	if results != nil {
		log.Errorf("Authorization process failed: %v", results)
		return results
	}
	common.Verbose(cmd.Context(), cmd, "Authorization process completed successfully")
	return nil
}

// openBrowser opens the specified URL in the default web browser
func openBrowser(url url.URL) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	case "darwin":
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url.String())
	return exec.Command(cmd, args...).Start()
}
