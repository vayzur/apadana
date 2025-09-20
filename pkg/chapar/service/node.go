package service

import (
	"context"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	"github.com/vayzur/apadana/pkg/chapar/storage/resources"
)

type NodeService struct {
	store *resources.NodeStore
}

func NewNodeService(store *resources.NodeStore) *NodeService {
	return &NodeService{store: store}
}

func (s *NodeService) GetNode(ctx context.Context, nodeID string) (*corev1.Node, error) {
	return s.store.GetNode(ctx, nodeID)
}

func (s *NodeService) DeleteNode(ctx context.Context, nodeID string) error {
	return s.store.DeleteNode(ctx, nodeID)
}

func (s *NodeService) CreateNode(ctx context.Context, node *corev1.Node) error {
	return s.store.CreateNode(ctx, node)
}

func (s *NodeService) GetNodes(ctx context.Context) ([]*corev1.Node, error) {
	return s.store.GetNodes(ctx)
}

func (s *NodeService) GetActiveNodes(ctx context.Context) ([]*corev1.Node, error) {
	nodes, err := s.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	activeNodes := make([]*corev1.Node, 0, len(nodes)) // preallocated, no zeroing

	n := len(nodes)
	for i := 0; i < n; i++ {
		node := nodes[i]
		if node.Status.Ready {
			activeNodes = append(activeNodes, node)
		}
	}

	return activeNodes, nil
}

func (s *NodeService) UpdateNodeStatus(ctx context.Context, nodeID string, status *corev1.NodeStatus) error {
	node, err := s.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	node.Status = *status
	return s.CreateNode(ctx, node)
}

func (s *NodeService) UpdateNodeMetadata(ctx context.Context, nodeID string, metadata *corev1.NodeMetadata) error {
	node, err := s.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	metadata.ID = node.Metadata.ID
	metadata.CreationTimestamp = node.Metadata.CreationTimestamp

	node.Metadata = *metadata
	return s.CreateNode(ctx, node)
}

func (s *NodeService) UpdateNodeSpec(ctx context.Context, nodeID string, spec *corev1.NodeSpec) error {
	node, err := s.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	node.Spec = *spec
	return s.CreateNode(ctx, node)
}
