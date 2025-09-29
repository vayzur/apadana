package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/internal/config"
	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	satrapconfigv1 "github.com/vayzur/apadana/pkg/satrap/config/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
	"github.com/vayzur/apadana/pkg/satrap/flock"
	satrapHeartbeatManager "github.com/vayzur/apadana/pkg/satrap/health"
	satrapRegisterManager "github.com/vayzur/apadana/pkg/satrap/register"
	satrapSyncManager "github.com/vayzur/apadana/pkg/satrap/sync"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &satrapconfigv1.SatrapConfig{}
	if err := config.Load(*configPath, cfg); err != nil {
		zlog.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	xrayClient, err := xray.New(&cfg.Xray)
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to connect xray")
	}
	defer func() {
		if err := xrayClient.Close(); err != nil {
			zlog.Error().Err(err).Msg("failed to close xray client")
		}
	}()

	apadanaClient := apadana.New(
		cfg.Cluster.Server,
		cfg.Cluster.Token,
		time.Second*5,
	)

	nodeName := cfg.GetName()
	if nodeName == "" {
		zlog.Fatal().
			Err(err).
			Str("component", "registerManager").
			Msg("failed to get node name: cfg.Name not set and system hostname unavailable")
	}

	if cfg.RegisterNode {
		registerManager := satrapRegisterManager.NewRegisterManager(
			apadanaClient,
		)

		labels := map[string]string{
			corev1.LabelHostname: nodeName,
			corev1.LabelOS:       runtime.GOOS,
			corev1.LabelArch:     runtime.GOARCH,
		}

		for k, v := range cfg.Labels {
			labels[k] = v
		}

		node := &corev1.Node{
			Metadata: metav1.ObjectMeta{
				Name:   nodeName,
				Labels: labels,
			},
		}

		rlock := flock.NewFlock("/tmp/satrap-register-manager.lock")
		if err := rlock.TryLock(); err == nil {
			// block until register node
			registerManager.RegisterWithAPIServer(ctx, node)
			defer rlock.Unlock()
		}
	}

	nodeStatus := &corev1.NodeStatus{
		Addresses: cfg.Addresses,
		Capacity:  corev1.NodeCapacity{MaxInbounds: cfg.MaxInbounds},
		Ready:     true,
	}

	hb := satrapHeartbeatManager.NewHeartbeatManager(
		apadanaClient,
		cfg.NodeStatusUpdateFrequency,
		nodeStatus,
	)

	syncManager := satrapSyncManager.NewSyncManager(
		xrayClient,
		apadanaClient,
		cfg.SyncFrequency,
		cfg.ConcurrentInboundSyncs,
		cfg.ConcurrentInboundGCSyncs,
		cfg.ConcurrentUserSyncs,
		cfg.ConcurrentUserGCSyncs,
	)

	hlock := flock.NewFlock("/tmp/satrap-heartbeat.lock")
	if err := hlock.TryLock(); err == nil {
		go hb.Run(ctx, nodeName)
		defer hlock.Unlock()
	}

	slock := flock.NewFlock("/tmp/satrap-sync-manager.lock")
	if err := slock.TryLock(); err == nil {
		go syncManager.Run(ctx, nodeName)
		defer slock.Unlock()
	}

	zlog.Info().Str("component", "satrap").Msg("started")
	<-ctx.Done()
}
