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

func (e *EtcdStorage) Get(ctx context.Context, key string, out *[]byte) error {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("%q: %w", key, err)
	}

	if len(resp.Kvs) == 0 {
		return errs.ErrResourceNotFound
	}

	*out = resp.Kvs[0].Value
	return nil
}

func (e *EtcdStorage) Create(ctx context.Context, key string, obj []byte, ttl uint64) error {
	var opts []clientv3.OpOption

	if ttl != 0 {
		lease, err := e.client.Grant(ctx, int64(ttl))
		if err != nil {
			return fmt.Errorf("create lease failed %q: %w", key, err)
		}
		opts = append(opts, clientv3.WithLease(lease.ID))
	}

	_, err := e.client.Put(ctx, key, string(obj), opts...)
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
		return errs.ErrResourceNotFound
	}

	return nil
}

func (e *EtcdStorage) GetList(ctx context.Context, prefix string, out *[][]byte) error {
	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	kvs := make([][]byte, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		kvs = append(kvs, kv.Value)
	}

	*out = kvs
	return nil
}

func (e *EtcdStorage) Count(ctx context.Context, key string) (uint32, error) {
	resp, err := e.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, fmt.Errorf("%q: %w", key, err)
	}
	return uint32(resp.Count), nil
}

func (e *EtcdStorage) ReadinessCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lastErr error
	for _, ep := range e.client.Endpoints() {
		_, err := e.client.Status(ctx, ep)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}
