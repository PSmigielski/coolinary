package clients

import (
	"maps"
	"regexp"
	"testing"

	"connectrpc.com/connect"

	kitchenv1 "coolinary/planner/gen/coolinary/kitchen/v1"
	"coolinary/planner/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/planner/internal/platform/rpc"
)

var jobRequestID = regexp.MustCompile(`^job-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func TestJobContext_SendsAFreshJobRequestIDAndTheHouseholdsPlanOnly(t *testing.T) {
	tests := []struct {
		name string
		plan rpc.PlanCode
	}{
		{name: "plus household", plan: rpc.PlanPlus},
		{name: "free household", plan: rpc.PlanFree},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fake{}
			c := New(Addrs{Kitchen: startFake(t, kitchenv1connect.KitchenServiceResolveBarcodeProcedure, f)})
			// A job started from a request must not send that request's user, household or id on.
			parent := incoming(t.Context(), gatewayMetadata)

			_, err := c.Kitchen.ResolveBarcode(JobContext(parent, tt.plan),
				connect.NewRequest(&kitchenv1.ResolveBarcodeRequest{}))

			if err != nil {
				t.Fatalf("call: %v", err)
			}
			got := f.calls()[0].metadata
			if !jobRequestID.MatchString(got["x-request-id"]) {
				t.Errorf("x-request-id %q, want job-<uuid>", got["x-request-id"])
			}
			want := map[string]string{"x-request-id": got["x-request-id"], "x-plan": string(tt.plan)}
			if !maps.Equal(got, want) {
				t.Errorf("kitchen received %v, want %v", got, want)
			}
		})
	}
}

func TestJobContext_EachJobContextGetsItsOwnRequestID(t *testing.T) {
	f := &fake{}
	c := New(Addrs{Kitchen: startFake(t, kitchenv1connect.KitchenServiceResolveBarcodeProcedure, f)})

	for range 2 {
		_, err := c.Kitchen.ResolveBarcode(JobContext(t.Context(), rpc.PlanFree),
			connect.NewRequest(&kitchenv1.ResolveBarcodeRequest{}))
		if err != nil {
			t.Fatalf("call: %v", err)
		}
	}

	calls := f.calls()
	first, second := calls[0].metadata["x-request-id"], calls[1].metadata["x-request-id"]
	if first == "" || first == second {
		t.Errorf("request ids %q and %q, want two different ones", first, second)
	}
}
