package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	"github.com/vayzur/apadana/pkg/chapar/storage"
	"github.com/vayzur/apadana/pkg/errs"
)

var nodePool = sync.Pool{
	New: func() any { return &corev1.Node{} },
}

type NodeStore struct {
	store storage.Interface
}

func NewNodeStore(store storage.Interface) *NodeStore {
	return &NodeStore{store: store}
}

func (s *NodeStore) GetNode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	key := fmt.Sprintf("/nodes/%s", nodeName)
	out := &[]byte{}

	if err := s.store.Get(ctx, key, out); err != nil {
		if errors.Is(err, errs.ErrResourceNotFound) {
			return nil, errs.ErrNodeNotFound
		}
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get node failed",
			map[string]string{
				"nodeName": nodeName,
			},
			err,
		)
	}

	node := &corev1.Node{}
	if err := json.Unmarshal(*out, node); err != nil {
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnmarshalFailed,
			"get node failed",
			map[string]string{
				"nodeName": nodeName,
			},
			err,
		)
	}

	return node, nil
}

func (s *NodeStore) DeleteNode(ctx context.Context, nodeName string) error {
	key := fmt.Sprintf("/nodes/%s", nodeName)
	if err := s.store.Delete(ctx, key); err != nil {
		if errors.Is(err, errs.ErrResourceNotFound) {
			return errs.ErrResourceNotFound
		}
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"delete node failed",
			map[string]string{
				"nodeName": nodeName,
			},
			err,
		)
	}
	return nil
}

func (s *NodeStore) CreateNode(ctx context.Context, node *corev1.Node) error {
	val, err := json.Marshal(node)
	if err != nil {
		return errs.New(
			errs.KindInternal,
			errs.ReasonMarshalFailed,
			"create node failed",
			nil,
			err,
		)
	}

	key := fmt.Sprintf("/nodes/%s", node.Metadata.Name)
	if err := s.store.Create(ctx, key, val, 0); err != nil {
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"create node failed",
			nil,
			err,
		)
	}

	return nil
}

func (s *NodeStore) GetNodes(ctx context.Context) ([]*corev1.Node, error) {
	key := "/nodes/"
	out := &[][]byte{}

	if err := s.store.GetList(ctx, key, out); err != nil {
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get nodes failed",
			nil,
			err,
		)
	}

	nodes := make([]*corev1.Node, 0, len(*out))

	for _, v := range *out {
		node := nodePool.Get().(*corev1.Node)
		*node = corev1.Node{}

		if err := json.Unmarshal(v, node); err != nil {
			zlog.Error().Err(err).Str("component", "store").Str("resource", "node").Msg("unmarshal failed")
			nodePool.Put(node)
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}
