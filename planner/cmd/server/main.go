// Command server runs the planner service. `server healthcheck` probes /healthz
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

	"coolinary/planner/gen/coolinary/planner/v1/plannerv1connect"
	"coolinary/planner/internal/config"
	"coolinary/planner/internal/planner"
	"coolinary/planner/internal/planner/store"
	"coolinary/planner/internal/platform/admin"
	"coolinary/planner/internal/platform/cache"
	"coolinary/planner/internal/platform/db"
	"coolinary/planner/internal/platform/log"
	"coolinary/planner/internal/platform/rpc"
	"coolinary/planner/migrations"
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

	svc := planner.New(store.New(pool), logger)
	grpcSrv := rpc.NewServer(cfg.GRPCPort, cfg.GRPCReflection, logger)
	path, handler := plannerv1connect.NewPlannerServiceHandler(svc, rpc.Options(logger)...)
	grpcSrv.Handle(plannerv1connect.PlannerServiceName, path, handler)

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
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(gctx), 10*time.Second)
		defer cancel()
		_ = grpcSrv.Shutdown(shutdownCtx)
		return adminSrv.Shutdown(shutdownCtx)
	})
	return g.Wait()
}
