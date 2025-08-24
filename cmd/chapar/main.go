package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := chaparconfigv1.ChaparConfig{}

	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Str("component", "config").Str("action", "load").Msg("failed")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	etcdCtx, etcdCancel := context.WithCancel(ctx)
	defer etcdCancel()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
		Context:     etcdCtx,
	})
	if err != nil {
		zlog.Fatal().Err(err).Str("etcd", "client").Str("action", "connect").Msg("failed")
	}
	defer func() {
		if err := etcdClient.Close(); err != nil {
			zlog.Fatal().Err(err).Str("etcd", "client").Str("action", "close").Msg("failed")
		}
	}()

	if err := etcd.CheckEtcdHealth(ctx, etcdClient); err != nil {
		zlog.Fatal().Err(err).Str("etcd", "client").Str("action", "health").Msg("failed")
	}

	etcdStorege := etcd.NewEtcdStorage(etcdClient)

	inboundStore := resources.NewInboundStore(etcdStorege)
	nodeStore := resources.NewNodeStore(etcdStorege)

	satrapClient := satrap.New(time.Second * 5)

	nodeService := service.NewNodeSerivce(nodeStore)
	inboundService := service.NewInboundService(inboundStore, nodeService, satrapClient)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	chapar := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, inboundService, nodeService)

	go func() {
		if cfg.TLS.Enabled {
			if err := chapar.StartTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile); err != nil {
				zlog.Fatal().Err(err).Str("component", "chapar").Str("action", "start").Msg("failed")
			}

		} else {
			if err := chapar.Start(); err != nil {
				zlog.Fatal().Err(err).Str("component", "chapar").Str("action", "start").Msg("failed")
			}
		}
	}()

	defer func() {
		if err := chapar.Stop(); err != nil {
			zlog.Fatal().Err(err).Str("component", "chapar").Str("action", "stop").Msg("failed")
		}
	}()

	zlog.Info().Str("component", "chapar").Str("action", "start").Msg("success")
	<-ctx.Done()
	zlog.Info().Str("component", "chapar").Str("action", "stop").Msg("success")
}
