package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	"github.com/vayzur/apadana/pkg/chapar/storage"
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

func (s *NodeStore) GetNode(ctx context.Context, nodeID string) (*corev1.Node, error) {
	key := fmt.Sprintf("/nodes/%s", nodeID)
	out := &[]byte{}

	if err := s.store.Get(ctx, key, out); err != nil {
		return nil, fmt.Errorf("get node %s: %w", nodeID, err)
	}

	node := &corev1.Node{}
	if err := json.Unmarshal(*out, node); err != nil {
		return nil, fmt.Errorf("unmarshal node %s: %w", nodeID, err)
	}

	return node, nil
}

func (s *NodeStore) DeleteNode(ctx context.Context, nodeID string) error {
	key := fmt.Sprintf("/nodes/%s", nodeID)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete node %s: %w", nodeID, err)
	}
	return nil
}

func (s *NodeStore) CreateNode(ctx context.Context, node *corev1.Node) error {
	val, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("marshal node %s: %w", node.Metadata.ID, err)
	}

	key := fmt.Sprintf("/nodes/%s", node.Metadata.ID)
	if err := s.store.Create(ctx, key, val, 0); err != nil {
		return fmt.Errorf("create node %s: %w", node.Metadata.ID, err)
	}

	return nil
}

func (s *NodeStore) GetNodes(ctx context.Context) ([]*corev1.Node, error) {
	key := "/nodes/"
	out := &[][]byte{}

	if err := s.store.GetList(ctx, key, out); err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
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
