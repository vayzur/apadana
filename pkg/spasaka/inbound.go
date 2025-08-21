package spasaka

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Spasaka) RunInboundMonitor(ctx context.Context, inboundMonitorPeriod time.Duration) {
	ticker := time.NewTicker(inboundMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "spasaka").Str("resource", "inbound").Str("action", "monitor").Msg("started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			nodes, err := c.nodeService.ListActiveNodes(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				zlog.Error().Err(err).Str("component", "spasaka").Str("resource", "nodes").Str("action", "list").Msg("failed")
				continue
			}

			var wg sync.WaitGroup

			for _, node := range nodes {
				wg.Add(1)
				currentNode := node
				go func(node *corev1.Node) {
					defer wg.Done()
					if ctx.Err() != nil {
						return
					}
					inbounds, err := c.inboundService.ListInbounds(ctx, node)
					if err != nil {
						if ctx.Err() != nil {
							return
						}
						zlog.Error().Err(err).Str("component", "spasaka").Str("resource", "inbounds").Str("action", "list").Msg("failed")
						return
					}
					now := time.Now()
					for _, inbound := range inbounds {
						if now.Sub(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
							if err := c.inboundService.DelInbound(ctx, node, inbound.Config.Tag); err != nil {
								if ctx.Err() != nil {
									return
								}
								zlog.Error().Err(err).Str("component", "spasaka").Str("resource", "inbound").Str("action", "delete").Str("nodeID", node.Metadata.ID).Str("tag", inbound.Config.Tag).Msg("failed")
								return
							}
						}
					}
				}(currentNode)
			}
			wg.Wait()
		}
	}
}
