package serviceauth

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedisClient(cfg Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Panicf("Redis connection failed: %v", err)
	}

	log.Println("[INIT] Connected to Redis:", cfg.RedisAddr)
	return rdb
}
