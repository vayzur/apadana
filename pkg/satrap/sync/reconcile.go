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

	for range m.concurrentInboundSyncs {
		go func() {
			for inb := range createInboundCh {
				if err := m.apadanaClient.CreateInbound(nodeID, inb); err != nil {
					continue
				}

				desiredUsers, err := m.apadanaClient.GetInboundUsers(nodeID, inb.Spec.Config.Tag, satrapv1.Active)
				if err != nil {
					continue
				}

				for _, user := range desiredUsers {
					createUserCh <- user
				}
			}
		}()
	}

	for range m.concurrentInboundExpireSyncs {
		go func() {
			for inb := range expireInboundCh {
				if err := m.apadanaClient.DeleteInbound(nodeID, inb.Spec.Config.Tag); err != nil {
					continue
				}
				zlog.Info().Str("component", "syncManager").Str("resource", "inbound").
					Str("nodeID", nodeID).Str("tag", inb.Spec.Config.Tag).Msg("expired")
			}
		}()
	}

	for range m.concurrentInboundGCSyncs {
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

	for range m.concurrentUserExpireSyncs {
		go func() {
			for user := range expireUserCh {
				if err := m.apadanaClient.DeleteInboundUser(nodeID, user.InboundTag, user.Email); err != nil {
					continue
				}
				zlog.Info().Str("component", "syncManager").Str("resource", "user").
					Str("nodeID", nodeID).Str("tag", user.InboundTag).Str("email", user.Email).Msg("expired")
			}
		}()
	}

	ticker := time.NewTicker(m.syncFrequency)
	defer ticker.Stop()

	wg := &sync.WaitGroup{}
	desiredInboundMap := make(map[string]*satrapv1.Inbound)

	zlog.Info().Str("component", "syncManager").Msg("started")

	for {
		select {
		case <-ctx.Done():
			close(createInboundCh)
			close(expireInboundCh)
			close(gcInboundCh)
			close(createUserCh)
			close(expireUserCh)
			return

		case <-ticker.C:
			desiredInbounds, err := m.apadanaClient.GetInbounds(nodeID, satrapv1.Active)
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
					desiredInboundMap[inbound.Spec.Config.Tag] = inbound
				}
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				expiredInbounds, err := m.apadanaClient.GetInbounds(nodeID, satrapv1.Expired)
				if err != nil {
					zlog.Error().Err(err).Str("component", "syncManager").Str("nodeID", nodeID).
						Msg("failed to get expired inbounds")
					return
				}

				for _, inbound := range expiredInbounds {
					expireInboundCh <- inbound
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag, inbound := range desiredInboundMap {
					if _, ok := currentInbounds[tag]; !ok {
						createInboundCh <- inbound
					}

					expiredUsers, err := m.apadanaClient.GetInboundUsers(nodeID, inbound.Spec.Config.Tag, satrapv1.Expired)
					if err != nil {
						continue
					}
					for _, user := range expiredUsers {
						expireUserCh <- user
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
