package inbound

import (
	"time"

	apadana "github.com/vayzur/apadana/pkg/client"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

type SyncManager struct {
	xrayClient    *xray.Client
	apadanaClient *apadana.Client
	syncFrequency time.Duration
}

func NewSyncManager(xrayClient *xray.Client, apadanaClient *apadana.Client, syncFrequency time.Duration) *SyncManager {
	return &SyncManager{
		xrayClient:    xrayClient,
		apadanaClient: apadanaClient,
		syncFrequency: syncFrequency,
	}
}
