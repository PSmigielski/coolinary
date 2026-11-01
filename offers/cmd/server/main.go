// Command server runs the offers service. `server healthcheck` probes /healthz
// (used as the Docker HEALTHCHECK in the distroless prod image).
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"coolinary/offers/gen/coolinary/offers/v1/offersv1connect"
	"coolinary/offers/internal/config"
	"coolinary/offers/internal/offers"
	"coolinary/offers/internal/offers/store"
	"coolinary/offers/internal/platform/admin"
	"coolinary/offers/internal/platform/cache"
	"coolinary/offers/internal/platform/db"
	"coolinary/offers/internal/platform/log"
	"coolinary/offers/internal/platform/rpc"
	"coolinary/offers/migrations"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		port := 8081
		if p := os.Getenv("ADMIN_PORT"); p != "" {
			fmt.Sscanf(p, "%d", &port)
		}
		if err := admin.Probe(port); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	logger := log.New(cfg.ServiceName, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool, migrations.FS, logger); err != nil {
		return err
	}
	redisCache, err := cache.Open(ctx, cfg.RedisCacheURL)
	if err != nil {
		return err
	}
	redisState, err := cache.Open(ctx, cfg.RedisStateURL)
	if err != nil {
		return err
	}

	svc := offers.New(store.New(pool), logger)
	grpcSrv := rpc.NewServer(cfg.GRPCPort, cfg.GRPCReflection, logger)
	path, handler := offersv1connect.NewOffersServiceHandler(svc, rpc.Options(logger)...)
	grpcSrv.Handle(offersv1connect.OffersServiceName, path, handler)

	adminSrv := admin.NewServer(cfg.AdminPort, []admin.Check{
		{Name: "postgres", Fn: pool.Ping},
		{Name: "redis_cache", Fn: func(ctx context.Context) error { return redisCache.Ping(ctx).Err() }},
		{Name: "redis_state", Fn: func(ctx context.Context) error { return redisState.Ping(ctx).Err() }},
	}, logger)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(grpcSrv.Serve)
	g.Go(adminSrv.Serve)
	g.Go(func() error {
		<-gctx.Done()
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = grpcSrv.Shutdown(shutdownCtx)
		return adminSrv.Shutdown(shutdownCtx)
	})
	return g.Wait()
}
