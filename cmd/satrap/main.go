package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/internal/config"
	"github.com/vayzur/apadana/internal/satrap/server"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	satrapconfigv1 "github.com/vayzur/apadana/pkg/satrap/config/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
	"github.com/vayzur/apadana/pkg/satrap/flock"
	"github.com/vayzur/apadana/pkg/satrap/health"
	satrapSyncManager "github.com/vayzur/apadana/pkg/satrap/sync"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := satrapconfigv1.SatrapConfig{}
	if err := config.Load(*configPath, &cfg); err != nil {
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
			zlog.Error().Err(err).Msg("xray client close failed")
		}
	}()

	serverAddr := fmt.Sprintf("%s:%d", cfg.BindAddress, cfg.Port)
	app := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, xrayClient)

	var scheme string

	if cfg.TLS.Enabled {
		scheme = "https"
		go func() {
			if err := app.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil {
				zlog.Fatal().Err(err).Msg("server failed")
			}
		}()
	} else {
		scheme = "http"
		go func() {
			if err := app.Start(); err != nil {
				zlog.Fatal().Err(err).Msg("server failed")
			}
		}()
	}

	zlog.Info().Str("addr", serverAddr).Msg("server started")

	apadanaClient := apadana.New(
		cfg.Cluster.Server,
		cfg.Cluster.Token,
		time.Second*5,
	)

	nodeStatus := &corev1.NodeStatus{
		Addresses: cfg.Addresses,
		Capacity:  corev1.NodeCapacity{MaxInbounds: cfg.MaxInbounds},
		ConnectionConfig: corev1.NodeConnectionConfig{
			Scheme: scheme,
			Port:   cfg.Port,
		},
		Ready: true,
	}

	hb := health.NewHeartbeatManager(
		apadanaClient,
		cfg.NodeStatusUpdateFrequency,
		nodeStatus,
	)

	syncManager := satrapSyncManager.NewSyncManager(
		xrayClient,
		apadanaClient,
		cfg.SyncFrequency,
		cfg.ConcurrentInboundSyncs,
		cfg.ConcurrentInboundExpireSyncs,
		cfg.ConcurrentInboundGCSyncs,
		cfg.ConcurrentUserSyncs,
		cfg.ConcurrentUserExpireSyncs,
	)

	if cfg.Cluster.Enabled {
		hlock := flock.NewFlock("/tmp/satrap-heartbeat.lock")
		if err := hlock.TryLock(); err == nil {
			go hb.Run(ctx, cfg.NodeID)
			defer hlock.Unlock()
		}

		slock := flock.NewFlock("/tmp/satrap-sync-manager.lock")
		if err := slock.TryLock(); err == nil {
			go syncManager.Run(ctx, cfg.NodeID)
			defer slock.Unlock()
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	defer func() {
		zlog.Info().Str("addr", serverAddr).Msg("shutting down server")
		if err := app.Shutdown(shutdownCtx); err != nil {
			zlog.Error().Err(err).Str("addr", serverAddr).Msg("server shutdown error")
		}
		zlog.Info().Msg("shutdown complete")
	}()

	<-ctx.Done()
}
