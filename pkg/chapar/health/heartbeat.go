package health

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"

	v1 "github.com/vayzur/apadana/pkg/api/v1"
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

func (h *HeartbeatManager) StartHeartbeat(nodeID string, ctx context.Context) {
	ticker := time.NewTicker(h.nodeStatusUpdateFrequency)
	defer ticker.Stop()
	nodeStatus := new(v1.NodeStatus)

	zlog.Info().Str("component", "health").Msg("heartbeat started")
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Str("component", "health").Msg("heartbeat stopped")
			return
		case <-ticker.C:
			nodeStatus.Status = true
			nodeStatus.LastHeartbeatTime = time.Now()

			if err := h.apadanaClient.UpdateNodeStatus(nodeID, nodeStatus); err != nil {
				zlog.Error().Err(err).Str("component", "health").Msg("heartbeat failed")
				continue
			}
		}
	}
}
