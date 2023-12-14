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

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
)

type PaginatedResources[T any] struct {
	Values   []T    `json:"values"`
	Page     int    `json:"page"`
	PageSize int    `json:"pagelen"`
	Size     int    `json:"size"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func (profile *Profile) Post(context context.Context, repository, uripath string, body interface{}, response interface{}) (err error) {
	return profile.send(context, http.MethodPost, repository, uripath, body, response)
}

func (profile *Profile) Get(context context.Context, repository, uripath string, response interface{}) (err error) {
	return profile.send(context, http.MethodGet, repository, uripath, nil, response)
}

func (profile *Profile) Put(context context.Context, repository, uripath string, body interface{}, response interface{}) (err error) {
	return profile.send(context, http.MethodPut, repository, uripath, body, response)
}

func (profile *Profile) Delete(context context.Context, repository, uripath string, response interface{}) (err error) {
	return profile.send(context, http.MethodDelete, repository, uripath, nil, response)
}

func (profile *Profile) Patch(context context.Context, repository, uripath string, body interface{}, response interface{}) (err error) {
	return profile.send(context, http.MethodPatch, repository, uripath, body, response)
}

// GetAllResources gets all resources using the given profile
func GetAll[T any](context context.Context, profile *Profile, repository string, uripath string) (resources []T, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "getall")

	log.Infof("Getting all resources for profile %s", profile.Name)

	for {
		var paginated PaginatedResources[T]

		err = profile.Get(
			context,
			repository,
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

func (profile *Profile) Download(context context.Context, repository, uripath, destination string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "download")

	if len(repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		repository = remote.Repository()
	}

	var authorization string

	if len(profile.User) > 0 {
		authorization = request.BasicAuthorization(profile.User, profile.Password)
	} else if len(profile.AccessToken) > 0 {
		authorization = request.BearerAuthorization(profile.AccessToken)
	} else if authorization, err = profile.authorize(context); err != nil {
		return err
	}

	if strings.HasPrefix(uripath, "/") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0%s", uripath)
	} else if !strings.HasPrefix(uripath, "http") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", repository, uripath)
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

	log.Infof("Downloading artifact %s to repository %s with profile %s", uripath, repository, profile.Name)
	options := &request.Options{
		Method:        http.MethodGet,
		URL:           core.Must(url.Parse(uripath)),
		Authorization: authorization,
		Timeout:       30 * time.Second,
		Logger:        log,
	}
	result, err := request.Send(options, nil)
	if err != nil {
		if result != nil {
			var bberr *BitBucketError
			if jerr := result.UnmarshalContentJSON(&bberr); jerr == nil {
				return bberr
			}
		}
		return err
	}
	filename := result.Headers.Get("Content-Disposition")
	if len(filename) == 0 {
		filename = filepath.Base(uripath)
	} else {
		filename = strings.TrimPrefix(filename, "attachment; filename=\"")
		filename = strings.TrimSuffix(filename, "\"")
	}
	return errors.RuntimeError.Wrap(os.WriteFile(filepath.Join(destination, filename), result.Data, 0644))
}

func (profile *Profile) Upload(context context.Context, repository, uripath, source string) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "upload")

	if len(repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		repository = remote.Repository()
	}

	var authorization string

	if len(profile.User) > 0 {
		authorization = request.BasicAuthorization(profile.User, profile.Password)
	} else if len(profile.AccessToken) > 0 {
		authorization = request.BearerAuthorization(profile.AccessToken)
	} else if authorization, err = profile.authorize(context); err != nil {
		return err
	}

	if strings.HasPrefix(uripath, "/") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0%s", uripath)
	} else if !strings.HasPrefix(uripath, "http") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", repository, uripath)
	}

	reader, err := os.Open(source)
	if err != nil {
		return errors.RuntimeError.Wrap(err)
	}
	defer reader.Close()

	log.Infof("Uploading artifact %s to repository %s with profile %s", source, repository, profile.Name)
	options := &request.Options{
		Method:        http.MethodPost,
		URL:           core.Must(url.Parse(uripath)),
		Authorization: authorization,
		Payload: map[string]string{
			">files": filepath.Base(source),
		},
		Attachment: reader,
		Timeout:    30 * time.Second,
		Logger:     log,
	}
	result, err := request.Send(options, nil)
	if err != nil {
		if result != nil {
			var bberr *BitBucketError
			if jerr := result.UnmarshalContentJSON(&bberr); jerr == nil {
				return bberr
			}
		}
	}
	return
}

func (profile *Profile) authorize(context context.Context) (authorization string, err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, "authorize")

	if err := profile.loadAccessToken(); err == nil {
		if !profile.isTokenExpired() {
			log.Infof("Using access token for profile %s", profile.Name)
			log.Debugf("Token exires on %s in %s", profile.TokenExpires.Format(time.RFC3339), time.Until(profile.TokenExpires))
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

func (profile *Profile) send(context context.Context, method, repository, uripath string, body interface{}, response interface{}) (err error) {
	log := logger.Must(logger.FromContext(context)).Child(nil, strings.ToLower(method))

	if len(repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		repository = remote.Repository()
	}

	var authorization string

	if len(profile.User) > 0 {
		authorization = request.BasicAuthorization(profile.User, profile.Password)
	} else if len(profile.AccessToken) > 0 {
		authorization = request.BearerAuthorization(profile.AccessToken)
	} else if authorization, err = profile.authorize(context); err != nil {
		return err
	}

	if strings.HasPrefix(uripath, "/") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0%s", uripath)
	} else if !strings.HasPrefix(uripath, "http") {
		uripath = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", repository, uripath)
	}

	options := &request.Options{
		Method:        method,
		URL:           core.Must(url.Parse(uripath)),
		Authorization: authorization,
		Payload:       body,
		Timeout:       30 * time.Second,
		Logger:        log,
	}
	result, err := request.Send(options, &response)
	if err != nil {
		if result != nil {
			var bberr *BitBucketError
			if jerr := result.UnmarshalContentJSON(&bberr); jerr == nil {
				return bberr
			}
		}
	}
	return
}
