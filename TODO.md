# TODO

Future work items, in no particular order.

## Graceful worker shutdown

`http.Server.Shutdown` (if/when added) only stops the HTTP server — it does
not tell the worker goroutine to stop processing. When `main` returns, the
worker is killed abruptly and may be cut off mid-job.

Wiring graceful worker shutdown means:

- Give `core.RunWorker` a `context.Context` parameter.
- Have the loop `select` between `queue.Dequeue()` and `<-ctx.Done()` so it
  can exit cleanly when the context is cancelled.
- Cancel that context from `main` after `httpServer.Shutdown` returns, and
  wait on a `done` channel before exiting so any in-flight job finishes.

Not required to get the API running, but worth doing once jobs do real
work or persistence is added.

## Real `ProcessJob` implementation

`core.ProcessJob` currently just sleeps for 500ms and returns nil. Replace
with the actual work this service is meant to do (room decoration / image
generation / whatever the product turns out to be). The `Payload` field on
`Job` will probably need to become structured rather than a plain string
once this happens.

## Persistence

Repo and queue are both in-memory, so a restart wipes all jobs. Swap the
in-memory repo for a real database (SQLite is the easiest first step;
Postgres if you want to learn the production-ish stack). Once persistence
exists, several other items in this file become more relevant (graceful
shutdown, retry, backpressure).

## Failed-job retry

If `ProcessJob` returns an error, the job is marked `FAILED` and never
retried. Add retry-with-backoff:

- New fields on `Job`: `Attempts int`, `MaxAttempts int`, maybe `LastError string`.
- On failure, if `Attempts < MaxAttempts`, increment and re-enqueue (with a
  delay if the queue grows up to support that) instead of marking `FAILED`.
- Only mark `FAILED` once attempts are exhausted.

## Queue backpressure

`InMemoryQueue` has a buffer of 10. If the queue fills up, `core.CreateJob`
blocks the HTTP handler indefinitely waiting for `Enqueue`. Two reasonable
fixes:

- Use a `select` with a `default` branch in `Enqueue` and return an error
  when full; have the handler respond with `503 Service Unavailable`.
- Or just make the buffer large enough that this won't happen in practice.

The first option is more correct, the second is simpler. Pick consciously.

## HTTP server timeouts

`http.Server` defaults to no read/write/idle timeouts, which lets a slow
client tie up a connection forever (slow-loris attack). Set
`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` on
the `http.Server` in `cmd/server/main.go`. ~4 lines.

## Repo and queue as interfaces

`internal/api` and `internal/core` currently depend on the concrete
`*infra.InMemoryJobRepo` and `*infra.InMemoryQueue`. Define small interfaces
in `core` (e.g. a `JobRepo` with `Get`, `Insert`, `UpdateStatus`, and a
`JobQueue` with `Enqueue`, `Dequeue`) and have the infra types satisfy them
implicitly.

Benefits:

- Lets you swap the in-memory implementation for a Postgres or Redis one
  without touching `internal/api` or `internal/core`.
- Lets you write handler/service tests with fakes rather than real infra.

## Tests for `core` and `api`

`internal/infra` has tests; `core` and `api` don't. Once the interfaces
above exist, handler tests with `httptest.NewRecorder` and service tests
with fake repo/queue become straightforward.

## Migrate `log` → `log/slog`

Go 1.21+ ships structured logging in the stdlib. Replace `log.Printf` calls
with `slog.Info` / `slog.Error` and pass key/value pairs (e.g.
`slog.Info("job processed", "job_id", id, "status", status)`). Pays off
once logs grow large enough that you want to filter by field.
