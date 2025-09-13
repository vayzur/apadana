package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
	"golang.org/x/sync/errgroup"

	"github.com/vayzur/apadana/pkg/chapar/storage/resources"
	satrap "github.com/vayzur/apadana/pkg/satrap/client"
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

func (s *InboundService) CountRuntimeInbounds(ctx context.Context, nodeID string) (*satrapv1.Count, error) {
	node, err := s.nodeService.GetNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	return s.satrapClient.CountInbounds(node)
}

func (s *InboundService) CountInbounds(ctx context.Context, nodeID string) (uint32, error) {
	return s.store.CountInbounds(ctx, nodeID)
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
	inboundsCount, err := s.CountInbounds(ctx, nodeID)
	if err != nil {
		return err
	}
	if inboundsCount >= node.Status.Capacity.MaxInbounds {
		return errs.ErrCapacityExceeded
	}
	if err := s.satrapClient.AddInbound(node, &inbound.Spec.Config); err != nil {
		return fmt.Errorf("create inbound runtime %s/%s: %w", nodeID, inbound.Spec.Config.Tag, err)
	}
	if err := s.store.CreateInbound(ctx, nodeID, inbound); err != nil {
		if rerr := s.satrapClient.RemoveInbound(node, inbound.Spec.Config.Tag); rerr != nil {
			return fmt.Errorf("create inbound rollback %s/%s failed: %w: %w", nodeID, inbound.Spec.Config.Tag, rerr, err)
		}
		return fmt.Errorf("create inbound store %s/%s: %w", nodeID, inbound.Spec.Config.Tag, err)
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

func (s *InboundService) CountUsers(ctx context.Context, nodeID, tag string) (uint32, error) {
	return s.store.CountUsers(ctx, nodeID, tag)
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
	var (
		node       *corev1.Node
		inb        *satrapv1.Inbound
		usersCount uint32
		nodeErr    error
		inbErr     error
		countErr   error
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		node, nodeErr = s.nodeService.GetNode(ctx, nodeID)
		return nodeErr
	})
	g.Go(func() error {
		inb, inbErr = s.GetInbound(ctx, nodeID, tag)
		return inbErr
	})
	g.Go(func() error {
		usersCount, countErr = s.CountUsers(ctx, nodeID, tag)
		return countErr
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if usersCount >= inb.Spec.Capacity.MaxUsers {
		return errs.ErrCapacityExceeded
	}

	if err := s.store.CreateUser(ctx, nodeID, tag, user); err != nil {
		return fmt.Errorf("create inbound user store %s/%s/%s: %w", nodeID, tag, user.Email, err)
	}

	if err := s.satrapClient.AddUser(node, tag, user); err != nil {
		return fmt.Errorf("create inbound user runtime %s/%s/%s: %w", nodeID, tag, user.Email, err)
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

func (s *InboundService) RenewInbound(ctx context.Context, nodeID, tag string, renew *satrapv1.Renew) error {
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

func (s *InboundService) RenewInboundUser(ctx context.Context, nodeID, tag, email string, renew *satrapv1.Renew) error {
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

func (s *InboundService) UpdateInboundMetadata(ctx context.Context, nodeID, tag string, metadata *satrapv1.Metadata) error {
	inbound, err := s.GetInbound(ctx, nodeID, tag)
	if err != nil {
		return err
	}

	inbound.Metadata = *metadata
	return s.store.CreateInbound(ctx, nodeID, inbound)
}

func (s *InboundService) UpdateUserMetadata(ctx context.Context, nodeID, tag, email string, metadata *satrapv1.Metadata) error {
	user, err := s.GetUser(ctx, nodeID, tag, email)
	if err != nil {
		return err
	}

	user.Metadata = *metadata
	return s.store.CreateUser(ctx, nodeID, tag, user)
}
