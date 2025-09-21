package inbound

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
)

func (m *SyncManager) Run(ctx context.Context, nodeName string) {
	createInboundCh := make(chan *satrapv1.Inbound, 256)
	gcInboundCh := make(chan string, 256)

	createUserCh := make(chan *satrapv1.InboundUser, 256)
	gcUserCh := make(chan *satrapv1.InboundUser, 256)

	for range m.concurrentInboundSyncs {
		go func() {
			for inb := range createInboundCh {
				if err := m.xrayClient.AddInbound(ctx, &inb.Spec.Config); err != nil {
					continue
				}

				desiredUsers, err := m.apadanaClient.GetInboundUsers(nodeName, inb.Spec.Config.Tag)
				if err != nil {
					continue
				}

				for _, user := range desiredUsers {
					createUserCh <- user
				}
			}
		}()
	}

	for range m.concurrentInboundGCSyncs {
		go func() {
			for tag := range gcInboundCh {
				if err := m.xrayClient.RemoveInbound(ctx, tag); err != nil {
					zlog.Error().Err(err).Str("component", "syncManager").Str("controller", "gc").
						Str("resource", "inbound").Str("action", "delete").
						Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
					continue
				}
			}
		}()
	}

	for range m.concurrentUserSyncs {
		go func() {
			for user := range createUserCh {
				account, err := user.ToAccount()
				if err != nil {
					continue
				}
				if err := m.xrayClient.AddUser(ctx, user.Spec.InboundTag, user.Spec.Email, account); err != nil {
					continue
				}
			}
		}()
	}

	for range m.concurrentUserGCSyncs {
		go func() {
			for user := range gcUserCh {
				if err := m.xrayClient.RemoveUser(ctx, user.Spec.InboundTag, user.Spec.Email); err != nil {
					continue
				}
			}
		}()
	}

	ticker := time.NewTicker(m.syncFrequency)
	defer ticker.Stop()

	wg := &sync.WaitGroup{}
	uwg := &sync.WaitGroup{}

	zlog.Info().Str("component", "syncManager").Msg("started")

	for {
		select {
		case <-ctx.Done():
			close(createInboundCh)
			close(gcInboundCh)
			close(createUserCh)
			close(gcUserCh)
			return

		case <-ticker.C:
			desiredInbounds, err := m.apadanaClient.GetInbounds(nodeName)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeName", nodeName).
					Msg("failed to get desired inbounds")
				continue
			}

			currentInbounds, err := m.xrayClient.ListInbounds(ctx)
			if err != nil {
				zlog.Error().Err(err).Str("component", "syncManager").Str("nodeName", nodeName).
					Msg("failed to get current inbounds")
				continue
			}

			desiredInboundsMap := make(map[string]*satrapv1.Inbound, len(desiredInbounds))

			for _, inbound := range desiredInbounds {
				if inbound != nil {
					desiredInboundsMap[inbound.Spec.Config.Tag] = inbound
				}
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				for _, inbound := range desiredInbounds {
					if _, ok := currentInbounds[inbound.Spec.Config.Tag]; !ok {
						createInboundCh <- inbound
						continue
					}

					desiredUsers, err := m.apadanaClient.GetInboundUsers(nodeName, inbound.Spec.Config.Tag)
					if err != nil {
						continue
					}

					currentUsers, err := m.xrayClient.ListUsers(ctx, inbound.Spec.Config.Tag)
					if err != nil {
						continue
					}

					desiredUsersMap := make(map[string]*satrapv1.InboundUser, len(desiredUsers))

					for _, user := range desiredUsers {
						if user != nil {
							desiredUsersMap[user.Spec.Email] = user
						}
					}

					uwg.Go(func() {
						for email := range currentUsers {
							if _, ok := desiredUsersMap[email]; !ok {
								gcUserCh <- &satrapv1.InboundUser{
									Spec: satrapv1.InboundUserSpec{
										InboundTag: inbound.Spec.Config.Tag,
										Email:      email,
									},
								}
							}
						}
					})

					uwg.Go(func() {
						for _, user := range desiredUsers {
							if _, ok := currentUsers[user.Spec.Email]; !ok {
								createUserCh <- user
							}
						}
					})

					uwg.Wait()
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				for tag := range currentInbounds {
					if _, ok := desiredInboundsMap[tag]; !ok {
						gcInboundCh <- tag
					}
				}
			}()

			wg.Wait()
		}
	}
}
