package handle

import (
	"context"
	"gin-blog/internal/model"
	"log"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
)

// Requires a Redis environment
func initRdb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       11,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	return rdb
}

func TestPageCache(t *testing.T) {
	rdb := initRdb()

	pages := []model.Page{
		{Name: "page1"},
		{Name: "page2"},
	}

	// Get cache directly
	// When not exists, returns redis.Nil error
	{
		cache, err := getPageCache(rdb)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, cache)
	}

	// Add cache and get it
	{
		err := addPageCache(rdb, pages)
		assert.Nil(t, err)

		cache, err := getPageCache(rdb)
		assert.Nil(t, err)
		assert.Equal(t, pages, cache)
	}

	// Remove cache and get it
	// When not exists, returns redis.Nil error
	{
		err := removePageCache(rdb)
		assert.Nil(t, err)

		cache, err := getPageCache(rdb)
		assert.Equal(t, redis.Nil, err)
		assert.Nil(t, cache)
	}

}

func TestConfigCache(t *testing.T) {
	rdb := initRdb()

	config := map[string]string{
		"name": "name",
		"url":  "url",
	}

	// Get cache directly
	// When not exists, returns empty map
	{
		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Empty(t, cache)
	}

	// Add cache and get it
	{
		err := addConfigCache(rdb, config)
		assert.Nil(t, err)

		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Equal(t, config, cache)
	}

	// Remove cache and get it
	// When not exists, returns empty map
	{
		err := removeConfigCache(rdb)
		assert.Nil(t, err)

		cache, err := getConfigCache(rdb)
		assert.Nil(t, err)
		assert.Empty(t, cache)
	}
}
