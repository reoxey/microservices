package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"product/core"
)

type redisCache struct {
	rdb *redis.Client
}

func (r *redisCache) SetJSON(ctx context.Context, s string, i interface{}, duration time.Duration) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	r.SetEx(ctx, s, data, duration)

	return nil
}

func (r *redisCache) GetJSON(ctx context.Context, s string, i interface{}) error {
	return json.Unmarshal(
		[]byte(r.Get(ctx, s)),
		i)
}

func (r *redisCache) Get(ctx context.Context, s string) string {
	return r.rdb.Get(ctx, s).Val()
}

func (r *redisCache) SetEx(ctx context.Context, s string, i interface{}, duration time.Duration) {
	r.rdb.Set(ctx, s, i, duration)
}

func Redis(addr string) core.Cache {
	cl := redis.NewClient(&redis.Options{
		Addr:     addr + ":6379",
		Password: "",
		DB:       0,
	})

	pong, e := cl.Ping(context.Background()).Result()
	fmt.Println(pong, e)

	return &redisCache{cl}
}
