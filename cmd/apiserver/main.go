package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/vayzur/apadana/internal/apiserver/server"
	"github.com/vayzur/apadana/internal/config"
	chapar "github.com/vayzur/apadana/pkg/chapar/client"
	"github.com/vayzur/apadana/pkg/httputil"
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

	cfg := config.APIServerConfig{}

	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Msg("config load failed")
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
		zlog.Fatal().Err(err).Msg("etcd connect failed")
	}
	defer func() {
		if err := etcdClient.Close(); err != nil {
			zlog.Error().Err(err).Msg("failed to close etcd client")
		}
	}()

	if err := etcd.CheckEtcdHealth(ctx, etcdClient); err != nil {
		zlog.Fatal().Err(err).Msg("etcd is not healthy")
	}

	etcdStorege := etcd.NewEtcdStorage(etcdClient)

	inboundStore := resources.NewInboundStore(etcdStorege)
	nodeStore := resources.NewNodeStore(etcdStorege)

	httpClient := httputil.New(time.Second * 5)
	chaparClient := chapar.New(httpClient)

	inboundService := service.NewInboundService(inboundStore, chaparClient)
	nodeService := service.NewNodeSerivce(nodeStore)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)

	apiserver := server.NewServer(serverAddr, cfg.Token, cfg.Prefork, inboundService, nodeService)

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

	zlog.Info().Str("component", "apiserver").Msg("apiserver started")
	defer zlog.Info().Str("component", "apiserver").Msg("apiserver stopped")

	<-ctx.Done()
	zlog.Info().Str("component", "apiserver").Msg("shutting down gracefully...")
}
