package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"

	satrap "github.com/vayzur/apadana/pkg/satrap/client"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

type InboundService struct {
	store        *resources.InboundStore
	nodeService  *NodeService
	satrapClient *satrap.Client
}

func NewInboundService(store *resources.InboundStore, nodeService *NodeService, satrapClient *satrap.Client) *InboundService {
	return &InboundService{
		store:        store,
		nodeService:  nodeService,
		satrapClient: satrapClient,
	}
}

func (s *InboundService) GetInboundsCount(ctx context.Context, nodeID string) (*satrapv1.Count, error) {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	return s.satrapClient.InboundsCount(node)
}

func (s *InboundService) GetInbound(ctx context.Context, nodeID, tag string) (*satrapv1.Inbound, error) {
	return s.store.GetInbound(ctx, nodeID, tag)
}

func (s *InboundService) DeleteInbound(ctx context.Context, nodeID, tag string) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	if err := s.store.DeleteUsers(ctx, nodeID, tag); err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("inbound delete store %s/%s: %w", nodeID, tag, err)
	}
	if err := s.store.DeleteInbound(ctx, nodeID, tag); err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("inbound delete store %s/%s: %w", nodeID, tag, err)
	}
	if err := s.satrapClient.RemoveInbound(node, tag); err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("inbound delete runtime %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundService) CreateInbound(ctx context.Context, nodeID string, inbound *satrapv1.Inbound) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	inboundsCount, err := s.GetInboundsCount(ctx, nodeID)
	if err != nil {
		return err
	}
	if inboundsCount.Value >= node.Status.Capacity.MaxInbounds {
		return errs.ErrNodeCapacity
	}
	if err := s.satrapClient.AddInbound(node, &inbound.Config); err != nil {
		return fmt.Errorf("inbound add runtime %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}
	if err := s.store.CreateInbound(ctx, nodeID, inbound); err != nil {
		if rerr := s.satrapClient.RemoveInbound(node, inbound.Config.Tag); rerr != nil {
			return fmt.Errorf("inbound add rollback %s/%s failed: %w: %w", nodeID, inbound.Config.Tag, rerr, err)
		}
		return fmt.Errorf("inbound add store %s/%s: %w", nodeID, inbound.Config.Tag, err)
	}
	return nil
}

func (s *InboundService) GetInbounds(ctx context.Context, nodeID string) ([]*satrapv1.Inbound, error) {
	return s.store.GetInbounds(ctx, nodeID)
}

func (s *InboundService) GetExpiredInbounds(ctx context.Context, nodeID string) ([]*satrapv1.Inbound, error) {
	inbounds, err := s.GetInbounds(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	expired := make([]*satrapv1.Inbound, 0, len(inbounds)) // preallocated, no zeroing
	now := time.Now()                                      // only once

	n := len(inbounds)
	for i := 0; i < n; i++ {
		inbound := inbounds[i]
		if now.Sub(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
			expired = append(expired, inbound)
		}
	}

	return expired, nil
}

func (s *InboundService) GetActiveInbounds(ctx context.Context, nodeID string) ([]*satrapv1.Inbound, error) {
	inbounds, err := s.GetInbounds(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	active := make([]*satrapv1.Inbound, 0, len(inbounds)) // preallocated, no zeroing
	now := time.Now()                                     // only once

	n := len(inbounds)
	for i := 0; i < n; i++ {
		inbound := inbounds[i]
		if now.Sub(inbound.Metadata.CreationTimestamp) < inbound.Metadata.TTL {
			active = append(active, inbound)
		}
	}

	return active, nil
}

func (s *InboundService) GetUser(ctx context.Context, nodeID, tag, email string) (*satrapv1.InboundUser, error) {
	return s.store.GetUser(ctx, nodeID, tag, email)
}

func (s *InboundService) DeleteUser(ctx context.Context, nodeID, tag, email string) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	if err := s.satrapClient.RemoveUser(node, tag, email); err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("inbound user delete runtime %s/%s: %w", nodeID, tag, err)
	}
	if err := s.store.DeleteUser(ctx, nodeID, tag, email); err != nil {
		return fmt.Errorf("inbound user delete store %s/%s/%s: %w", nodeID, tag, email, err)
	}
	return nil
}

func (s *InboundService) CreateUser(ctx context.Context, nodeID, tag string, user *satrapv1.InboundUser) error {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}
	if err := s.satrapClient.AddUser(node, tag, user); err != nil {
		return fmt.Errorf("create inbound user runtime %s/%s/%s: %w", nodeID, tag, user.Email, err)
	}
	if err := s.store.CreateUser(ctx, nodeID, tag, user); err != nil {
		if rerr := s.satrapClient.RemoveUser(node, tag, user.Email); rerr != nil && !errors.Is(rerr, errs.ErrNotFound) {
			return fmt.Errorf("create inbound user rollback %s/%s/%s failed: %w: %w", nodeID, tag, user.Email, rerr, err)
		}
		return fmt.Errorf("create inbound user store %s/%s/%s: %w", nodeID, tag, user.Email, err)
	}
	return nil
}

func (s *InboundService) GetUsers(ctx context.Context, nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	return s.store.GetUsers(ctx, nodeID, tag)
}

func (s *InboundService) GetExpiredUsers(ctx context.Context, nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	users, err := s.GetUsers(ctx, nodeID, tag)
	if err != nil {
		return nil, err
	}

	expired := make([]*satrapv1.InboundUser, 0, len(users)) // preallocated, no zeroing
	now := time.Now()

	n := len(users)
	for i := 0; i < n; i++ {
		user := users[i]
		if now.Sub(user.Metadata.CreationTimestamp) >= user.Metadata.TTL {
			expired = append(expired, user)
		}
	}

	return expired, nil
}

func (s *InboundService) GetActiveUsers(ctx context.Context, nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	users, err := s.GetUsers(ctx, nodeID, tag)
	if err != nil {
		return nil, err
	}

	active := make([]*satrapv1.InboundUser, 0, len(users)) // preallocated, no zeroing
	now := time.Now()

	n := len(users)
	for i := 0; i < n; i++ {
		user := users[i]
		if now.Sub(user.Metadata.CreationTimestamp) < user.Metadata.TTL {
			active = append(active, user)
		}
	}

	return active, nil
}

func (s *InboundService) InboundRenew(ctx context.Context, nodeID, tag string, renew *satrapv1.Renew) error {
	inbound, err := s.GetInbound(ctx, nodeID, tag)
	if err != nil {
		return err
	}
	inbound.Metadata.TTL = renew.TTL
	if err := s.store.CreateInbound(ctx, nodeID, inbound); err != nil {
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
	if err := s.store.CreateUser(ctx, nodeID, tag, user); err != nil {
		return fmt.Errorf("inbound user renew store %s/%s: %w", nodeID, tag, err)
	}
	return nil
}
