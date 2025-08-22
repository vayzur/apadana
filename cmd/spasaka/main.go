package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/vayzur/apadana/internal/config"
	"github.com/vayzur/apadana/pkg/leader"
	"github.com/vayzur/apadana/pkg/service"
	spasakaconfigv1 "github.com/vayzur/apadana/pkg/spasaka/config/v1"
	"github.com/vayzur/apadana/pkg/spasaka/controller"

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

	cfg := spasakaconfigv1.SpasakaConfig{}

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

	etcdStorage := etcd.NewEtcdStorage(etcdClient)

	sessionCtx, sessionCancel := context.WithCancel(ctx)
	defer sessionCancel()

	etcdSession, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(10), concurrency.WithContext(sessionCtx))
	if err != nil {
		zlog.Fatal().Err(err).Str("etcd", "session").Str("action", "create").Msg("failed")
	}

	defer func() {
		if err := etcdSession.Close(); err != nil {
			zlog.Fatal().Err(err).Str("etcd", "session").Str("action", "close").Msg("failed")
		}
	}()

	nodeStore := resources.NewNodeStore(etcdStorage)
	nodeService := service.NewNodeSerivce(nodeStore)

	spasakaManager := controller.NewSpasaka(nodeService)

	val := "spasaka"

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := leader.Run(ctx, etcdSession, "/lock/node-monitor", val, func(leaderCtx context.Context) {
			spasakaManager.RunNodeMonitor(leaderCtx, cfg.NodeMonitorPeriod, cfg.NodeMonitorGracePeriod)
		}); err != nil && ctx.Err() == nil {
			zlog.Error().Err(err).Str("spasaka", "leader").Str("resource", "node").Str("action", "run").Msg("failed")
		}
	}()

	zlog.Info().Str("component", "spasaka").Str("action", "start").Msg("success")

	<-ctx.Done()

	wg.Wait()
	zlog.Info().Str("component", "spasaka").Str("action", "stop").Msg("success")
}
