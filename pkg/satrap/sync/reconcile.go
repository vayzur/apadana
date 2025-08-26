package inbound

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
)

func (m *SyncManager) Run(ctx context.Context, nodeID string) {
	createInboundCh := make(chan *satrapv1.Inbound, 256)
	expireInboundCh := make(chan *satrapv1.Inbound, 256)
	gcInboundCh := make(chan string, 256)

	createUserCh := make(chan *satrapv1.InboundUser, 256)
	expireUserCh := make(chan *satrapv1.InboundUser, 256)
	gcUserCh := make(chan *satrapv1.InboundUser, 256)

	for range m.concurrentInboundSyncs {
		go func() {
			for inb := range createInboundCh {
				if err := m.apadanaClient.CreateInbound(nodeID, inb); err != nil {
					continue
				}

				go func(inb *satrapv1.Inbound) {
					desiredUsers, err := m.apadanaClient.GetInboundUsers(nodeID, inb.Config.Tag)
					if err != nil {
						return
					}

					for _, user := range desiredUsers {
						if time.Since(user.Metadata.CreationTimestamp) >= user.Metadata.TTL {
							expireUserCh <- user
							continue
						}
						createUserCh <- user
					}
				}(inb)
			}
		}()
	}

	for range m.concurrentExpireSyncs {
		go func() {
			for inb := range expireInboundCh {
				if err := m.apadanaClient.DeleteInbound(nodeID, inb.Config.Tag); err != nil {
					continue
				}
				zlog.Info().Str("component", "syncManager").Str("resource", "inbound").
					Str("nodeID", nodeID).Str("tag", inb.Config.Tag).Msg("expired")
			}
		}()
	}

	for range m.concurrentGCSyncs {
		go func() {
			for tag := range gcInboundCh {
				if err := m.xrayClient.RemoveInbound(ctx, tag); err != nil {
					zlog.Error().Err(err).Str("component", "syncManager").Str("controller", "gc").
						Str("resource", "inbound").Str("action", "delete").
						Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
					continue
				}
			}
		}()
	}

	for range m.concurrentUserSyncs {
		go func() {
			for user := range createUserCh {
				if err := m.apadanaClient.CreateInboundUser(nodeID, user.InboundTag, user); err != nil {
					continue
				}
			}
		}()
	}

	for range m.concurrentExpireSyncs {
		go func() {
			for user := range expireUserCh {
				if err := m.apadanaClient.DeleteInboundUser(nodeID, user.InboundTag, user.Email); err != nil {
					continue
				}
				zlog.Info().Str("component", "syncManager").Str("resource", "user").
					Str("nodeID", nodeID).Str("tag", user.InboundTag).Str("user", user.Email).Msg("expired")
			}
		}()
	}

	for range m.concurrentGCSyncs {
		go func() {
			for user := range gcUserCh {
				if err := m.xrayClient.RemoveUser(ctx, user.InboundTag, user.Email); err != nil {
					zlog.Error().Err(err).Str("component", "syncManager").Str("controller", "gc").
						Str("resource", "user").Str("action", "delete").
						Str("nodeID", nodeID).Str("tag", user.InboundTag).Str("user", user.Email).Msg("failed")
					continue
				}
			}
		}()
	}

	ticker := time.NewTicker(m.syncFrequency)
	defer ticker.Stop()

	wg := &sync.WaitGroup{}
	desiredInboundMap := make(map[string]*satrapv1.Inbound)

	zlog.Info().Str("component", "syncManager").Str("action", "tick").Msg("started")

	for {
		select {
		case <-ctx.Done():
			close(createInboundCh)
			close(expireInboundCh)
			close(gcInboundCh)
			close(createUserCh)
			close(expireUserCh)
			close(gcUserCh)
			return

		case <-ticker.C:
			desiredInbounds, err := m.apadanaClient.GetInbounds(nodeID)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeID", nodeID).
					Msg("failed to get desired inbounds")
				continue
			}

			currentInbounds, err := m.xrayClient.ListInbounds(ctx)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeID", nodeID).
					Msg("failed to get current inbounds")
				continue
			}

			clear(desiredInboundMap)
			for _, inbound := range desiredInbounds {
				if inbound != nil {
					desiredInboundMap[inbound.Config.Tag] = inbound
				}
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag, inbound := range desiredInboundMap {
					if _, ok := currentInbounds[tag]; !ok {
						createInboundCh <- inbound
						continue
					}

					if time.Since(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
						expireInboundCh <- inbound
						continue
					}
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag := range currentInbounds {
					if _, ok := desiredInboundMap[tag]; !ok {
						gcInboundCh <- tag
					}
				}
			}()

			wg.Wait()
		}
	}
}
