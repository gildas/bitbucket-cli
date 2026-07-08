package remote_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gildas/bitbucket-cli/cmd/remote"
	"github.com/stretchr/testify/assert"
)

func TestCanGetRepositoryNameWithGitAt(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = git@bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = git@bitbucket.org:myworkspace/bitbucket-cli
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())

	r, err = remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())
}

func TestCanGetRepositoryNameWithHTTPS(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = https://bitbucket.org/myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = https://bitbucket.org/myworkspace/bitbucket-cli
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())

	r, err = remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())
}

func TestCanGetWorkspaceNameWithGitAt(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = git@bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = git@bitbucket.org:myworkspace/bitbucket-cli
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace", r.WorkspaceName())

	r, err = remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace", r.WorkspaceName())
}

func TestCanGetWorkspaceNameWithHTTPS(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = https://bitbucket.org/myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = https://bitbucket.org/myworkspace/bitbucket-cli
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace", r.WorkspaceName())

	r, err = remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "myworkspace", r.WorkspaceName())
}
