package controller

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
)

func (c *ControllerManager) RunNodeMonitor(ctx context.Context, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "controller").Msg("node monitor started")
	defer zlog.Info().Str("component", "controller").Msg("node monitor stopped")

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

			now := time.Now()
			for _, node := range nodes {
				if now.Sub(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
					node.Status.Status = false
					if err := c.nodeService.PutNode(ctx, node); err != nil {
						if ctx.Err() != nil {
							return
						}
						zlog.Error().Err(err).Str("component", "controller").Msg("failed to update node status")
						continue
					}
				}
			}
		}
	}
}
