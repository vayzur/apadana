package controller

import (
	"context"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Spasaka) RunNodeMonitor(ctx context.Context, concurrentNodeSyncs int, nodeMonitorPeriod, nodeMonitorGracePeriod time.Duration) {
	ticker := time.NewTicker(nodeMonitorPeriod)
	defer ticker.Stop()

	zlog.Info().Str("component", "spasaka").Str("resource", "node").Str("action", "monitor").Msg("started")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			nodes, err := c.apadanaClient.GetActiveNodes()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				zlog.Error().Err(err).Str("component", "spasaka").Str("resource", "nodes").Str("action", "list").Msg("failed")
				continue
			}

			var wg sync.WaitGroup
			tasksChan := make(chan *corev1.Node, len(nodes))

			now := time.Now()

			for i := 0; i < concurrentNodeSyncs; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for node := range tasksChan {
						if now.Sub(node.Status.LastHeartbeatTime) >= nodeMonitorGracePeriod {
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

			for _, node := range nodes {
				tasksChan <- node
			}

			close(tasksChan)
			wg.Wait()
		}
	}
}
