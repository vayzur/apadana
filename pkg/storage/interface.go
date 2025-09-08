package storage

import "context"

type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Create(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
	GetList(ctx context.Context, key string) (map[string][]byte, error)
	Count(ctx context.Context, key string) (uint32, error)
	ReadinessCheck() error
}
