package storage

import "context"

type Interface interface {
	Get(ctx context.Context, key string, out *[]byte) error
	Create(ctx context.Context, key string, obj []byte, ttl uint64) error
	Delete(ctx context.Context, key string) error
	GetList(ctx context.Context, key string, out *[][]byte) error
	Count(ctx context.Context, key string) (uint32, error)
	ReadinessCheck() error
}
