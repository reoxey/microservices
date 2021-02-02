package profile

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, k string) string
	GetJSON(ctx context.Context, k string, v interface{}) error
	SetEx(ctx context.Context, k string, v interface{}, t time.Duration)
	SetJSON(ctx context.Context, k string, v interface{}, t time.Duration) error
}
