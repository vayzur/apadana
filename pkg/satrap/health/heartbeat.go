package health

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	apadana "github.com/vayzur/apadana/pkg/client"
)

type HeartbeatManager struct {
	apadanaClient             *apadana.Client
	nodeStatusUpdateFrequency time.Duration
}

func NewHeartbeatManager(apadanaClient *apadana.Client, nodeStatusUpdateFrequency time.Duration) *HeartbeatManager {
	return &HeartbeatManager{
		apadanaClient:             apadanaClient,
		nodeStatusUpdateFrequency: nodeStatusUpdateFrequency,
	}
}

func (h *HeartbeatManager) Run(ctx context.Context, nodeID string) {
	ticker := time.NewTicker(h.nodeStatusUpdateFrequency)
	defer ticker.Stop()
	nodeStatus := &corev1.NodeStatus{
		Ready: true,
	}

	zlog.Info().Str("component", "health").Str("resource", "node").Str("action", "heartbeat").Msg("started")
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			nodeStatus.LastHeartbeatTime = time.Now()

			if err := h.apadanaClient.UpdateNodeStatus(nodeID, nodeStatus); err != nil {
				zlog.Error().Err(err).Str("component", "health").Str("resource", "node").Str("action", "heartbeat").Msg("failed")
				continue
			}
		}
	}
}
