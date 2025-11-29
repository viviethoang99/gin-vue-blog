package handle

import (
	"context"
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"time"
	"github.com/go-redis/redis/v9"
)

// redis context
var rctx = context.Background()

// Page

// Cache page list into Redis
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, g.PAGE, string(data), 0).Err()
}

// Remove page list cache from Redis
func removePageCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.PAGE).Err()
}

// Get page list cache from Redis
// rdb.Get returns redis.Nil error if key does not exist
func getPageCache(rdb *redis.Client) (cache []model.Page, err error) {
	s, err := rdb.Get(rctx, g.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}

	return cache, nil
}

// Config

// Cache blog config into Redis
func addConfigCache(rdb *redis.Client, config map[string]string) error {
	return rdb.HMSet(rctx, g.CONFIG, config).Err()
}

// Remove blog config cache from Redis
func removeConfigCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.CONFIG).Err()
}

// Get blog config cache from Redis
// rdb.HGetAll returns empty map if key does not exist (no redis.Nil)
func getConfigCache(rdb *redis.Client) (cache map[string]string, err error) {
	return rdb.HGetAll(rctx, g.CONFIG).Result()
}

// email
func SetMailInfo (rdb *redis.Client,info string,expire time.Duration) error{
	return rdb.Set(rctx,info,true,expire).Err()
}
func GetMailInfo (rdb *redis.Client,info string) (bool,error){
	return rdb.Get(rctx,info).Bool()
}
func DeleteMailInfo(rdb *redis.Client,info string) error{
	return rdb.Del(rctx,info).Err()
}