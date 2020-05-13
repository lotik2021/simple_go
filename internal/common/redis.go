package common

import (
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/prometheus/common/log"
)

var (
	redisClient *redis.Client
)

func init() {
	GetRedisConnection()
}

func GetRedisConnection() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.C.Redis.URL,
		DB:   0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return redisClient
}

func CloseRedis() error {
	return redisClient.Close()
}
