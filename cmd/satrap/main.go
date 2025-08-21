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
	"github.com/vayzur/apadana/pkg/satrap/manager/lease"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := satrapconfigv1.Config{}

	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Str("component", "config").Str("action", "load").Msg("failed")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	xrayAddr := fmt.Sprintf("%s:%d", cfg.Xray.Address, cfg.Xray.Port)
	xrayClient, err := xray.New(xrayAddr)
	if err != nil {
		zlog.Fatal().Err(err).Str("client", "xray").Str("action", "connect").Msg("failed")
	}

	defer func() {
		if err := xrayClient.Close(); err != nil {
			zlog.Error().Err(err).Str("client", "xray").Str("action", "stop").Msg("failed")
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
	)

	inboundLeaseManager := lease.NewInboundLeaseManager(
		apadanaClient,
		cfg.InboundTTLCheckPeriod,
	)

	if cfg.Cluster.Enabled {
		hbLock := flock.NewFlock("/tmp/satrap-heartbeat.lock")
		if err := hbLock.TryLock(); err == nil {
			go hb.Run(ctx, cfg.NodeID)
			defer hbLock.Unlock()
		}

		lmLock := flock.NewFlock("/tmp/satrap-inbound-lease-manager.lock")
		if err := lmLock.TryLock(); err == nil {
			go inboundLeaseManager.Run(ctx, cfg.NodeID)
			defer lmLock.Unlock()
		}
	}

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	satrap := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, xrayClient)

	go func() {
		if cfg.TLS.Enabled {
			if err := satrap.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil {
				zlog.Fatal().Err(err).Str("component", "satrap").Str("action", "start").Msg("failed")
			}

		} else {
			if err := satrap.Start(); err != nil {
				zlog.Fatal().Err(err).Str("component", "satrap").Str("action", "start").Msg("failed")
			}
		}
	}()

	defer func() {
		if err := satrap.Stop(); err != nil {
			zlog.Error().Err(err).Str("server", "satrap").Str("action", "stop").Msg("failed")
		}
	}()

	zlog.Info().Str("component", "satrap").Str("action", "start").Msg("success")
	<-ctx.Done()
	zlog.Info().Str("component", "satrap").Str("action", "stop").Msg("success")
}
