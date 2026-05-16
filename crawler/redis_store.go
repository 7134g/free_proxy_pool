package crawler

import (
	"context"
	"free_proxy_pool/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func redisEnabled() bool {
	return config.RedisClient != nil
}

func redisSyncScore(link string, score int) {
	if !redisEnabled() {
		return
	}
	go func() {
		if err := config.RedisClient.ZAdd(context.Background(), config.Cfg.Redis.Key, redis.Z{
			Score:  float64(score),
			Member: link,
		}).Err(); err != nil {
			log.Printf("Redis ZAdd failed: %v\n", err)
		}
	}()
}

func redisZRem(link string) {
	if !redisEnabled() {
		return
	}
	go func() {
		if err := config.RedisClient.ZRem(context.Background(), config.Cfg.Redis.Key, link).Err(); err != nil {
			log.Printf("Redis ZRem failed: %v\n", err)
		}
	}()
}

func recoverFromRedis() {
	if !redisEnabled() {
		return
	}
	list := redisLoadProxies()
	if len(list) == 0 {
		return
	}
	CacheProxyData.lock.Lock()
	for _, z := range list {
		link, ok := z.Member.(string)
		if !ok {
			continue
		}
		p := newProxy(link)
		p.Score = int(z.Score)
		CacheProxyData.body[link] = p
	}
	CacheProxyData.lock.Unlock()
	CacheProxyData.sort()
	log.Printf("Recovered %d proxies from Redis\n", len(list))
}

func redisLoadProxies() []redis.Z {
	if !redisEnabled() {
		return nil
	}
	members, err := config.RedisClient.ZRevRangeByScoreWithScores(context.Background(), config.Cfg.Redis.Key, &redis.ZRangeBy{
		Min: "1",
		Max: "+inf",
	}).Result()
	if err != nil {
		log.Printf("Redis ZRevRangeByScoreWithScores failed: %v\n", err)
		return nil
	}
	return members
}
