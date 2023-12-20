package profile

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/repository"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
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
func GetAll[T any](context context.Context, cmd *cobra.Command, profile *Profile, uripath string) (resources []T, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getall")

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

// Dowmload downloads a resource to a destination folder
func (profile *Profile) Download(context context.Context, cmd *cobra.Command, uripath, destination string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "download")

	log.Infof("Downloading %s", uripath)
	options := &request.Options{Method: http.MethodGet}
	result, err := profile.send(context, cmd, options, uripath, nil)
	if err != nil {
		return err
	}
	filename := result.Headers.Get("Content-Disposition")
	if len(filename) == 0 {
		filename = filepath.Base(uripath)
	} else {
		filename = strings.TrimPrefix(filename, "attachment; filename=\"")
		filename = strings.TrimSuffix(filename, "\"")
	}
	if len(destination) == 0 {
		destination = "."
	}
	if !strings.HasSuffix(destination, "/") {
		destination += "/"
	}
	if err = os.MkdirAll(destination, 0755); err != nil {
		return errors.RuntimeError.Wrap(err)
	}
	log.Infof("Writing data to %s", filepath.Join(destination, filename))
	return errors.RuntimeError.Wrap(os.WriteFile(filepath.Join(destination, filename), result.Data, 0644))
}

// Upload uploads a resource from a source file
func (profile *Profile) Upload(context context.Context, cmd *cobra.Command, uripath, source string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "upload")

	reader, err := os.Open(source)
	if err != nil {
		return errors.RuntimeError.Wrap(err)
	}
	defer reader.Close()

	log.Infof("Uploading %s", source)
	options := &request.Options{
		Method: http.MethodPost,
		Payload: map[string]string{
			">files": filepath.Base(source),
		},
	}
	_, err = profile.send(context, cmd, options, uripath, nil)
	return
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

	if strings.HasPrefix(uripath, "/") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0%s", uripath)
	} else if !strings.HasPrefix(uripath, "http") {
		repository, err := repository.GetRepository(context, cmd)
		if err != nil {
			return nil, err
		}
		log.Infof("Using repository %s", repository)
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", repository, uripath)
	}

	options.URL, err = url.Parse(uripath)
	if err != nil {
		return nil, err
	}
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Second
	}
	if options.Logger == nil {
		options.Logger = log
	}
	options.ResponseBodyLogSize = 16 * 1024
	result, err = request.Send(options, &response)
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
