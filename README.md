# room-decorator

A small Go service that accepts jobs over an HTTP API and processes them
asynchronously on a background worker. Built as a learning project — storage
and queueing are in-memory, so all state is lost on restart.

## Layout

```
cmd/server/        HTTP server + worker, the single binary
internal/api/      HTTP handlers and routing
internal/core/     business logic (job service, worker loop)
internal/infra/    in-memory job repository and queue
internal/models/   data types
```

## Running

```
go run ./cmd/server
```

The server listens on `:8080`. The worker runs in a background goroutine
inside the same process, so it shares the in-memory repo and queue with the
HTTP handlers.

## API

### Create a job

```
curl -X POST localhost:8080/jobs \
  -H 'Content-Type: application/json' \
  -d '{"payload":"hello"}'
```

Returns `201 Created` with the new job. The job starts in `PENDING`,
transitions to `PROCESSING` when the worker picks it up, and ends in
`SUCCESS` (or `FAILED`) once `ProcessJob` returns.

### Get a job

```
curl localhost:8080/jobs/<id>
```

Returns `200 OK` with the job, or `404 Not Found` if the id is unknown.

## Tests

```
go test ./...
```
