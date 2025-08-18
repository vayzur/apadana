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
	apadana "github.com/vayzur/apadana/pkg/client"
	"github.com/vayzur/apadana/pkg/httputil"
	"github.com/vayzur/apadana/pkg/satrap/flock"
	"github.com/vayzur/apadana/pkg/satrap/health"
	xray "github.com/vayzur/apadana/pkg/satrap/xray/client"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := config.SatrapConfig{}

	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Msg("config load failed")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	xrayAddr := fmt.Sprintf("%s:%d", cfg.Xray.Address, cfg.Xray.Port)
	xrayClient, err := xray.New(xrayAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("xray connect failed")
	}

	defer xrayClient.Close()

	httpClient := httputil.New(time.Second * 5)
	apadanaClient := apadana.New(
		httpClient,
		cfg.Cluster.Server,
		cfg.Cluster.Token,
	)

	hb := health.NewHeartbeatManager(
		apadanaClient,
		cfg.NodeStatusUpdateFrequency,
	)

	if cfg.Cluster.Enabled {
		lock := flock.NewFlock("/tmp/chapar-heartbeat.lock")

		if err := lock.TryLock(); err == nil {
			go hb.StartHeartbeat(cfg.NodeID, ctx)
			defer lock.Unlock()
		}
	}

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	apiserver := server.NewServer(serverAddr, cfg.Token, xrayClient)

	go func() {
		if cfg.TLS.Enabled {
			zlog.Fatal().Err(apiserver.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile))
		} else {
			zlog.Fatal().Err(apiserver.Start())
		}
	}()

	defer func() {
		if err := apiserver.Stop(); err != nil {
			zlog.Error().Err(err).Msg("failed to stop apiserver")
		}
	}()

	zlog.Info().Str("component", "chapar").Msg("server started")
	<-ctx.Done()
	zlog.Info().Str("component", "chapar").Msg("shutting down gracefully")
}
