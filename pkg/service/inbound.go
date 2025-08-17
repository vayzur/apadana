package service

import (
	"context"
	"fmt"

	v1 "github.com/vayzur/apadana/pkg/api/v1"
	chapar "github.com/vayzur/apadana/pkg/chapar/client"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

type InboundService struct {
	store        *resources.InboundStore
	chaparClient *chapar.Client
}

func NewInboundService(store *resources.InboundStore, chaparClient *chapar.Client) *InboundService {
	return &InboundService{
		store:        store,
		chaparClient: chaparClient,
	}
}

func (s *InboundService) GetInbound(ctx context.Context, node *v1.Node, tag string) (*v1.Inbound, error) {
	return s.store.GetInbound(ctx, node.Metadata.ID, tag)
}

func (s *InboundService) DelInbound(ctx context.Context, node *v1.Node, tag string) error {
	if err := s.chaparClient.RemoveInbound(node, tag); err != nil {
		return fmt.Errorf("spark delete inbound %s/%s: %w", node.Metadata.ID, tag, err)
	}

	return s.store.DelInbound(ctx, node.Metadata.ID, tag)
}

func (s *InboundService) AddInbound(ctx context.Context, inbound *v1.Inbound, node *v1.Node) error {
	if err := s.chaparClient.AddInbound(&inbound.Config, node); err != nil {
		return fmt.Errorf("spark add inbound %s/%s: %w", node.Metadata.ID, inbound.Config.Tag, err)
	}

	if err := s.store.PutInbound(ctx, node.Metadata.ID, inbound); err != nil {
		if err := s.chaparClient.RemoveInbound(node, inbound.Config.Tag); err != nil {
			return fmt.Errorf("spark add inbound rollback %s/%s: %w", node.Metadata.ID, inbound.Config.Tag, err)
		}
		return err
	}
	return nil
}

func (s *InboundService) ListInbounds(ctx context.Context, node *v1.Node) ([]*v1.Inbound, error) {
	return s.store.ListInbounds(ctx, node.Metadata.ID)
}
