package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type redisRepo struct {
	rConn *redis.Client
}
type RepositoryStorage interface {
	Set(key, value string, ctx context.Context) error
	Get(key string, ctx context.Context) (string, error)
}

func NewRedisRepo(rds *redis.Client) RepositoryStorage {
	return &redisRepo{
		rConn: rds,
	}

}
func (r *redisRepo) Set(key, value string, ctx context.Context) error {
	err := r.rConn.Set(ctx, key, value, 0).Err()
	return err
}

// SetWithTTL

func (r *redisRepo) Get(key string, ctx context.Context) (string, error) {
	val, err := r.rConn.Get(ctx, key).Result()

	return val, err
}
