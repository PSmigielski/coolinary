// Command smoketest checks gRPC calls from this Go service to the Bun notifier:
// a unary call, a deadline, a cancel and (with -idle) a call after an idle gap.
// Run inside the dev container: `make smoke`.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"

	notifierv1 "coolinary/kitchen/gen/coolinary/notifier/v1"
	"coolinary/kitchen/gen/coolinary/notifier/v1/notifierv1connect"
	"coolinary/kitchen/internal/platform/rpc"
)

func main() {
	addr := flag.String("notifier", envOr("NOTIFIER_GRPC_ADDR", "notifier:50051"), "notifier address")
	idle := flag.Duration("idle", 0, "wait this long, then call again on the same connection")
	flag.Parse()

	httpClient := rpc.NewHTTPClient()
	client := notifierv1connect.NewNotifierServiceClient(httpClient, rpc.BaseURL(*addr), rpc.ClientOptions()...)
	health := grpchealth.NewClient(httpClient, rpc.BaseURL(*addr), connect.WithGRPC())

	failed := false
	step := func(name string, fn func() error) {
		start := time.Now()
		if err := fn(); err != nil {
			failed = true
			fmt.Printf("FAIL %-28s %v (%s)\n", name, err, time.Since(start).Round(time.Millisecond))
			return
		}
		fmt.Printf("ok   %-28s (%s)\n", name, time.Since(start).Round(time.Millisecond))
	}

	send := func(ctx context.Context) error {
		_, err := client.SendEmail(ctx, connect.NewRequest(&notifierv1.SendEmailRequest{
			To: "smoke@coolinary.test", Subject: "Go -> Bun smoke test", TextBody: "sent by kitchen/cmd/smoketest",
		}))
		return err
	}

	step("health check", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		res, err := health.Check(ctx, &grpchealth.CheckRequest{})
		if err == nil && res.Status != grpchealth.StatusServing {
			err = fmt.Errorf("status %v", res.Status)
		}
		return err
	})
	step("unary SendEmail", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return send(ctx)
	})
	step("deadline is enforced", func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := send(ctx)
		if connect.CodeOf(err) == connect.CodeDeadlineExceeded {
			return nil
		}
		return fmt.Errorf("want DeadlineExceeded, got %v", err)
	})
	step("cancel is propagated", func() error {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := send(ctx)
		if connect.CodeOf(err) == connect.CodeCanceled || errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("want Canceled, got %v", err)
	})
	if *idle > 0 {
		fmt.Printf("...  idle for %s\n", *idle)
		time.Sleep(*idle)
		step("unary after idle", func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return send(ctx)
		})
	}
	if failed {
		os.Exit(1)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
