package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vayzur/apadana/internal/chapar/server"
	"github.com/vayzur/apadana/internal/config"
	chaparconfigv1 "github.com/vayzur/apadana/pkg/chapar/config/v1"
	"github.com/vayzur/apadana/pkg/chapar/service"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/pkg/chapar/storage/etcd"
	"github.com/vayzur/apadana/pkg/chapar/storage/resources"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := &chaparconfigv1.ChaparConfig{}
	if err := config.Load(*configPath, cfg); err != nil {
		zlog.Fatal().
			Err(err).
			Str("component", "config").
			Str("path", *configPath).
			Msg("failed to load configuration")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	etcdClient, err := etcd.NewClient(&cfg.Etcd, ctx)
	if err != nil {
		zlog.Fatal().
			Err(err).
			Str("component", "etcd").
			Msg("failed to connect")
	}
	defer func() {
		zlog.Info().
			Str("component", "etcd").
			Msg("closing client")
		if err := etcdClient.Close(); err != nil {
			zlog.Error().
				Err(err).
				Str("component", "etcd").
				Msg("client close error")
		}
	}()

	etcdStorage := etcd.NewEtcdStorage(etcdClient)
	if err := etcdStorage.ReadinessCheck(); err != nil {
		zlog.Error().
			Err(err).
			Str("component", "etcd").
			Msg("readiness check failed: not ready")
	}

	inboundStore := resources.NewInboundStore(etcdStorage)
	nodeStore := resources.NewNodeStore(etcdStorage)
	nodeService := service.NewNodeService(nodeStore)
	inboundService := service.NewInboundService(inboundStore)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	app := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, inboundService, nodeService)

	go func() {
		var err error
		if cfg.TLS.Enabled {
			err = app.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			err = app.Start()
		}
		if err != nil {
			zlog.Fatal().
				Err(err).
				Str("component", "server").
				Str("addr", serverAddr).
				Msg("failed to start")
		}
	}()

	zlog.Info().
		Str("component", "server").
		Str("addr", serverAddr).
		Msg("started successfully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	defer func() {
		zlog.Info().
			Str("component", "server").
			Str("addr", serverAddr).
			Msg("shutting down")
		if err := app.Shutdown(shutdownCtx); err != nil {
			zlog.Error().
				Err(err).
				Str("component", "server").
				Str("addr", serverAddr).
				Msg("shutdown error")
		}
		zlog.Info().
			Str("component", "server").
			Msg("shutdown complete")
	}()

	<-ctx.Done()
}
