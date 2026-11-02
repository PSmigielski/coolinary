package clients

import (
	"context"
	"uuid"

	"coolinary/planner/internal/platform/rpc"
)

// JobContext is the context a background job calls other services with for one household: a fresh x-request-id
// (job-<uuid>) and the household's plan as identity reported it (PLN-01-05). Metadata of a request in ctx is dropped.
func JobContext(ctx context.Context, plan rpc.PlanCode) context.Context {
	md := map[string]string{"x-request-id": "job-" + uuid.New().String(), "x-plan": string(plan)}
	return rpc.WithMetadata(ctx, func(k string) string { return md[k] })
}
