package storage

import "context"

type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Create(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
	GetList(ctx context.Context, key string) (map[string][]byte, error)
	ReadinessCheck() error
}
