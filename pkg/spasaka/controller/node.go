package controller

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Spasaka) RunNodeMonitor(ctx context.Context, concurrentNodeSyncs int, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	nodesChan := make(chan *corev1.Node)

	for range concurrentNodeSyncs {
		go func() {
			for node := range nodesChan {
				if time.Since(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
					node.Status.Ready = false
					if err := c.apadanaClient.UpdateNodeStatus(node.Metadata.ID, &node.Status); err != nil {
						if ctx.Err() != nil {
							return
						}
						zlog.Error().Err(err).Str("component", "spasaka").Str("resource", "node").Str("action", "update").Str("nodeID", node.Metadata.ID).Msg("failed")
						continue
					}
				}
			}
		}()
	}

	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "spasaka").Str("resource", "node").Str("action", "monitor").Msg("started")

	for {
		select {
		case <-ctx.Done():
			close(nodesChan)
			return
		case <-ticker.C:
			nodes, err := c.apadanaClient.GetActiveNodes()
			if err != nil {
				if ctx.Err() != nil {
					close(nodesChan)
					return
				}
				continue
			}

			zlog.Info().Int("count", len(nodes)).Msg("retrieved")

			for _, node := range nodes {
				if ctx.Err() != nil {
					close(nodesChan)
					return
				}
				nodesChan <- node
			}
		}
	}
}
