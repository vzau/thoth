package cache

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DisgoOrg/log"
	"github.com/go-redis/redis/v8"
	gocache "github.com/robfig/go-cache"
	"github.com/vzau/common/utils"
)

type cache struct {
	Driver     string
	InMemory   *gocache.Cache
	Redis      *redis.Client
	DefaultTTL time.Duration
}

var Cache *cache

func BuildCache(defaultExpiration int) error {
	Cache = &cache{}

	driver := utils.Getenv("CACHE_DRIVER", "")
	if strings.EqualFold(driver, "redis") {
		Cache.Driver = driver
		db, err := strconv.ParseInt(utils.Getenv("REDIS_DB", "0"), 10, 64)
		if err != nil {
			db = 0
		}
		Cache.Redis = redis.NewClient(&redis.Options{
			Addr:     utils.Getenv("REDIS_ADDR", "localhost:6379"),
			Password: utils.Getenv("REDIS_PASSWORD", ""),
			DB:       int(db),
		})
		Cache.DefaultTTL = time.Duration(defaultExpiration) * time.Second
	} else if strings.EqualFold(driver, "memory") {
		Cache.Driver = driver
		Cache.DefaultTTL = time.Duration(defaultExpiration) * time.Second
		Cache.InMemory = gocache.New(Cache.DefaultTTL, time.Minute)
	} else {
		return errors.New("invalid cache driver")
	}

	return nil
}

func (c *cache) Get(key string) (interface{}, bool) {
	log.Debug("Getting key: %s", key)
	if strings.EqualFold(c.Driver, "redis") {
		value, err := c.Redis.Get(context.Background(), key).Result()
		if err == redis.Nil {
			return nil, false
		}
		if err != nil {
			return nil, false
		}
		return value, true
	} else if strings.EqualFold(c.Driver, "memory") {
		value, found := c.InMemory.Get(key)
		if !found {
			return nil, false
		}
		return value, true
	}
	return nil, false
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.DefaultTTL
	}

	if strings.EqualFold(c.Driver, "redis") {
		err := c.Redis.Set(context.Background(), key, value, ttl).Err()
		if err != nil {
			return err
		}
	} else if strings.EqualFold(c.Driver, "memory") {
		c.InMemory.Set(key, value, ttl)
	}
	return nil
}

func (c *cache) Delete(key string) error {
	log.Debug("Deleting key: %s", key)
	if strings.EqualFold(c.Driver, "redis") {
		err := c.Redis.Del(context.Background(), key).Err()
		if err != nil {
			return err
		}
	} else if strings.EqualFold(c.Driver, "memory") {
		c.InMemory.Delete(key)
	}

	return nil
}
