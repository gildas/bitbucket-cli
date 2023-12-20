package remote_test

import (
	"strings"
	"testing"

	"bitbucket.org/gildas_cherruel/bb/cmd/remote"
	"github.com/stretchr/testify/assert"
)

func TestCanGetRepositoryWithGitAt(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = git@bitbucket.org:gildas_cherruel/bb.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = git@bitbucket.org:gildas_cherruel/bb
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.Get(strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "gildas_cherruel/bb", r.RepositoryName())

	r, err = remote.Get(strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "gildas_cherruel/bb", r.RepositoryName())
}

func TestCanGetRepositoryWithHTTPS(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = https://bitbucket.org/gildas_cherruel/bb.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = https://bitbucket.org/gildas_cherruel/bb
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	r, err := remote.Get(strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "gildas_cherruel/bb", r.RepositoryName())

	r, err = remote.Get(strings.NewReader(payload), "alternate")
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "gildas_cherruel/bb", r.RepositoryName())
}
