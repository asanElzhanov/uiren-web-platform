package database

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidTTL = errors.New("invalid ttl")
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
	DataTTL  time.Duration
}

type RedisDB struct {
	Client  *redis.Client
	DataTTL time.Duration
}

func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, ttl *time.Duration) error {
	var dataTTL time.Duration
	if ttl == nil {
		dataTTL = r.DataTTL
	} else {
		dataTTL = *ttl
	}
	return r.Client.Set(ctx, key, value, dataTTL).Err()
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisDB) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func GetRedisDatabase(ctx context.Context, config RedisConfig) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &RedisDB{
		Client:  client,
		DataTTL: config.DataTTL,
	}, nil
}
