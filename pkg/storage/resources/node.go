package resources

import (
	"context"
	"encoding/json"
	"fmt"

	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/apadana/pkg/api/v1"
	"github.com/vayzur/apadana/pkg/storage"
)

type NodeStore struct {
	store storage.Storage
}

func NewNodeStore(store storage.Storage) *NodeStore {
	return &NodeStore{store: store}
}

func (s *NodeStore) GetNode(ctx context.Context, nodeID string) (*v1.Node, error) {
	key := fmt.Sprintf("/nodes/%s", nodeID)
	resp, err := s.store.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get node %s: %w", nodeID, err)
	}
	var node v1.Node
	if err := json.Unmarshal(resp, &node); err != nil {
		return nil, fmt.Errorf("unmarshal node %s: %w", nodeID, err)
	}

	return &node, nil
}

func (s *NodeStore) DelNode(ctx context.Context, nodeID string) error {
	key := fmt.Sprintf("/nodes/%s", nodeID)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete node %s: %w", nodeID, err)
	}
	return nil
}

func (s *NodeStore) PutNode(ctx context.Context, node *v1.Node) error {
	val, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("marshal node %s: %w", node.Metadata.ID, err)
	}

	key := fmt.Sprintf("/nodes/%s", node.Metadata.ID)
	if err := s.store.Put(ctx, key, string(val)); err != nil {
		return fmt.Errorf("put node %s: %w", node.Metadata.ID, err)
	}

	return nil
}

func (s *NodeStore) ListNodes(ctx context.Context) ([]*v1.Node, error) {
	prefix := "/nodes/"
	resp, err := s.store.List(ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	var nodes []*v1.Node
	for k, v := range resp {
		var node v1.Node
		if err := json.Unmarshal(v, &node); err != nil {
			zlog.Error().Err(err).Str("component", "node").Str("nodeID", k).Msg("unmarshal failed")
			continue
		}
		nodes = append(nodes, &node)
	}

	return nodes, nil
}
