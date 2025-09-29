package controller

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
)

func (c *Spasaka) RunNodeMonitor(ctx context.Context, concurrentNodeSyncs int, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	nodesCh := make(chan *corev1.Node)

	for range concurrentNodeSyncs {
		go func() {
			for node := range nodesCh {
				if time.Since(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
					node.Status.Ready = false
					if err := c.apadanaClient.UpdateNodeStatus(node.Metadata.Name, &node.Status); err != nil {
						if ctx.Err() != nil {
							return
						}
						continue
					}
				}
			}
		}()
	}

	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "nodeController").Msg("started")

	for {
		select {
		case <-ctx.Done():
			close(nodesCh)
			return
		case <-ticker.C:
			nodes, err := c.apadanaClient.GetActiveNodes()
			if err != nil {
				if ctx.Err() != nil {
					close(nodesCh)
					return
				}
				continue
			}

			zlog.Info().Str("component", "nodeController").Int("count", len(nodes)).Msg("retrieved")

			for _, node := range nodes {
				if ctx.Err() != nil {
					close(nodesCh)
					return
				}
				nodesCh <- node
			}
		}
	}
}
