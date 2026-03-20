package profile_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"bitbucket.org/gildas_cherruel/bb/cmd/profile"
	"github.com/gildas/go-logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ID string `json:"id"`
}

func TestGetAll_OriginalQueryIsPreservedForNextMissingParams(t *testing.T) {
	oldCurrent := profile.Current
	defer func() { profile.Current = oldCurrent }()

	const filter = `target.ref_name="my-branch"`
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if r.URL.Path == "/pipelines" {
			if r.URL.Query().Get("page") == "" {
				assert.Equal(t, filter, q, "initial request should include original q")
				resp := map[string]interface{}{
					"values": []map[string]string{{"id": "1"}},
					"next":   server.URL + "/pipelines?page=2&pagelen=1",
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			assert.Equal(t, filter, q, "second request should include original q even when next omits it")
			resp := map[string]interface{}{
				"values": []map[string]string{{"id": "2"}},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	apiRoot, err := url.Parse(server.URL)
	assert.NoError(t, err)
	profile.Current = &profile.Profile{APIRoot: apiRoot, DefaultPageLength: 0, AccessToken: "fake-token"}

	cmd := &cobra.Command{}
	cmd.Flags().String("profile", "", "")
	cmd.Flags().Int("page-length", 0, "")
	ctx := logger.Create("test").ToContext(context.Background())
	items, err := profile.GetAll[testItem](ctx, cmd, server.URL+"/pipelines?pagelen=1&q="+url.QueryEscape(filter))
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "1", items[0].ID)
	assert.Equal(t, "2", items[1].ID)
}

func TestGetAll_DoesNotOverwriteExistingNextParams(t *testing.T) {
	oldCurrent := profile.Current
	defer func() { profile.Current = oldCurrent }()

	const originalFilter = `target.ref_name="original"`
	const nextFilter = `target.ref_name="different"`
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if r.URL.Path == "/pipelines" {
			if r.URL.Query().Get("page") == "" {
				assert.Equal(t, originalFilter, q)
				resp := map[string]interface{}{
					"values": []map[string]string{{"id": "1"}},
					"next":   server.URL + "/pipelines?page=2&pagelen=1&q=" + url.QueryEscape(nextFilter),
				}
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			assert.Equal(t, nextFilter, q, "existing q on next URL must not be overwritten")
			resp := map[string]interface{}{
				"values": []map[string]string{{"id": "2"}},
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	apiRoot, err := url.Parse(server.URL)
	assert.NoError(t, err)
	profile.Current = &profile.Profile{APIRoot: apiRoot, DefaultPageLength: 0, AccessToken: "fake-token"}

	cmd := &cobra.Command{}
	cmd.Flags().String("profile", "", "")
	cmd.Flags().Int("page-length", 0, "")
	ctx := logger.Create("test").ToContext(context.Background())
	items, err := profile.GetAll[testItem](ctx, cmd, server.URL+"/pipelines?pagelen=1&q="+url.QueryEscape(originalFilter))
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "1", items[0].ID)
	assert.Equal(t, "2", items[1].ID)
}
