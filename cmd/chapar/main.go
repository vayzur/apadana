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
	satrap "github.com/vayzur/apadana/pkg/satrap/client"
	"github.com/vayzur/apadana/pkg/service"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/pkg/storage/etcd"
	"github.com/vayzur/apadana/pkg/storage/resources"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := chaparconfigv1.ChaparConfig{}
	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	etcdClient, err := etcd.NewClient(&cfg.Etcd, ctx)
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to connect etcd")
	}

	defer func() {
		zlog.Info().Msg("closing etcd client")
		if err := etcdClient.Close(); err != nil {
			zlog.Error().Err(err).Msg("etcd client close error")
		}
	}()

	etcdStorage := etcd.NewEtcdStorage(etcdClient)

	if err := etcdStorage.ReadinessCheck(); err != nil {
		zlog.Fatal().Err(err).Msg("etcd not ready")
	}

	inboundStore := resources.NewInboundStore(etcdStorage)
	nodeStore := resources.NewNodeStore(etcdStorage)
	satrapClient := satrap.New(time.Second * 5)
	nodeService := service.NewNodeService(nodeStore)
	inboundService := service.NewInboundService(inboundStore, nodeService, satrapClient)

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
