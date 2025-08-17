package distlock

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func RunAsLeader(ctx context.Context, session *concurrency.Session, key string, callback func(context.Context)) error {
	election := concurrency.NewElection(session, key)

	zlog.Info().Str("key", key).Msg("trying to become leader...")

	// Try to become leader (blocks until we get it)
	if err := election.Campaign(ctx, "controller"); err != nil {
		return err
	}

	zlog.Info().Str("key", key).Msg("I AM THE LEADER! Starting work...")

	go callback(ctx)

	select {
	case <-ctx.Done():
		zlog.Info().Str("key", key).Msg("context cancelled - stepping down")
	case <-session.Done():
		zlog.Warn().Str("key", key).Msg("etcd session lost - stepping down")
	}

	resignCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	election.Resign(resignCtx)
	cancel()

	zlog.Info().Str("key", key).Msg("stepped down from leadership")
	return nil
}
