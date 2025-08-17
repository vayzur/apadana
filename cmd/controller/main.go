package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/vayzur/apadana/internal/config"
	chapar "github.com/vayzur/apadana/pkg/chapar/client"
	"github.com/vayzur/apadana/pkg/controller"
	"github.com/vayzur/apadana/pkg/httputil"
	"github.com/vayzur/apadana/pkg/leader"
	"github.com/vayzur/apadana/pkg/service"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/apadana/pkg/storage/etcd"
	"github.com/vayzur/apadana/pkg/storage/resources"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := config.ControllerConfig{}

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

	etcdStorage := etcd.NewEtcdStorage(etcdClient)

	sessionCtx, sessionCancel := context.WithCancel(ctx)
	defer sessionCancel()

	etcdSession, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(10), concurrency.WithContext(sessionCtx))
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to create etcd session")
	}

	defer func() {
		if err := etcdSession.Close(); err != nil {
			zlog.Error().Err(err).Msg("failed to close etcd session")
		}
	}()

	inboundStore := resources.NewInboundStore(etcdStorage)
	nodeStore := resources.NewNodeStore(etcdStorage)

	httpClient := httputil.New(time.Second * 5)
	chaparClient := chapar.New(httpClient)

	inboundService := service.NewInboundService(inboundStore, chaparClient)
	nodeService := service.NewNodeSerivce(nodeStore)

	controllerManager := controller.NewControllerManager(nodeService, inboundService)

	val := "controller"

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := leader.Run(ctx, etcdSession, "/lock/node-controller", val, func(leaderCtx context.Context) {
			controllerManager.RunNodeMonitor(leaderCtx, cfg.NodeMonitorPeriod, cfg.NodeMonitorGracePeriod)
		}); err != nil && ctx.Err() == nil {
			zlog.Error().Err(err).Msg("node controller leadership failed")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := leader.Run(ctx, etcdSession, "/lock/inbound-controller", val, func(leaderCtx context.Context) {
			controllerManager.RunInboundMonitor(leaderCtx, cfg.InboundMonitorPeriod)
		}); err != nil && ctx.Err() == nil {
			zlog.Error().Err(err).Msg("inbound controller leadership failed")
		}
	}()

	zlog.Info().Str("component", "controller").Msg("controller manager started")

	<-ctx.Done()

	wg.Wait()
	zlog.Info().Str("component", "controller").Msg("shutting down gracefully...")
}
