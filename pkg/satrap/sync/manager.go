package inbound

import (
	"time"

	apadana "github.com/vayzur/apadana/pkg/client"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

type SyncManager struct {
	xrayClient                   *xray.Client
	apadanaClient                *apadana.Client
	syncFrequency                time.Duration
	concurrentInboundSyncs       uint32
	concurrentInboundExpireSyncs uint32
	concurrentInboundGCSyncs     uint32
	concurrentUserSyncs          uint32
	concurrentUserExpireSyncs    uint32
}

func NewSyncManager(
	xrayClient *xray.Client,
	apadanaClient *apadana.Client,
	syncFrequency time.Duration,
	concurrentInboundSyncs,
	concurrentInboundExpireSyncs,
	concurrentInboundGCSyncs,
	concurrentUserSyncs,
	concurrentUserExpireSyncs uint32,
) *SyncManager {
	return &SyncManager{
		xrayClient:                   xrayClient,
		apadanaClient:                apadanaClient,
		syncFrequency:                syncFrequency,
		concurrentInboundSyncs:       concurrentInboundSyncs,
		concurrentInboundExpireSyncs: concurrentInboundExpireSyncs,
		concurrentInboundGCSyncs:     concurrentInboundGCSyncs,
		concurrentUserSyncs:          concurrentUserSyncs,
		concurrentUserExpireSyncs:    concurrentUserExpireSyncs,
	}
}
