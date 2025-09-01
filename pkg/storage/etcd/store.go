package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vayzur/apadana/pkg/errs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdStorage struct {
	client *clientv3.Client
}

func NewEtcdStorage(client *clientv3.Client) *EtcdStorage {
	return &EtcdStorage{
		client: client,
	}
}

func (e *EtcdStorage) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", key, err)
	}

	if len(resp.Kvs) == 0 {
		return nil, errs.ErrNotFound
	}

	return resp.Kvs[0].Value, nil
}

func (e *EtcdStorage) Create(ctx context.Context, key string, value string) error {
	_, err := e.client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("%q: %w", key, err)
	}
	return nil
}

func (e *EtcdStorage) Delete(ctx context.Context, key string) error {
	opts := []clientv3.OpOption{}

	if strings.HasSuffix(key, "/") {
		opts = append(opts, clientv3.WithPrefix())
	}

	resp, err := e.client.Delete(ctx, key, opts...)
	if err != nil {
		return fmt.Errorf("%q: %w", key, err)
	}

	if resp.Deleted == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (e *EtcdStorage) GetList(ctx context.Context, prefix string) (map[string][]byte, error) {
	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("%q: %w", prefix, err)
	}

	result := make(map[string][]byte, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), prefix)
		result[key] = kv.Value
	}

	return result, nil
}

func (e *EtcdStorage) ReadinessCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := e.client.Status(ctx, e.client.Endpoints()[0])
	return err
}
