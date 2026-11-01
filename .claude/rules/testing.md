# Tests: test first, and what a good test looks like

Applies to every service. What must be tested, where test files go and how tests reach Postgres and Redis is in
`docs/service-conventions.md` "Tests". This file is how to build with tests and how to write them.

## Test-driven development

Build every business rule, acceptance criterion and bug fix red, green, refactor:

1. **Red.** Write one test for the next rule, named after it. Run it and watch it fail **for the reason you
   expect**: an assertion about the missing behaviour, not a compile error, a typo or a missing import. A skipped
   test is not red (tests that need Postgres or Redis skip unless `make dev` is up). A test that passes before the
   code exists tests nothing: fix the test.
2. **Green.** Write the least code that makes it pass. Don't add behaviour no test asks for yet.
3. **Refactor.** With everything green, clean up names, duplication and structure in the code and the tests, then
   run the tests again.

Then the next rule: the simple case first, then edge cases and error paths.

- Contracts and generated code come first (`make proto`, `make sqlc`, `bun run db:generate`), so tests compile
  against the real types.
- A bug fix starts with a test that reproduces the bug and fails, then the fix.
- Never change a test just to make it pass (weaker assertion, deleted case, added skip). If the test is right, fix
  the code. If the task's rule seems wrong, that is a decision: ask, or block the task (`CLAUDE.md`).
- Not test-first: generated code, wiring in `main.go` / `src/index.ts`, settings in config, and a migration on its
  own (the code that uses the table is tested against it). `mobile` has no test runner yet; `make lint s=mobile` is
  its check. Lint, `make smoke` and the running stack cover these.
- Run tests with `make test s=<svc>` (it sets the test database variables). For a quicker loop on tests that need
  nothing running: `~/.bun/bin/bun test <file>` in a Bun service, `uv run pytest tests/<file> -k <name>` in `ai/`.

## What a good test looks like

- **Behaviour, not implementation.** Call the code the way its callers do and check what comes out: the returned
  message, the error code, the rows written, the event emitted. Entry points per service: the unexported method
  behind each RPC (Go, `go.md`), the plain functions the handlers call (Bun gRPC), `createApp().handle(...)`
  (gateway), the `stub` provider and the mapping to the proto response (`ai`, `fastapi.md`). Don't assert on
  private helpers, call order or internal state: a refactor that keeps the behaviour must not break a test.
- **One behaviour per test**, named after the rule so a failure reads as a broken rule
  (`TestConsumeStock_OldestExpiryFirst`, `"rejects a target under 1,200 kcal"`). In Go table tests every row has a
  `name` and tests one case.
- **Arrange, act, assert**, in that order. Keep setup short: a helper that builds valid input and takes only the
  fields the test cares about beats a copied fixture.
- **Exact assertions.** Compare the whole value, not "not null" or "length > 0". For errors check the code (the
  connect or gRPC status, the error reason from `docs/EPICs/README.md`), not just that something failed, and that
  the message carries no personal data.
- **Edges and failures.** Values at, just under and just over each limit; empty lists; missing optional fields;
  duplicates; Free versus Plus; another household's data (must stay invisible). Every way a dependency can fail:
  timeout, unavailable, AI output that breaks a rule (the non-AI fallback is returned).
- **Deterministic.** Nothing in an assertion depends on the current time or randomness: pass in a clock or the
  values. Never `sleep` to wait for something; wait for the condition, with a timeout. Tests don't depend on run
  order or on each other.
- **Fakes at the edges, real code inside.** Replace what is outside the service with small hand-written fakes:
  other services' gRPC clients, the LLM (`stub` provider), Open Food Facts and retailer sites (recorded fixtures),
  Expo Push, the clock. Use the real test Postgres and Redis for queries, transactions and locks instead of
  mocking them; a mocked query proves nothing about the SQL.
- **Own your data.** Each test creates the rows it needs with fresh random ids and removes them afterwards
  (`afterAll`, `t.Cleanup`). Never rely on dev or seed data, or on rows another test made.
- **Fake data that looks fake.** `@example.test` addresses and made-up names; no real personal or health data,
  even in fixtures.
- **Fast.** Unit tests need nothing running and take milliseconds. Use the database only for what needs it.
- **Test your code, not your tools.** Don't test generated code, the framework or a library, only how you use
  them. Where a type should make wrong code impossible, prove it at compile time (`@ts-expect-error` checked by
  `tsc`, as in `identity/src/events/outbox.test.ts`; `var _ Interface = (*T)(nil)` in Go).
- **A flaky test is a bug.** Find the race or the hidden dependency and fix it; never retry or skip to get green.
- **Coverage is not the goal.** Each rule needs a test that fails when the rule breaks; code without behaviour
  doesn't need one.
