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

func TestCanGetRepositoryNameWithSSHHostname(t *testing.T) {
	payload := `
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = git@ssh.bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "alternate"]
	url = git@ssh.bitbucket.org:myworkspace/bitbucket-cli
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "ssh"]
	url = ssh://git@ssh.bitbucket.org/myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[remote "nouser"]
	url = ssh.bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "dev"]
	remote = origin
	merge = refs/heads/dev
	`
	for _, name := range []string{"origin", "alternate", "ssh", "nouser"} {
		r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), name)
		assert.NoError(t, err, "remote %s", name)
		if assert.NotNil(t, r, "remote %s", name) {
			assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName(), "remote %s", name)
			assert.Equal(t, "myworkspace", r.WorkspaceName(), "remote %s", name)
		}
	}
}

func TestCanFindFirstBitbucketRemoteWithSSHHostname(t *testing.T) {
	payload := `
[remote "github"]
	url = git@github.com:someone/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/github/*
[remote "bitbucket"]
	url = git@ssh.bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/bitbucket/*
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "")
	assert.NoError(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())
	}
}

func TestShouldRejectNonBitbucketRemote(t *testing.T) {
	payload := `
[remote "origin"]
	url = git@github.com:someone/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
	`
	_, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.Error(t, err)
}

func TestCanGetRepositoryNameWithSSHHostnameAndPort(t *testing.T) {
	payload := `
[remote "origin"]
	url = ssh://git@ssh.bitbucket.org:22/myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/origin/*
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "origin")
	assert.NoError(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())
	}
}

func TestShouldRejectSSHHostnameSuffix(t *testing.T) {
	payload := `
[remote "scp"]
	url = git@ssh.bitbucket.org.example.com:other/project.git
	fetch = +refs/heads/*:refs/remotes/scp/*
[remote "ssh"]
	url = ssh://git@ssh.bitbucket.org.example.com/other/project.git
	fetch = +refs/heads/*:refs/remotes/ssh/*
	`
	for _, name := range []string{"scp", "ssh"} {
		_, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), name)
		assert.Error(t, err, "remote %s", name)
	}
	_, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "")
	assert.Error(t, err)
}

func TestShouldSkipSSHHostnameSuffixWhenFindingFirstBitbucketRemote(t *testing.T) {
	payload := `
[remote "impostor"]
	url = git@ssh.bitbucket.org.example.com:other/project.git
	fetch = +refs/heads/*:refs/remotes/impostor/*
[remote "bitbucket"]
	url = git@ssh.bitbucket.org:myworkspace/bitbucket-cli.git
	fetch = +refs/heads/*:refs/remotes/bitbucket/*
	`
	r, err := remote.GetRemoteFromReader(context.Background(), strings.NewReader(payload), "")
	assert.NoError(t, err)
	if assert.NotNil(t, r) {
		assert.Equal(t, "myworkspace/bitbucket-cli", r.RepositoryName())
	}
}
