package stores

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const KEY = "passwords"

type RedisStore struct {
	RDB *redis.Client
}

func (r *RedisStore) Initialize(ctx context.Context, errorRate float64, capacity int) error {
	err := r.RDB.FlushAll(ctx).Err()
	if err != nil {
		return err
	}

	err = r.RDB.BFReserve(ctx, KEY, errorRate, int64(capacity)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) Add(ctx context.Context, passwords []any) error {
	return r.RDB.BFMAdd(ctx, KEY, passwords...).Err()
}

func (r *RedisStore) Exists(ctx context.Context, password []byte) (bool, error) {
	res := r.RDB.BFExists(ctx, KEY, password)
	if err := res.Err(); err != nil {
		return false, err
	}
	return res.Val(), nil
}
