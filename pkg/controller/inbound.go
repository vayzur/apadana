package controller

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/apadana/pkg/api/v1"
)

func (c *ControllerManager) RunInboundMonitor(ctx context.Context, inboundMonitorPeriod time.Duration) {
	ticker := time.NewTicker(inboundMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "controller").Msg("inbound monitor started")

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
				zlog.Error().Err(err).Str("component", "controller").Msg("failed to get nodes")
				continue
			}

			var wg sync.WaitGroup

			for _, node := range nodes {
				wg.Add(1)
				currentNode := node
				go func(node *v1.Node) {
					defer wg.Done()

					if ctx.Err() != nil {
						return
					}

					inbounds, err := c.inboundService.ListInbounds(ctx, node)
					if err != nil {
						if ctx.Err() != nil {
							return
						}
						zlog.Error().Err(err).Str("component", "controller").Str("nodeID", node.Metadata.ID).Msg("get inbounds failed")
						return
					}
					now := time.Now()
					for _, inbound := range inbounds {
						if now.Sub(inbound.Metadata.CreationTimestamp) >= inbound.Metadata.TTL {
							if err := c.inboundService.DelInbound(ctx, node, inbound.Config.Tag); err != nil {
								if ctx.Err() != nil {
									return
								}
								zlog.Error().Err(err).Str("component", "controller").Str("nodeID", node.Metadata.ID).Str("tag", inbound.Config.Tag).Msg("delete inbound failed")
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
