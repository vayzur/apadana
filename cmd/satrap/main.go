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

	xrayAddr := fmt.Sprintf("%s:%d", cfg.Xray.Address, cfg.Xray.Port)
	xrayClient, err := xray.New(xrayAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to connect xray")
	}
	defer func() {
		if err := xrayClient.Close(); err != nil {
			zlog.Error().Err(err).Msg("xray client close failed")
		}
	}()

	apadanaClient := apadana.New(
		cfg.Cluster.Server,
		cfg.Cluster.Token,
		time.Second*5,
	)

	hb := health.NewHeartbeatManager(
		apadanaClient,
		cfg.NodeStatusUpdateFrequency,
		cfg.MaxInbounds,
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

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	app := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, xrayClient)

	go func() {
		var err error
		if cfg.TLS.Enabled {
			err = app.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			err = app.Start()
		}
		if err != nil {
			zlog.Fatal().Err(err).Msg("server failed")
		}
	}()

	zlog.Info().Str("addr", serverAddr).Msg("server started")

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
