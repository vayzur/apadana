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
	concurrentInboundSyncs       int32
	concurrentInboundExpireSyncs int32
	concurrentInboundGCSyncs     int32
	concurrentUserSyncs          int32
	concurrentUserExpireSyncs    int32
}

func NewSyncManager(
	xrayClient *xray.Client,
	apadanaClient *apadana.Client,
	syncFrequency time.Duration,
	concurrentInboundSyncs,
	concurrentInboundExpireSyncs,
	concurrentInboundGCSyncs,
	concurrentUserSyncs,
	concurrentUserExpireSyncs int32,
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
