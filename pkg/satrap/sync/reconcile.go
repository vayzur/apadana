package inbound

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
)

func (m *SyncManager) Run(ctx context.Context, nodeID string) {
	add := make(chan *satrapv1.Inbound, 256)
	expire := make(chan *satrapv1.Inbound, 256)
	gc := make(chan string, 256)

	for range m.concurrentInboundSyncs {
		go func() {
			for inb := range add {
				if err := m.apadanaClient.CreateInbound(nodeID, inb); err != nil {
					continue
				}
			}
		}()
	}

	for range m.concurrentExpireSyncs {
		go func() {
			for inb := range expire {
				if err := m.apadanaClient.DeleteInbound(nodeID, inb.Config.Tag); err != nil {
					continue
				}
				zlog.Info().Str("component", "syncManager").Str("resource", "inbound").Str("nodeID", nodeID).Str("tag", inb.Config.Tag).Msg("expired")
			}
		}()
	}

	for range m.concurrentGCSyncs {
		go func() {
			for tag := range gc {
				if err := m.xrayClient.RemoveInbound(ctx, tag); err != nil {
					zlog.Error().Err(err).Str("component", "syncManager").Str("controller", "gc").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
					continue
				}
			}
		}()
	}

	ticker := time.NewTicker(m.syncFrequency)
	defer ticker.Stop()

	wg := &sync.WaitGroup{}

	zlog.Info().Str("component", "syncManager").Str("action", "tick").Msg("started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			desiredInbounds, err := m.apadanaClient.GetInbounds(nodeID)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeID", nodeID).Msg("failed to get desired inbounds")
				continue
			}

			currentInbounds, err := m.xrayClient.ListInbounds(ctx)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeID", nodeID).Msg("failed to get current inbounds")
				continue
			}

			desiredMap := make(map[string]*satrapv1.Inbound, len(desiredInbounds))
			for _, inbound := range desiredInbounds {
				if inbound != nil {
					desiredMap[inbound.Config.Tag] = inbound
				}
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag, inbound := range desiredMap {
					if _, ok := currentInbounds[tag]; !ok {
						add <- inbound
					} else if time.Since(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
						expire <- inbound
					}
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag := range currentInbounds {
					if _, ok := desiredMap[tag]; !ok {
						gc <- tag
					}
				}
			}()

			wg.Wait()
		}
	}
}
