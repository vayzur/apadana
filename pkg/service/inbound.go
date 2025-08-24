package service

import (
	"context"
	"fmt"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"

	satrap "github.com/vayzur/apadana/pkg/satrap/client"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

type InboundService struct {
	store        *resources.InboundStore
	nodeService  *NodeSerivce
	satrapClient *satrap.Client
}

func NewInboundService(store *resources.InboundStore, nodeService *NodeSerivce, satrapClient *satrap.Client) *InboundService {
	return &InboundService{
		store:        store,
		nodeService:  nodeService,
		satrapClient: satrapClient,
	}
}

func (s *InboundService) InboundsCount(ctx context.Context, nodeID string) (*satrapv1.Count, error) {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	return s.satrapClient.InboundsCount(node)
}

func (s *InboundService) GetInbound(ctx context.Context, nodeID, tag string) (*satrapv1.Inbound, error) {
	return s.store.GetInbound(ctx, nodeID, tag)
}

func (s *InboundService) DelInbound(ctx context.Context, nodeID, tag string) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	if err := s.satrapClient.RemoveInbound(node, tag); err != nil {
		return fmt.Errorf("inbound delete %s/%s: %w", nodeID, tag, err)
	}
	if err := s.store.DelInbound(ctx, nodeID, tag); err != nil {
		return fmt.Errorf("inbound delete store %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundService) AddInbound(ctx context.Context, nodeID string, inbound *satrapv1.Inbound) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	inboundsCount, err := s.InboundsCount(ctx, nodeID)
	if err != nil {
		return err
	}

	if inboundsCount.Value >= node.Status.Capacity.MaxInbounds {
		return errs.ErrNodeCapacity
	}

	if err := s.satrapClient.AddInbound(node, &inbound.Config); err != nil {
		return fmt.Errorf("inbound add %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}

	if err := s.store.PutInbound(ctx, nodeID, inbound); err != nil {
		if rerr := s.satrapClient.RemoveInbound(node, inbound.Config.Tag); rerr != nil {
			return fmt.Errorf("inbound add rollback %s/%s failed: %w: %w", nodeID, inbound.Config.Tag, rerr, err)
		}
		return fmt.Errorf("inbound add store %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}
	return nil
}

func (s *InboundService) ListInbounds(ctx context.Context, nodeID string) ([]*satrapv1.Inbound, error) {
	return s.store.ListInbounds(ctx, nodeID)
}

func (s *InboundService) GetUser(ctx context.Context, nodeID, tag, email string) (*satrapv1.InboundUser, error) {
	return s.store.GetUser(ctx, nodeID, tag, email)
}

func (s *InboundService) DelUser(ctx context.Context, nodeID, tag, email string) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	if err := s.satrapClient.RemoveUser(node, tag, email); err != nil {
		return fmt.Errorf("inbound user delete %s/%s: %w", nodeID, tag, err)
	}
	if err := s.store.DelUser(ctx, nodeID, tag, email); err != nil {
		return fmt.Errorf("inbound user delete store %s/%s/%s: %w", nodeID, tag, email, err)
	}
	return nil
}

func (s *InboundService) AddUser(ctx context.Context, nodeID, tag string, user *satrapv1.InboundUser) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	if err := s.satrapClient.AddUser(node, tag, user); err != nil {
		return fmt.Errorf("inbound user add %s/%s/%s: %w", nodeID, tag, user.Email, err)
	}

	if err := s.store.PutUser(ctx, nodeID, tag, user); err != nil {
		if rerr := s.satrapClient.RemoveUser(node, tag, user.Email); rerr != nil {
			return fmt.Errorf("inbound user add rollback %s/%s/%s failed: %w: %w", nodeID, tag, user.Email, rerr, err)
		}
		return fmt.Errorf("inbound user add store %s/%s/%s: %w", nodeID, tag, user.Email, err)
	}
	return nil
}

func (s *InboundService) ListUsers(ctx context.Context, nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	return s.store.ListUsers(ctx, nodeID, tag)
}

func (s *InboundService) InboundRenew(ctx context.Context, nodeID, tag string, renew *satrapv1.Renew) error {
	inbound, err := s.GetInbound(ctx, nodeID, tag)
	if err != nil {
		return err
	}

	inbound.Metadata.TTL = renew.TTL

	if err := s.store.PutInbound(ctx, nodeID, inbound); err != nil {
		return fmt.Errorf("inbound renew store %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundService) InboundUserRenew(ctx context.Context, nodeID, tag, email string, renew *satrapv1.Renew) error {
	user, err := s.GetUser(ctx, nodeID, tag, email)
	if err != nil {
		return err
	}

	user.Metadata.TTL = renew.TTL

	if err := s.store.PutUser(ctx, nodeID, tag, user); err != nil {
		return fmt.Errorf("inbound user renew store %s/%s: %w", nodeID, tag, err)
	}
	return nil
}
