package profile

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
)

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

func (profile *Profile) send(context context.Context, method, repository, uripath string, body interface{}, response interface{}) (err error) {
	log := logger.Must(logger.FromContext(context, Log)).Child(nil, strings.ToLower(method))

	if len(repository) == 0 {
		remote, err := remote.GetFromGitConfig("origin")
		if err != nil {
			return err
		}
		repository = remote.Repository()
	}

	options := &request.Options{
		Method:        method,
		URL:           core.Must(url.Parse(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", repository, uripath))),
		Authorization: request.BearerAuthorization(profile.AccessToken),
		Timeout:       30 * time.Second,
		Logger:        log,
	}
	result, err := request.Send(options, &response)
	if err != nil {
		return err
	}
	log.Record("result", result).Infof("Success")
	return
}
