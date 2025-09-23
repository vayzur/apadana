package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"

	"github.com/vayzur/apadana/pkg/chapar/storage/resources"
)

type InboundService struct {
	store *resources.InboundStore
}

func NewInboundService(store *resources.InboundStore) *InboundService {
	return &InboundService{
		store: store,
	}
}

func (s *InboundService) CountInbounds(ctx context.Context, nodeName string) (uint32, error) {
	return s.store.CountInbounds(ctx, nodeName)
}

func (s *InboundService) GetInbound(ctx context.Context, nodeName, tag string) (*satrapv1.Inbound, error) {
	return s.store.GetInbound(ctx, nodeName, tag)
}

func (s *InboundService) DeleteInbound(ctx context.Context, nodeName, tag string) error {
	if err := s.store.DeleteUsers(ctx, nodeName, tag); err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return err
	}
	if err := s.store.DeleteInbound(ctx, nodeName, tag); err != nil && !errors.Is(err, errs.ErrInboundNotFound) {
		return err
	}
	return nil
}

func (s *InboundService) CreateInbound(ctx context.Context, nodeName string, inbound *satrapv1.Inbound) error {
	existingInbound, _ := s.GetInbound(ctx, nodeName, inbound.Spec.Config.Tag)
	if existingInbound != nil {
		return errs.ErrInboundConflict
	}

	inbound.Metadata.UID = uuid.NewString()
	inbound.Metadata.CreationTimestamp = time.Now()

	if err := s.store.CreateInbound(ctx, nodeName, inbound); err != nil {
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
	if err := s.store.DeleteUser(ctx, nodeName, tag, email); err != nil {
		return err
	}
	return nil
}

func (s *InboundService) CreateUser(ctx context.Context, nodeName, tag string, user *satrapv1.InboundUser) error {
	existingUser, _ := s.GetUser(ctx, nodeName, user.Spec.InboundTag, user.Spec.Email)
	if existingUser != nil {
		return errs.ErrUserConflict
	}

	user.Metadata.UID = uuid.NewString()
	user.Metadata.CreationTimestamp = time.Now()

	if err := s.store.CreateUser(ctx, nodeName, tag, user); err != nil {
		return err
	}

	return nil
}

func (s *InboundService) GetUsers(ctx context.Context, nodeName, tag string) ([]*satrapv1.InboundUser, error) {
	return s.store.GetUsers(ctx, nodeName, tag)
}

func (s *InboundService) UpdateInboundMetadata(ctx context.Context, nodeName, tag string, newMetadata *metav1.ObjectMeta) error {
	inbound, err := s.GetInbound(ctx, nodeName, tag)
	if err != nil {
		return err
	}

	newMetadata.Name = inbound.Metadata.Name
	newMetadata.UID = inbound.Metadata.UID
	newMetadata.CreationTimestamp = inbound.Metadata.CreationTimestamp

	inbound.Metadata = *newMetadata
	return s.store.CreateInbound(ctx, nodeName, inbound)
}

func (s *InboundService) UpdateInboundSpec(ctx context.Context, nodeName, tag string, newSpec *satrapv1.InboundSpec) error {
	inbound, err := s.GetInbound(ctx, nodeName, tag)
	if err != nil {
		return err
	}

	newSpec.Config = inbound.Spec.Config

	inbound.Spec = *newSpec
	return s.store.CreateInbound(ctx, nodeName, inbound)
}

func (s *InboundService) UpdateUserMetadata(ctx context.Context, nodeName, tag, email string, newMetadata *metav1.ObjectMeta) error {
	user, err := s.GetUser(ctx, nodeName, tag, email)
	if err != nil {
		return err
	}

	newMetadata.Name = user.Metadata.Name
	newMetadata.UID = user.Metadata.UID
	newMetadata.CreationTimestamp = user.Metadata.CreationTimestamp

	user.Metadata = *newMetadata
	return s.store.CreateUser(ctx, nodeName, tag, user)
}

func (s *InboundService) UpdateUserSpec(ctx context.Context, nodeName, tag, email string, newSpec *satrapv1.InboundUserSpec) error {
	user, err := s.GetUser(ctx, nodeName, tag, email)
	if err != nil {
		return err
	}

	newSpec.Type = user.Spec.Type
	newSpec.InboundTag = user.Spec.InboundTag
	newSpec.Email = user.Spec.Email
	newSpec.Account = user.Spec.Account

	user.Spec = *newSpec
	return s.store.CreateUser(ctx, nodeName, tag, user)
}
