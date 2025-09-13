package inbound

import (
	"time"

	apadana "github.com/vayzur/apadana/pkg/client"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

type SyncManager struct {
	xrayClient               *xray.Client
	apadanaClient            *apadana.Client
	syncFrequency            time.Duration
	concurrentInboundSyncs   uint32
	concurrentInboundGCSyncs uint32
	concurrentUserSyncs      uint32
	concurrentUserGCSyncs    uint32
}

func NewSyncManager(
	xrayClient *xray.Client,
	apadanaClient *apadana.Client,
	syncFrequency time.Duration,
	concurrentInboundSyncs,
	concurrentInboundGCSyncs,
	concurrentUserSyncs,
	concurrentUserGCSyncs uint32,
) *SyncManager {
	return &SyncManager{
		xrayClient:               xrayClient,
		apadanaClient:            apadanaClient,
		syncFrequency:            syncFrequency,
		concurrentInboundSyncs:   concurrentInboundSyncs,
		concurrentInboundGCSyncs: concurrentInboundGCSyncs,
		concurrentUserSyncs:      concurrentUserSyncs,
		concurrentUserGCSyncs:    concurrentUserGCSyncs,
	}
}
