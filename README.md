# room-decorator

A small Go service that accepts jobs over an HTTP API and processes them
asynchronously on a background worker. Jobs are persisted in Postgres
(Supabase); the queue is still in-memory, so jobs already enqueued at the
moment of a restart are dropped, but jobs themselves survive.

## Layout

```
cmd/server/        HTTP server + worker, the single binary
internal/api/      HTTP handlers and routing
internal/core/     business logic (JobRepo interface, job service, worker loop)
internal/infra/    Postgres job repository and in-memory queue
internal/models/   data types and domain errors
migrations/        golang-migrate SQL files
```

## Setup

### Prerequisites

- Go (matches `go.mod`)
- A Supabase Postgres project
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI for
  applying schema migrations
- [`direnv`](https://direnv.net/) (recommended) for auto-loading
  `DATABASE_URL` when you `cd` into the repo

### Configuring `DATABASE_URL`

The server and the migrate CLI both read `DATABASE_URL` from the environment.
For local development, use the Supabase **session pooler** connection string:

```
postgresql://postgres.<PROJECT_REF>:<PASSWORD>@aws-0-<REGION>.pooler.supabase.com:5432/postgres?sslmode=require
```

Put it in `.env` (gitignored). With direnv set up, an `.envrc` containing
`dotenv` will load it automatically. Otherwise, `export` it in your shell
before running anything.

### Applying migrations

```
migrate -path migrations -database "$DATABASE_URL" up
```

## Running

```
go run ./cmd/server
```

The server listens on `:8080`. On startup it builds a `pgxpool` connection
pool and pings the database — a missing or invalid `DATABASE_URL` causes a
fast, loud exit. The worker runs in a background goroutine inside the same
process and shares the repo + queue with the HTTP handlers.

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

## Migrations workflow

Create a new migration:

```
migrate create -ext sql -dir migrations -seq <descriptive_name>
```

Edit the generated `*.up.sql` and `*.down.sql`, then apply:

```
migrate -path migrations -database "$DATABASE_URL" up
```

Other useful commands:

```
migrate -path migrations -database "$DATABASE_URL" version   # current version
migrate -path migrations -database "$DATABASE_URL" down 1    # roll back one
migrate -path migrations -database "$DATABASE_URL" force <v> # recover from dirty state
```

## Tests

```
go test ./...
```

Currently only the in-memory queue is covered. Postgres-backed repository
tests will need a real (or testcontainers-managed) database — not yet wired
up.
