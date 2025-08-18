package spasaka

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Spasaka) RunNodeMonitor(ctx context.Context, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "spasaka").Msg("node monitor started")
	defer zlog.Info().Str("component", "spasaka").Msg("node monitor stopped")

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
				zlog.Error().Err(err).Str("component", "spasaka").Msg("failed to get nodes")
				continue
			}

			var wg sync.WaitGroup

			now := time.Now()
			for _, node := range nodes {
				wg.Add(1)
				currentNode := node
				go func(node *corev1.Node) {
					defer wg.Done()
					if now.Sub(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
						node.Status.Status = false
						if err := c.nodeService.PutNode(ctx, node); err != nil {
							if ctx.Err() != nil {
								return
							}
							zlog.Error().Err(err).Str("component", "spasaka").Msg("failed to update node status")
							return
						}
					}
				}(currentNode)
			}
			wg.Wait()
		}
	}
}
