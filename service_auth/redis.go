package serviceauth

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	GetServiceKey(service string) (string, error)
	SetServiceKey(service string, key string) error
	HasServiceKey(service string) (bool, error)
}

type redisClientImpl struct {
	client *redis.Client
}

func NewRedisClientWrapper(client *redis.Client) RedisClient {
	return &redisClientImpl{client: client}
}

// GetServiceKey zwraca przypisany klucz API danego serwisu
func (r *redisClientImpl) GetServiceKey(service string) (string, error) {
	key := fmt.Sprintf("serviceauth:key:%s", service)
	return r.client.Get(context.Background(), key).Result()
}

// SetServiceKey ustawia klucz API dla danego serwisu (bez TTL na razie)
func (r *redisClientImpl) SetServiceKey(service string, apiKey string) error {
	key := fmt.Sprintf("serviceauth:key:%s", service)
	return r.client.Set(context.Background(), key, apiKey, 0).Err()
}

// HasServiceKey sprawdza, czy dany serwis ma przypisany klucz API
func (r *redisClientImpl) HasServiceKey(service string) (bool, error) {
	key := fmt.Sprintf("serviceauth:key:%s", service)
	res, err := r.client.Exists(context.Background(), key).Result()
	return res > 0, err
}
