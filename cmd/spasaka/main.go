package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/vayzur/apadana/internal/config"
	"github.com/vayzur/apadana/pkg/chapar/storage/etcd"
	apadana "github.com/vayzur/apadana/pkg/client"
	"github.com/vayzur/apadana/pkg/leader"
	spasakaconfigv1 "github.com/vayzur/apadana/pkg/spasaka/config/v1"
	"github.com/vayzur/apadana/pkg/spasaka/controller"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	cfg := spasakaconfigv1.SpasakaConfig{}

	if err := config.Load(*configPath, &cfg); err != nil {
		zlog.Fatal().Err(err).Msg("failed to load config")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
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

	etcdSession, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(10), concurrency.WithContext(ctx))
	if err != nil {
		zlog.Error().Err(err).Msg("etcd new session failed")
	}

	defer func() {
		zlog.Info().Msg("closing etcd session")
		if err := etcdSession.Close(); err != nil {
			zlog.Error().Err(err).Msg("etcd session close error")
		}
	}()

	apadanaClient := apadana.New(cfg.Cluster.Server, cfg.Cluster.Token, time.Second*5)
	spasakaManager := controller.NewSpasaka(apadanaClient)

	val := "spasaka"

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := leader.Run(ctx, etcdSession, "/lock/node-controller", val, func(leaderCtx context.Context) {
			spasakaManager.RunNodeMonitor(leaderCtx, cfg.ConcurrentNodeSyncs, cfg.NodeMonitorPeriod, cfg.NodeMonitorGracePeriod)
		}); err != nil && ctx.Err() == nil {
			zlog.Error().Err(err).Msg("failed to start node controller")
		}
	}()

	<-ctx.Done()
	zlog.Info().Msg("shutting down")
	wg.Wait()
	zlog.Info().Msg("shutdown complete")
}
