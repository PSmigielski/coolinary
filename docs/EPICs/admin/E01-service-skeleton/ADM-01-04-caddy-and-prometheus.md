# ADM-01-04 Caddy route and Prometheus scrape

| | |
|---|---|
| Service | admin |
| Epic | [ADM-01 Service skeleton](README.md) |
| Release | R1 |
| Size | S |
| Status | todo |
| Depends on | ADM-01-02 |
| Unblocks | ADM-02-01 |

## Goal

In production admin is reachable at `admin.<domain>` over HTTPS (decision 9), never indexed or framed, and its
metrics land in Grafana with the other services.

## Business rules

1. `platform/caddy/Caddyfile`: block `admin.{$DOMAIN}` with `import common`, `reverse_proxy admin:3000`, and extra
   headers: `X-Robots-Tag "noindex, nofollow"`, `Content-Security-Policy "frame-ancestors 'none'"`,
   `X-Frame-Options "DENY"`, `Cache-Control "no-store"` on `/admin/*`.
2. Request body limit 10 MB for `admin.{$DOMAIN}` (import files are at most 2 MB; flyer pages are only downloaded).
3. DNS: `admin.<domain>` must point to the VPS; the root README production section lists it next to `api.`,
   `grafana.` and `langfuse.`.
4. `platform/prometheus/prometheus.yml`: target `admin:8081` with labels `{ service: admin, runtime: bun }`.
5. The Grafana "Coolinary Services" dashboard shows admin without edits (it groups by `service`); if a panel lists
   services explicitly, add admin.

## Contract

None.

## Data

None.

## Errors

None.

## Acceptance criteria

- [ ] Local prod smoke test: `https://admin.coolinary.localhost/healthz` is not exposed (only port 3000 routes are
      proxied; `/healthz` lives on 8081), `https://admin.coolinary.localhost/openapi` answers.
- [ ] Response headers contain `X-Robots-Tag` and `frame-ancestors 'none'`.
- [ ] Prometheus shows target `admin` up.
