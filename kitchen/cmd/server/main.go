// Command server runs the kitchen service. `server healthcheck` probes /healthz
// (used as the Docker HEALTHCHECK in the distroless prod image).
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"coolinary/kitchen/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/kitchen/internal/config"
	"coolinary/kitchen/internal/kitchen"
	"coolinary/kitchen/internal/kitchen/store"
	"coolinary/kitchen/internal/platform/admin"
	"coolinary/kitchen/internal/platform/cache"
	"coolinary/kitchen/internal/platform/db"
	"coolinary/kitchen/internal/platform/events"
	"coolinary/kitchen/internal/platform/log"
	"coolinary/kitchen/internal/platform/rpc"
	"coolinary/kitchen/migrations"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		port := 8081
		if p, err := strconv.Atoi(os.Getenv("ADMIN_PORT")); err == nil {
			port = p
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

	svc := kitchen.New(store.New(pool), redisCache, logger)
	grpcSrv := rpc.NewServer(cfg.GRPCPort, cfg.GRPCReflection, logger)
	path, handler := kitchenv1connect.NewKitchenServiceHandler(svc, rpc.Options(logger)...)
	grpcSrv.Handle(kitchenv1connect.KitchenServiceName, path, handler)
	backOffice := kitchen.NewAdmin(pool, logger)
	path, handler = kitchenv1connect.NewKitchenAdminServiceHandler(backOffice, rpc.AdminOptions(logger)...)
	grpcSrv.Handle(kitchenv1connect.KitchenAdminServiceName, path, handler)

	adminSrv := admin.NewServer(cfg.AdminPort, []admin.Check{
		{Name: "postgres", Fn: pool.Ping},
		{Name: "redis_cache", Fn: func(ctx context.Context) error { return redisCache.Ping(ctx).Err() }},
		{Name: "redis_state", Fn: func(ctx context.Context) error { return redisState.Ping(ctx).Err() }},
	}, logger)
	relay := events.NewRelay(pool, redisState, logger, events.RelayOptions{})
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("hostname for the event consumer: %w", err)
	}
	// Handlers are registered by the tasks that need them (consumer.Handle before Run).
	consumer := events.NewConsumer(pool, redisState, logger, events.ConsumerOptions{
		Streams: cfg.EventStreams, Group: cfg.ServiceName, Name: hostname,
	})

	g, gctx := errgroup.WithContext(ctx)
	g.Go(grpcSrv.Serve)
	g.Go(adminSrv.Serve)
	g.Go(func() error { return relay.Run(gctx) })
	g.Go(func() error { return consumer.Run(gctx) })
	g.Go(func() error { backOffice.RunCleanup(gctx); return nil })
	g.Go(func() error {
		<-gctx.Done()
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(gctx), 10*time.Second)
		defer cancel()
		_ = grpcSrv.Shutdown(shutdownCtx)
		return adminSrv.Shutdown(shutdownCtx)
	})
	return g.Wait()
}
