package service

import (
	"context"
	"fmt"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"

	satrap "github.com/vayzur/apadana/pkg/satrap/client"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

type InboundService struct {
	store        *resources.InboundStore
	satrapClient *satrap.Client
}

func NewInboundService(store *resources.InboundStore, satrapClient *satrap.Client) *InboundService {
	return &InboundService{
		store:        store,
		satrapClient: satrapClient,
	}
}

func (s *InboundService) InboundsCount(ctx context.Context, node *corev1.Node) (*satrapv1.Count, error) {
	return s.satrapClient.InboundsCount(node)
}

func (s *InboundService) GetInbound(ctx context.Context, node *corev1.Node, tag string) (*satrapv1.Inbound, error) {
	return s.store.GetInbound(ctx, node.Metadata.ID, tag)
}

func (s *InboundService) DelInbound(ctx context.Context, node *corev1.Node, tag string) error {
	if err := s.satrapClient.RemoveInbound(node, tag); err != nil {
		return fmt.Errorf("inbound delete %s/%s: %w", node.Metadata.ID, tag, err)
	}
	if err := s.store.DelInbound(ctx, node.Metadata.ID, tag); err != nil {
		return fmt.Errorf("inbound delete store %s/%s: %w", node.Metadata.ID, tag, err)
	}
	return nil
}

func (s *InboundService) AddInbound(ctx context.Context, inbound *satrapv1.Inbound, node *corev1.Node) error {
	if err := s.satrapClient.AddInbound(&inbound.Config, node); err != nil {
		return fmt.Errorf("inbound add %s/%s: %w", node.Metadata.ID, inbound.Config.Tag, err)
	}

	if err := s.store.PutInbound(ctx, node.Metadata.ID, inbound); err != nil {
		if rerr := s.satrapClient.RemoveInbound(node, inbound.Config.Tag); rerr != nil {
			return fmt.Errorf("inbound add rollback %s/%s failed: %w: %w", node.Metadata.ID, inbound.Config.Tag, rerr, err)
		}
		return fmt.Errorf("inbound add store %s/%s: %w", node.Metadata.ID, inbound.Config.Tag, err)
	}
	return nil
}

func (s *InboundService) ListInbounds(ctx context.Context, node *corev1.Node) ([]*satrapv1.Inbound, error) {
	return s.store.ListInbounds(ctx, node.Metadata.ID)
}

func (s *InboundService) GetUser(ctx context.Context, node *corev1.Node, tag, email string) (*satrapv1.InboundUser, error) {
	return s.store.GetUser(ctx, node.Metadata.ID, tag, email)
}

func (s *InboundService) DelUser(ctx context.Context, node *corev1.Node, tag, email string) error {
	if err := s.satrapClient.RemoveUser(node, tag, email); err != nil {
		return fmt.Errorf("inbound user delete %s/%s: %w", node.Metadata.ID, tag, err)
	}
	if err := s.store.DelUser(ctx, node.Metadata.ID, tag, email); err != nil {
		return fmt.Errorf("inbound user delete store %s/%s/%s: %w", node.Metadata.ID, tag, email, err)
	}
	return nil
}

func (s *InboundService) AddUser(ctx context.Context, node *corev1.Node, tag string, user *satrapv1.InboundUser) error {
	if err := s.satrapClient.AddUser(node, tag, user); err != nil {
		return fmt.Errorf("inbound user add %s/%s/%s: %w", node.Metadata.ID, tag, user.Email, err)
	}
	if err := s.store.PutUser(ctx, node.Metadata.ID, tag, user); err != nil {
		if rerr := s.satrapClient.RemoveUser(node, tag, user.Email); rerr != nil {
			return fmt.Errorf("inbound user add rollback %s/%s/%s failed: %w: %w", node.Metadata.ID, tag, user.Email, rerr, err)
		}
		return fmt.Errorf("inbound user add store %s/%s/%s: %w", node.Metadata.ID, tag, user.Email, err)
	}
	return nil
}

func (s *InboundService) ListUsers(ctx context.Context, node *corev1.Node, tag string) ([]*satrapv1.InboundUser, error) {
	return s.store.ListUsers(ctx, node.Metadata.ID, tag)
}
