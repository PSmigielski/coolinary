// Package events applies other services' events to planner's database (docs/EPICs/README.md "Events",
// docs/service-conventions.md "Events").
//
// The Consumer (PLN-01-01) reads the streams in EVENT_STREAMS as consumer group planner. Each registered Handler runs
// in one transaction with the record that makes a second delivery of the same event a no-op.
//
//	consumer.Handle("household.plan_changed", 1, func(ctx context.Context, q *store.Queries, e events.Envelope) error {
//		... decode e.Payload, write with q ...
//	})
//
// The consumer is a copy of kitchen's (KIT-01-04) with the same API, envelope type and metric names; it differs where
// PLN-01-01's rules say so (Q232). Planner publishes no events (PLN-01), so kitchen's outbox, relay and event types
// aren't copied; a later epic that needs to publish copies them.
package events
