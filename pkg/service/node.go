package service

import (
	"context"

	v1 "github.com/vayzur/apadana/pkg/api/v1"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

type NodeSerivce struct {
	store *resources.NodeStore
}

func NewNodeSerivce(store *resources.NodeStore) *NodeSerivce {
	return &NodeSerivce{store: store}
}

func (s *NodeSerivce) GetNode(ctx context.Context, nodeID string) (*v1.Node, error) {
	return s.store.GetNode(ctx, nodeID)
}

func (s *NodeSerivce) DelNode(ctx context.Context, nodeID string) error {
	return s.store.DelNode(ctx, nodeID)
}

func (s *NodeSerivce) PutNode(ctx context.Context, node *v1.Node) error {
	return s.store.PutNode(ctx, node)
}

func (s *NodeSerivce) ListNodes(ctx context.Context) ([]*v1.Node, error) {
	return s.store.ListNodes(ctx)
}

func (s *NodeSerivce) ListActiveNodes(ctx context.Context) ([]*v1.Node, error) {
	nodes, err := s.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	var activeNodes []*v1.Node
	for _, node := range nodes {
		if node.Status.Status {
			activeNodes = append(activeNodes, node)
		}
	}

	return activeNodes, nil
}

func (s *NodeSerivce) UpdateNodeStatus(ctx context.Context, nodeID string, status *v1.NodeStatus) error {
	node, err := s.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	node.Status = *status
	return s.PutNode(ctx, node)
}
