package profile

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type PaginatedResources[T any] struct {
	Values   []T    `json:"values"`
	Page     int    `json:"page"`
	PageSize int    `json:"pagelen"`
	Size     int    `json:"size"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

// Post posts a resource
func (profile *Profile) Post(context context.Context, cmd *cobra.Command, uripath string, body interface{}, response interface{}) (err error) {
	options := &request.Options{Method: http.MethodPost, Payload: body}
	_, err = profile.send(context, cmd, options, uripath, response)
	return
}

// Get gets a resource
func (profile *Profile) Get(context context.Context, cmd *cobra.Command, uripath string, response interface{}) (err error) {
	options := &request.Options{Method: http.MethodGet}
	_, err = profile.send(context, cmd, options, uripath, response)
	return
}

// GetRaw gets a resource without unmarshaling it
func (profile *Profile) GetRaw(context context.Context, cmd *cobra.Command, uripath string) (raw io.Reader, err error) {
	options := &request.Options{
		Method: http.MethodGet,
		Accept: "*/*",
	}
	result, err := profile.send(context, cmd, options, uripath, nil)
	return result.Reader(), err
}

// Put puts/updates a resource
func (profile *Profile) Put(context context.Context, cmd *cobra.Command, uripath string, body interface{}, response interface{}) (err error) {
	options := &request.Options{Method: http.MethodPut, Payload: body}
	_, err = profile.send(context, cmd, options, uripath, response)
	return
}

// Delete deletes a resource
func (profile *Profile) Delete(context context.Context, cmd *cobra.Command, uripath string, response interface{}) (err error) {
	options := &request.Options{Method: http.MethodDelete}
	_, err = profile.send(context, cmd, options, uripath, response)
	return
}

// Patch patches a resource
func (profile *Profile) Patch(context context.Context, cmd *cobra.Command, uripath string, body interface{}, response interface{}) (err error) {
	options := &request.Options{Method: http.MethodPatch, Payload: body}
	_, err = profile.send(context, cmd, options, uripath, response)
	return
}

// GetAllResources gets all resources of the given type
//
// The Current profile will be set to the profile of the command
func GetAll[T any](context context.Context, cmd *cobra.Command, uripath string) (resources []T, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getall")

	profile, err := GetProfileFromCommand(context, cmd)
	if err != nil {
		log.Errorf("Failed to get profile.", err)
		return nil, err
	}
	Current = profile // Make sure the current profile is set

	log.Infof("Getting all resources for profile %s", profile.Name)
	for {
		var paginated PaginatedResources[T]

		err = profile.Get(
			context,
			cmd,
			uripath,
			&paginated,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, paginated.Values...)
		log.Debugf("Got %d resources", len(paginated.Values))
		log.Debugf("Next page:     %s", paginated.Next)
		log.Debugf("Previous page: %s", paginated.Previous)
		if len(paginated.Next) == 0 {
			break
		}
		uripath = paginated.Next
	}
	return resources, nil
}

// Download downloads a resource to a destination folder
//
// # The destination folder is the current folder if not specified
//
// If the profile has its Progress flag set to true, it will show a progress bar.
// Otherwise, if the command has a flag --progress, it will show a progress bar.
func (profile *Profile) Download(context context.Context, cmd *cobra.Command, uripath, destination string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "download")

	if len(destination) == 0 {
		destination = "."
	}
	if !strings.HasSuffix(destination, "/") {
		destination += "/"
	}
	if err = os.MkdirAll(destination, 0755); err != nil {
		return errors.RuntimeError.Wrap(err)
	}
	writer, err := os.CreateTemp(destination, "artifact-")
	if err != nil {
		return errors.RuntimeError.Wrap(err)
	}

	log.Debugf("Downloading data to %s", writer.Name())
	options := &request.Options{
		Method:              http.MethodGet,
		Timeout:             15 * time.Minute,
		ResponseBodyLogSize: -1, // we are not interested in the file content
	}
	showProgress := profile.Progress
	if cmd != nil && cmd.Flags().Changed("progress") {
		showProgress, _ = cmd.Flags().GetBool("progress")
	}
	if showProgress {
		options.ProgressWriter = profile.getProgressWriter(1, "Downloading")
	}
	result, err := profile.send(context, cmd, options, uripath, writer)
	if err != nil {
		_ = writer.Close()
		return err
	}
	if err = writer.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close file %s: %s\n", writer.Name(), err)
		log.Errorf("Failed to close file %s: %s", writer.Name(), err)
	}
	log.Debugf("Downloaded %d bytes", result.Length)

	filename := result.Headers.Get("Content-Disposition")
	if len(filename) == 0 {
		filename = filepath.Base(uripath)
	} else {
		filename = strings.TrimPrefix(filename, "attachment; filename=\"")
		filename = strings.TrimSuffix(filename, "\"")
	}
	log.Infof("Renaming %s  into %s", writer.Name(), filepath.Join(destination, filename))
	return errors.RuntimeError.Wrap(os.Rename(writer.Name(), filepath.Join(destination, filename)))
}

// Upload uploads a resource from a source file
//
// If the profile has its Progress flag set to true, it will show a progress bar.
// Otherwise, if the command has a flag --progress, it will show a progress bar.
func (profile *Profile) Upload(context context.Context, cmd *cobra.Command, uripath, source string) (err error) {
	reader, err := os.Open(source)
	if err != nil {
		return errors.RuntimeError.Wrap(err)
	}
	defer reader.Close()

	options := &request.Options{
		Method: http.MethodPost,
		Payload: map[string]string{
			">files": filepath.Base(source),
		},
		Attachment:         reader,
		Timeout:            15 * time.Minute,
		RequestBodyLogSize: -1, // we are not interested in the file content
	}
	showProgress := profile.Progress
	if cmd != nil && cmd.Flags().Changed("progress") {
		showProgress, _ = cmd.Flags().GetBool("progress")
	}
	if showProgress {
		var size int64 = -1
		if stat, err := reader.Stat(); err == nil {
			size = stat.Size()
		}
		options.ProgressWriter = profile.getProgressWriter(size, "Upoading")
	}
	_, err = profile.send(context, cmd, options, uripath, nil)
	return
}

func (profile *Profile) getProgressWriter(size int64, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription("[cyan]"+description+"[reset] "),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() { fmt.Fprint(os.Stderr, "\n") }),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

func (profile *Profile) CodeGrantCallback(resultchan chan error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Must(logger.FromContext(r.Context())).Child(nil, nil, "profile", profile.Name)

		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/favicon.ico" {
			_, _ = w.Write([]byte{})
			return
		}

		log.Infof("Received callback from BitBucket")
		code := r.URL.Query().Get("code")
		if len(code) == 0 {
			log.Errorf("No code in the callback")
			http.Error(w, "No code in the callback", http.StatusBadRequest)
			return
		}
		log.Infof("Received code %s", code)

		log.Infof("Requesting authorization token for profile %s", profile.Name)
		result, err := request.Send(&request.Options{
			Method:        http.MethodPost,
			Authorization: request.BasicAuthorization(profile.ClientID, profile.ClientSecret),
			URL:           core.Must(url.Parse("https://bitbucket.org/site/oauth2/access_token")),
			Payload:       map[string]string{"grant_type": "authorization_code", "code": code},
			Timeout:       30 * time.Second,
			Logger:        log,
		}, nil)
		if err != nil {
			if result != nil {
				var errorResponse struct {
					Error            string `json:"error"`
					ErrorDescription string `json:"error_description"`
				}
				if jerr := result.UnmarshalContentJSON(&errorResponse); jerr == nil {
					var details *errors.Error
					status := http.StatusInternalServerError
					if errors.As(err, &details) {
						status = details.Code
					}
					http.Error(w, errorResponse.ErrorDescription, status)
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			resultchan <- err
			return
		}
		profile.saveAccessToken(result.Data)
		_, _ = w.Write([]byte("Authorization Code received. You can close this window."))
		resultchan <- nil
	})
}

func (profile *Profile) authorize(context context.Context) (authorization string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "authorize")

	if err := profile.loadAccessToken(); err == nil {
		if !profile.isTokenExpired() {
			log.Infof("Using access token for profile %s", profile.Name)
			log.Debugf("Token expires on %s in %s", profile.TokenExpires.Format(time.RFC3339), time.Until(profile.TokenExpires))
			return request.BearerAuthorization(profile.AccessToken), nil
		}
	}
	payload := map[string]string{}
	if len(profile.RefreshToken) > 0 {
		log.Warnf("Access token for profile %s is expired", profile.Name)
		payload["grant_type"] = "refresh_token"
		payload["refresh_token"] = profile.RefreshToken
	} else {
		payload["grant_type"] = "client_credentials"
	}

	log.Infof("Authorizing profile %s", profile.Name)
	result, err := request.Send(&request.Options{
		Method:        http.MethodPost,
		Authorization: request.BasicAuthorization(profile.ClientID, profile.ClientSecret),
		URL:           core.Must(url.Parse("https://bitbucket.org/site/oauth2/access_token")),
		Payload:       payload,
		Timeout:       30 * time.Second,
		Logger:        log,
	}, nil)
	if err != nil {
		if result != nil {
			var errorResponse struct {
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
			}
			if jerr := result.UnmarshalContentJSON(&errorResponse); jerr == nil {
				var details *errors.Error
				if errors.As(err, &details) {
					return "", errors.NewSentinel(details.Code, errorResponse.Error, errorResponse.ErrorDescription)
				}
				return "", errors.NewSentinel(500, errorResponse.Error, errorResponse.ErrorDescription)
			}
		}
		return
	}
	profile.saveAccessToken(result.Data)
	return request.BearerAuthorization(profile.AccessToken), nil
}

func (profile *Profile) send(context context.Context, cmd *cobra.Command, options *request.Options, uripath string, response interface{}) (result *request.Content, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, strings.ToLower(options.Method))

	if len(profile.User) > 0 {
		options.Authorization = request.BasicAuthorization(profile.User, profile.Password)
	} else if len(profile.AccessToken) > 0 {
		options.Authorization = request.BearerAuthorization(profile.AccessToken)
	} else if options.Authorization, err = profile.authorize(context); err != nil {
		return nil, err
	}

	apiRoot := profile.APIRoot
	if apiRoot == nil {
		apiRoot = &url.URL{Scheme: "https", Host: "api.bitbucket.org"}
	}

	if strings.HasPrefix(uripath, "/") {
		components := strings.Split(uripath, "?")
		options.URL = apiRoot.JoinPath("2.0", components[0])
		if len(components) > 1 {
			options.URL.RawQuery = components[1]
		}
	} else if !strings.HasPrefix(uripath, "http") {
		repositoryName, err := profile.getRepositoryFullname(context, cmd)
		if err != nil {
			return nil, err
		}
		log.Infof("Using repository %s", repositoryName)
		components := strings.Split(uripath, "?")
		options.URL = apiRoot.JoinPath("2.0", "repositories", repositoryName, components[0])
		if len(components) > 1 {
			options.URL.RawQuery = components[1]
		}
	} else {
		if options.URL, err = url.Parse(uripath); err != nil {
			return nil, err
		}
	}

	if options.Timeout == 0 {
		options.Timeout = 30 * time.Second
	}
	if options.Logger == nil {
		options.Logger = log
	}
	if options.ResponseBodyLogSize == 0 {
		options.ResponseBodyLogSize = 16 * 1024
	}
	if options.ProgressWriter != nil {
		log.Warnf("[B] We have a ProgressWriter for uploading content")
	}
	log.Infof("Sending %s request to %s", options.Method, options.URL)
	result, err = request.Send(options, response)
	if err != nil {
		if result != nil {
			var bberr *BitBucketError
			if jerr := result.UnmarshalContentJSON(&bberr); jerr == nil {
				return result, bberr
			}
		}
	}
	return
}

func (profile Profile) getRepositoryFullname(context context.Context, cmd *cobra.Command) (string, error) {
	log := logger.Must(logger.FromContext(context)).Child("profile", "getrepositoryname")

	fullName := ""
	if cmd != nil && cmd.Flag("repository") != nil {
		fullName = cmd.Flag("repository").Value.String()
	}
	if len(fullName) == 0 {
		log.Debugf("No repository name given, trying to get it from the current git repository")
		remote, err := remote.GetFromGitConfig(context, "origin")
		if err != nil {
			return "", errors.Join(errors.NotFound.With("current repository"), err)
		}
		fullName = remote.RepositoryName()
	}
	return fullName, nil
}
