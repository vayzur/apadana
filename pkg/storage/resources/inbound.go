package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/storage"
)

var (
	inboundPool = sync.Pool{
		New: func() any { return &satrapv1.Inbound{} },
	}
	userPool = sync.Pool{
		New: func() any { return &satrapv1.InboundUser{} },
	}
)

type InboundStore struct {
	store storage.Storage
}

func NewInboundStore(store storage.Storage) *InboundStore {
	return &InboundStore{store: store}
}

func (s *InboundStore) GetInbound(ctx context.Context, nodeID, tag string) (*satrapv1.Inbound, error) {
	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, tag)
	resp, err := s.store.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get inbound %s/%s: %w", nodeID, tag, err)
	}
	var inbound satrapv1.Inbound
	if err := json.Unmarshal(resp, &inbound); err != nil {
		return nil, fmt.Errorf("unmarshal inbound %s/%s: %w", nodeID, tag, err)
	}

	return &inbound, nil
}

func (s *InboundStore) CreateInbound(ctx context.Context, nodeID string, inbound *satrapv1.Inbound) error {
	val, err := json.Marshal(inbound)
	if err != nil {
		return fmt.Errorf("marshal inbound %s/%s: %w", nodeID, inbound.Spec.Config.Tag, err)
	}

	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, inbound.Spec.Config.Tag)
	if err := s.store.Create(ctx, key, string(val)); err != nil {
		return fmt.Errorf("create inbound %s/%s: %w", nodeID, inbound.Spec.Config.Tag, err)
	}

	return nil
}

func (s *InboundStore) DeleteInbound(ctx context.Context, nodeID, tag string) error {
	key := fmt.Sprintf("/inbounds/%s/%s", nodeID, tag)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete inbound %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundStore) GetInbounds(ctx context.Context, nodeID string) ([]*satrapv1.Inbound, error) {
	key := fmt.Sprintf("/inbounds/%s/", nodeID)
	resp, err := s.store.GetList(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("list inbounds %s: %w", nodeID, err)
	}

	inbounds := make([]*satrapv1.Inbound, 0, len(resp))

	for k, v := range resp {
		inbound := inboundPool.Get().(*satrapv1.Inbound)
		*inbound = satrapv1.Inbound{}

		if err := json.Unmarshal(v, inbound); err != nil {
			zlog.Error().Err(err).Str("component", "inbound").Str("nodeID", nodeID).Str("tag", k).Msg("unmarshal failed")
			inboundPool.Put(inbound)
			continue
		}
		inbounds = append(inbounds, inbound)
	}

	return inbounds, nil
}

func (s *InboundStore) CountInbounds(ctx context.Context, nodeID string) (uint32, error) {
	key := fmt.Sprintf("/inbounds/%s/", nodeID)
	count, err := s.store.Count(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("count inbounds %s: %w", nodeID, err)
	}

	return count, nil
}

func (s *InboundStore) GetUser(ctx context.Context, nodeID, tag, email string) (*satrapv1.InboundUser, error) {
	key := fmt.Sprintf("/inboundUsers/%s/%s/%s", nodeID, tag, email)
	resp, err := s.store.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("get inbound user %s/%s/%s: %w", nodeID, tag, email, err)
	}
	var u satrapv1.InboundUser
	if err := json.Unmarshal(resp, &u); err != nil {
		return nil, fmt.Errorf("unmarshal inbound user %s/%s/%s: %w", nodeID, tag, email, err)
	}

	return &u, nil
}

func (s *InboundStore) CreateUser(ctx context.Context, nodeID, tag string, inboundUser *satrapv1.InboundUser) error {
	val, err := json.Marshal(inboundUser)
	if err != nil {
		return fmt.Errorf("marshal inbound user %s/%s: %w", nodeID, tag, err)
	}

	key := fmt.Sprintf("/inboundUsers/%s/%s/%s", nodeID, tag, inboundUser.Email)
	if err := s.store.Create(ctx, key, string(val)); err != nil {
		return fmt.Errorf("create inbound user %s/%s/%s: %w", nodeID, tag, inboundUser.Email, err)
	}

	return nil
}

func (s *InboundStore) DeleteUser(ctx context.Context, nodeID, tag, email string) error {
	key := fmt.Sprintf("/inboundUsers/%s/%s/%s", nodeID, tag, email)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete inbound user %s/%s/%s: %w", nodeID, tag, email, err)
	}
	return nil
}

func (s *InboundStore) DeleteUsers(ctx context.Context, nodeID, tag string) error {
	key := fmt.Sprintf("/inboundUsers/%s/%s/", nodeID, tag)
	if err := s.store.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete inbound users %s/%s: %w", nodeID, tag, err)
	}
	return nil
}

func (s *InboundStore) GetUsers(ctx context.Context, nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	key := fmt.Sprintf("/inboundUsers/%s/%s/", nodeID, tag)
	resp, err := s.store.GetList(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("list inbound users %s: %w", nodeID, err)
	}

	users := make([]*satrapv1.InboundUser, 0, len(resp))

	for k, v := range resp {
		user := userPool.Get().(*satrapv1.InboundUser)
		*user = satrapv1.InboundUser{}

		if err := json.Unmarshal(v, user); err != nil {
			zlog.Error().Err(err).Str("component", "inbound").Str("nodeID", nodeID).Str("tag", k).Msg("user unmarshal failed")
			userPool.Put(user)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *InboundStore) CountUsers(ctx context.Context, nodeID, tag string) (uint32, error) {
	key := fmt.Sprintf("/inboundUsers/%s/%s/", nodeID, tag)
	count, err := s.store.Count(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("count inbound users %s: %w", nodeID, err)
	}

	return count, nil
}
