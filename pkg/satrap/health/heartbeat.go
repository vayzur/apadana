package health

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"

	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	apadana "github.com/vayzur/apadana/pkg/client"
)

type HeartbeatManager struct {
	apadanaClient             *apadana.Client
	nodeStatusUpdateFrequency time.Duration
	nodeStatus                *corev1.NodeStatus
}

func NewHeartbeatManager(
	apadanaClient *apadana.Client,
	nodeStatusUpdateFrequency time.Duration,
	nodeStatus *corev1.NodeStatus,
) *HeartbeatManager {
	return &HeartbeatManager{
		apadanaClient:             apadanaClient,
		nodeStatusUpdateFrequency: nodeStatusUpdateFrequency,
		nodeStatus:                nodeStatus,
	}
}

func (h *HeartbeatManager) Run(ctx context.Context, nodeName string) {
	ticker := time.NewTicker(h.nodeStatusUpdateFrequency)
	defer ticker.Stop()

	zlog.Info().Str("component", "heartbeatManager").Msg("started")
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.nodeStatus.LastHeartbeatTime = time.Now()

			if err := h.apadanaClient.UpdateNodeStatus(nodeName, h.nodeStatus); err != nil {
				zlog.Error().Err(err).Str("component", "heartbeatManager").Str("resource", "node").Str("action", "heartbeat").Str("nodeName", nodeName).Msg("failed")
				continue
			}
		}
	}
}
