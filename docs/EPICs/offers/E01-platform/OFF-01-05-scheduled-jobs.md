# OFF-01-05 Scheduled job runner

| | |
|---|---|
| Service | offers |
| Epic | [OFF-01 Platform](README.md) |
| Release | R2 |
| Size | M |
| Depends on | — |
| Unblocks | OFF-02-03, OFF-03-01, OFF-05-04, OFF-05-05, OFF-07-02, OFF-07-07, OFF-08-03, OFF-11-03 |

## Goal

offers runs its recurring work (OSM refresh, flyer checks, batch polling, publishing, expiry, spot-check sampling,
clean-ups) on a schedule, in the chains' local time, with no job lost or run twice when the container restarts.

## Business rules

1. Jobs are registered in code with a name, a schedule and a handler. Schedules are cron expressions evaluated in a
   named time zone (`Europe/Warsaw` for PL jobs, `Europe/Berlin` for DE jobs, `UTC` otherwise).
2. A tick every 30 s finds due jobs. Each run is claimed with `SELECT … FOR UPDATE SKIP LOCKED` on its `job_runs` row,
   so a second process (or a restarted one) never runs the same due slot twice.
3. A missed slot (container down) runs once on start-up if it is less than 6 hours late; older missed slots are
   recorded as `skipped`.
4. One run of a job at a time: if the previous run of the same job is still `running`, the new slot is `skipped`.
5. A run that has been `running` for longer than its `timeout` (default 30 min) is marked `failed` with
   reason `timeout` and its lock released.
6. A run can be triggered by admin (`TriggerJob`, OFF-01-06); it is recorded with `trigger = 'admin'` and the admin id.
7. Each run stores status (`running`, `done`, `failed`, `skipped`), start, end and a short error text. Runs older than
   30 days are deleted.
8. Three failed runs of the same job in a row send a founder alert (`SendAlert`, key `offers.job.<name>`).

## Contract

`contracts/coolinary/offers/v1/admin.proto` (shell in OFF-01-06):

```proto
rpc ListJobs(ListJobsRequest) returns (ListJobsResponse);      // name, schedule, last run, next run
rpc TriggerJob(TriggerJobRequest) returns (TriggerJobResponse); // name -> run id
```

## Data

`offers/migrations/000NN_jobs.sql`

```sql
CREATE TABLE job_runs (
    id          bigserial PRIMARY KEY,
    job         text NOT NULL,
    slot        timestamptz NOT NULL,          -- the scheduled time this run is for
    trigger     text NOT NULL DEFAULT 'schedule' CHECK (trigger IN ('schedule', 'admin', 'startup')),
    admin_user_id text,
    status      text NOT NULL CHECK (status IN ('running', 'done', 'failed', 'skipped')),
    started_at  timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz,
    error       text,
    UNIQUE (job, slot)
);
CREATE INDEX job_runs_recent ON job_runs (job, started_at DESC);
```

## Errors

| Case | Behaviour |
|---|---|
| Unknown job name in `TriggerJob` | `NOT_FOUND` |
| Job already running | `FAILED_PRECONDITION` "job is running" |

## Acceptance criteria

- [ ] A job scheduled `0 6 * * *` in `Europe/Warsaw` runs at 06:00 Warsaw time on both sides of the DST change.
- [ ] Killing the container mid-run leaves the run `failed` with `timeout` after the timeout, and the next slot runs.
- [ ] Restarting the container 2 h after a missed slot runs it once with `trigger = 'startup'`; 7 h after, it is
      `skipped`.
- [ ] Metric `jobs_runs_total{job,status}` and `jobs_last_success_timestamp{job}` exist.
