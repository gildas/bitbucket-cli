package common

import (
	"time"

	"github.com/gildas/go-cache"
	"github.com/gildas/go-core"
)

func NewCache[T any]() *cache.Cache[T] {
	return cache.New[T]("bitbucket", cache.CacheOptionPersistent).WithExpiration(core.GetEnvAsDuration("BITBUCKET_CLI_CACHE_DURATION", 5*time.Minute)).WithEncryptionKey([]byte(core.GetEnvAsString("BITBUCKET_CLI_CACHE_ENCRYPTIONKEY", "")))
}
