package lease

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
)

type InboundLeaseManager struct {
	apadanaClient  *apadana.Client
	ttlCheckPeriod time.Duration
}

func NewInboundLeaseManager(apadanaClient *apadana.Client, ttlCheckPeriod time.Duration) *InboundLeaseManager {
	return &InboundLeaseManager{
		apadanaClient:  apadanaClient,
		ttlCheckPeriod: ttlCheckPeriod,
	}
}

func (m *InboundLeaseManager) Tick(ctx context.Context, nodeID string) error {
	inbounds, err := m.apadanaClient.GetInbounds(nodeID)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	now := time.Now()
	for _, inbound := range inbounds {
		wg.Add(1)
		current := inbound
		go func(inb *satrapv1.Inbound) {
			defer wg.Done()
			if now.Sub(inb.Metadata.CreationTimestamp) >= inb.Metadata.TTL {
				if err := m.apadanaClient.DeleteInbound(nodeID, inb.Config.Tag); err != nil {
					zlog.Error().Err(err).Str("component", "lease").Str("resource", "inbound").Str("action", "delete").Str("tag", inb.Config.Tag).Msg("failed")
					return
				}
				zlog.Info().Str("component", "lease").Str("resource", "inbound").Str("tag", inb.Config.Tag).Msg("expired")
			}
		}(current)
	}

	wg.Wait()
	return nil
}

func (m *InboundLeaseManager) Run(ctx context.Context, nodeID string) {
	ticker := time.NewTicker(m.ttlCheckPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "lease").Str("resource", "inbound").Str("action", "tick").Msg("started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.Tick(ctx, nodeID); err != nil {
				zlog.Error().Err(err).Str("component", "lease").Str("resource", "inbound").Str("action", "tick").Msg("failed")
			}
		}
	}
}
