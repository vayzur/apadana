package inbound

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
)

func (m *SyncManager) Tick(ctx context.Context, nodeID string) error {
	desiredInbounds, err := m.apadanaClient.GetInbounds(nodeID)
	if err != nil {
		return err
	}

	currentInbounds, err := m.xrayClient.ListInbounds(ctx)
	if err != nil {
		return err
	}

	desiredMap := make(map[string]*satrapv1.Inbound, len(desiredInbounds))
	for _, inbound := range desiredInbounds {
		if inbound != nil {
			desiredMap[inbound.Config.Tag] = inbound
		}
	}

	now := time.Now()
	wg := &sync.WaitGroup{}

	go func() {
		for tag, inbound := range desiredMap {
			if _, ok := currentInbounds[tag]; !ok {
				wg.Add(1)
				go func(tag string, inb *satrapv1.Inbound) {
					defer wg.Done()
					if err := m.apadanaClient.CreateInbound(nodeID, inbound); err != nil {
						return
					}
				}(tag, inbound)
			} else if now.Sub(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
				wg.Add(1)
				go func(inb *satrapv1.Inbound) {
					defer wg.Done()
					if err := m.apadanaClient.DeleteInbound(nodeID, inb.Config.Tag); err != nil {
						zlog.Error().Err(err).Str("component", "syncManager").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", inb.Config.Tag).Msg("failed")
						return
					}
					zlog.Info().Str("component", "syncManager").Str("resource", "inbound").Str("nodeID", nodeID).Str("tag", inb.Config.Tag).Msg("expired")
				}(inbound)
			}
		}
	}()

	go func() {
		for tag := range currentInbounds {
			if _, ok := desiredMap[tag]; !ok {
				wg.Add(1)
				go func(tag string) {
					defer wg.Done()
					if err := m.xrayClient.RemoveInbound(ctx, tag); err != nil {
						return
					}
				}(tag)
			}
		}
	}()

	wg.Wait()
	return nil
}

func (m *SyncManager) Run(ctx context.Context, nodeID string) {
	ticker := time.NewTicker(m.syncFrequency)
	defer ticker.Stop()

	zlog.Info().Str("component", "syncManager").Str("action", "tick").Msg("started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.Tick(ctx, nodeID); err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("action", "tick").Msg("failed")
			}
		}
	}
}
