package auth

import (
	"context"

	"github.com/Dreamacro/clash/common/cache"
	"github.com/go-redis/redis/v8"
)

type redisAuthenticator struct {
	redis *redis.Client
	key   string
	cache *cache.LruCache
}

func NewRedisAuthenticator(url, key string, cacheSeconds int64) Authenticator {
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(opts)

	authenticator := &redisAuthenticator{
		redisClient,
		key,
		cache.New(cache.WithAge(cacheSeconds)),
	}
	return authenticator
}

func (ra *redisAuthenticator) Verify(user string, pass string) bool {
	var realPass string
	var err error
	cachedPass, found := ra.cache.Get(user)
	if found {
		realPass = cachedPass.(string)
	} else {
		realPass, err = ra.redis.HGet(context.Background(), ra.key, user).Result()
		if err != nil {
			return false
		}
		ra.cache.Set(user, realPass)
	}
	return pass == realPass
}

func (ra *redisAuthenticator) Users() []string {
	keys, err := ra.redis.HKeys(context.Background(), ra.key).Result()
	if err != nil {
		panic(err)
	}
	return keys
}
