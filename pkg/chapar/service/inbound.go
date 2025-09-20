package service

import (
	"context"
	"errors"

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

func (s *InboundService) CountInbounds(ctx context.Context, nodeName string) (uint32, error) {
	return s.store.CountInbounds(ctx, nodeName)
}

func (s *InboundService) GetInbound(ctx context.Context, nodeName, tag string) (*satrapv1.Inbound, error) {
	return s.store.GetInbound(ctx, nodeName, tag)
}

func (s *InboundService) DeleteInbound(ctx context.Context, nodeName, tag string) error {
	node, err := s.nodeService.GetNode(ctx, nodeName)
	if err != nil {
		return err
	}
	if err := s.satrapClient.RemoveInbound(node, tag); err != nil && !errors.Is(err, errs.ErrInboundNotFound) {
		return err
	}
	if err := s.store.DeleteUsers(ctx, nodeName, tag); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return err
	}
	if err := s.store.DeleteInbound(ctx, nodeName, tag); err != nil && !errors.Is(err, errs.ErrInboundNotFound) {
		return err
	}
	return nil
}

func (s *InboundService) CreateInbound(ctx context.Context, nodeName string, inbound *satrapv1.Inbound) error {
	node, err := s.nodeService.GetNode(ctx, nodeName)
	if err != nil {
		return err
	}
	inboundsCount, err := s.CountInbounds(ctx, nodeName)
	if err != nil {
		return err
	}
	if inboundsCount >= node.Status.Capacity.MaxInbounds {
		return errs.ErrNodeCapacityExceeded
	}
	if err := s.satrapClient.AddInbound(node, &inbound.Spec.Config); err != nil {
		return err
	}
	if err := s.store.CreateInbound(ctx, nodeName, inbound); err != nil {
		if rerr := s.satrapClient.RemoveInbound(node, inbound.Spec.Config.Tag); rerr != nil {
			return rerr
		}
		return err
	}
	return nil
}

func (s *InboundService) GetInbounds(ctx context.Context, nodeName string) ([]*satrapv1.Inbound, error) {
	return s.store.GetInbounds(ctx, nodeName)
}

func (s *InboundService) GetUser(ctx context.Context, nodeName, tag, email string) (*satrapv1.InboundUser, error) {
	return s.store.GetUser(ctx, nodeName, tag, email)
}

func (s *InboundService) CountUsers(ctx context.Context, nodeName, tag string) (uint32, error) {
	return s.store.CountUsers(ctx, nodeName, tag)
}

func (s *InboundService) DeleteUser(ctx context.Context, nodeName, tag, email string) error {
	node, err := s.nodeService.GetNode(ctx, nodeName)
	if err != nil {
		return err
	}
	if err := s.satrapClient.RemoveUser(node, tag, email); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return err
	}
	if err := s.store.DeleteUser(ctx, nodeName, tag, email); err != nil {
		return err
	}
	return nil
}

func (s *InboundService) CreateUser(ctx context.Context, nodeName, tag string, user *satrapv1.InboundUser) error {
	var (
		node       *corev1.Node
		inb        *satrapv1.Inbound
		usersCount uint32
		nodeErr    error
		inbErr     error
		countErr   error
	)

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		node, nodeErr = s.nodeService.GetNode(groupCtx, nodeName)
		return nodeErr
	})
	g.Go(func() error {
		inb, inbErr = s.GetInbound(groupCtx, nodeName, tag)
		return inbErr
	})
	g.Go(func() error {
		usersCount, countErr = s.CountUsers(groupCtx, nodeName, tag)
		return countErr
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if usersCount >= inb.Spec.Capacity.MaxUsers {
		return errs.ErrInboundCapacityExceeded
	}

	if err := s.store.CreateUser(ctx, nodeName, tag, user); err != nil {
		return err
	}

	if err := s.satrapClient.AddUser(node, tag, user); err != nil {
		return err
	}

	return nil
}

func (s *InboundService) GetUsers(ctx context.Context, nodeName, tag string) ([]*satrapv1.InboundUser, error) {
	return s.store.GetUsers(ctx, nodeName, tag)
}

func (s *InboundService) RenewInbound(ctx context.Context, nodeName, tag string, renew *satrapv1.Renew) error {
	inbound, err := s.GetInbound(ctx, nodeName, tag)
	if err != nil {
		return err
	}
	inbound.Metadata.TTL = renew.TTL
	if err := s.store.CreateInbound(ctx, nodeName, inbound); err != nil {
		return err
	}
	return nil
}

func (s *InboundService) RenewInboundUser(ctx context.Context, nodeName, tag, email string, renew *satrapv1.Renew) error {
	user, err := s.GetUser(ctx, nodeName, tag, email)
	if err != nil {
		return err
	}
	user.Metadata.TTL = renew.TTL
	if err := s.store.CreateUser(ctx, nodeName, tag, user); err != nil {
		return err
	}
	return nil
}

func (s *InboundService) UpdateInboundMetadata(ctx context.Context, nodeName, tag string, metadata *satrapv1.Metadata) error {
	inbound, err := s.GetInbound(ctx, nodeName, tag)
	if err != nil {
		return err
	}
	metadata.CreationTimestamp = inbound.Metadata.CreationTimestamp
	inbound.Metadata = *metadata
	return s.store.CreateInbound(ctx, nodeName, inbound)
}

func (s *InboundService) UpdateUserMetadata(ctx context.Context, nodeName, tag, email string, metadata *satrapv1.Metadata) error {
	user, err := s.GetUser(ctx, nodeName, tag, email)
	if err != nil {
		return err
	}
	metadata.CreationTimestamp = user.Metadata.CreationTimestamp
	user.Metadata = *metadata
	return s.store.CreateUser(ctx, nodeName, tag, user)
}
