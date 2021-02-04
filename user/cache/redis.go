package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"user/core"
)

type redisCache struct {
	rdb *redis.Client
}

func (r *redisCache) SetJSON(ctx context.Context, k string, v interface{}, duration time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.SetEx(ctx, k, data, duration)

	return nil
}

func (r *redisCache) GetJSON(ctx context.Context, k string, v interface{}) error {
	return json.Unmarshal(
		[]byte(r.Get(ctx, k)),
		v)
}

func (r *redisCache) Get(ctx context.Context, k string) string {
	return r.rdb.Get(ctx, k).Val()
}

func (r *redisCache) SetEx(ctx context.Context, k string, v interface{}, duration time.Duration) {
	r.rdb.Set(ctx, k, v, duration)
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
