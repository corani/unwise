# Copilot / AI Agent Instructions for unwise

This file gives focused, actionable guidance for AI coding agents working on the unwise repository (a small Go HTTP service that implements a subset of the Readwise API).

Summary
- Purpose: small Go service exposing a Readwise-compatible API used by Moon+ Reader and Obsidian.
- Key components: `cmd/unwise` (main entrypoint), `internal/config` (env/config loader), `internal/storage` (storage interface), `internal/storage/sqlite` (SQLite implementation), and `internal/web` (HTTP handlers, routing, and types).

Architecture & data flow
- Startup: `cmd/unwise/main.go` calls `config.MustLoad()`, creates a storage backend (`sqlite.New`) and constructs a web `Server` via `web.New`.
- HTTP: `internal/web/app.go` builds the Fiber app and mounts routes. Handlers are methods on `web.Server` in `internal/web/handlers.go`.
- Storage abstraction: `internal/storage.Storage` interface defines `AddBook`, `AddHighlight`, `UpdateHighlight`, `DeleteHighlight`, `ListBooks`, `ListHighlights`, `ListHighlightsByBook`. The `sqlite` package provides the persistent storage implementation.
- Request/response shapes: see `internal/web/types.go` for JSON shapes used by Moon+ Reader (`CreateHighlightRequest`) and Obsidian (`ListHighlightsResponse`, `ListBooksResponse`).
- Web UI: A browser-based interface at `/ui/` allows viewing, editing, and deleting highlights. UI uses vanilla JavaScript with Bootstrap 5 for modals and styling. See `static/index.html` and `static/js/app.js`.

Developer workflows
- Build: run `./build.sh -b` — this sets `cfg/VERSION`, `cfg/HASH`, builds `bin/unwise` and uses `go build` flags from the script.
- Test: run `./build.sh -t` (build then `go test -race ./...` with coverage). Use `./build.sh -g` to generate mocks (requires `mockery`).
- Docker: `./build.sh -d` or `docker-compose -f docker/docker-compose.yml up` or `docker run -p 3123:3123 -e TOKEN=my-token ghcr.io/corani/unwise:latest`.
- Environment: configuration via environment variables or `.env` (loaded by `internal/config.Load()`). Key env vars: `LOGLEVEL`, `REST_ADDR`, `REST_PATH`, `DATA_PATH`, `TOKEN`, `DROP_TABLE`.

Project-specific conventions & patterns
- Use `internal/config`'s `Config` and `conf.Logger` for logging; prefer `conf.Logger.Info/Errorf` over fmt.Println.
- Time format: timestamps are RFC3339 strings (look for time.Parse(time.RFC3339) usage in `internal/storage/sqlite`).
- Defaults and normalization: `storage.AddBook` sets a default title `storage.DefaultTitle` when the incoming title is empty.
- Error handling: handlers return Go errors which are converted to JSON by `web.Server.HandleError`; prefer wrapping errors with `fmt.Errorf("%w: ...", fiber.ErrBadRequest)` where a Fiber `*fiber.Error` is intended.
- Concurrency: The `sqlite` package handles concurrency at the database level.

Integration points & external dependencies
- HTTP framework: Fiber v2 (`github.com/gofiber/fiber/v2`). Handlers attach middleware for `logger`, `helmet`, and `keyauth` (Token auth). Review `internal/web/app.go` if changing routing or middleware ordering.
- SQLite: uses `modernc.org/sqlite`; connections are opened in `internal/storage/sqlite` and DDL is applied during `Init`.
- Env parsing: `github.com/caarlos0/env/v11` + `joho/godotenv` for `.env` files in `internal/config`.

Tests & mocking
- Tests live next to packages (e.g. `internal/storage/*_test.go`, `cmd/unwise/*_test.go`). Tests use `sqlmock` for `sqlite` behavior.
- Test structure: Use table-driven tests with a `setup func(*testing.T, sqlmock.Sqlmock)` parameter for mock configuration. This pattern keeps tests maintainable and consistent across the codebase.
- When adding tests that depend on time ranges, pay attention to `parseISO8601Datetime` and functions that treat zero `time.Time` as special (see `ListBooks` / `ListHighlights`).

Examples & quick references
- Add highlight flow (handler): `internal/web/handlers.go:HandleCreateHighlights` — parse `CreateHighlightRequest`, call `stor.AddBook` then `stor.AddHighlight`, and build `CreateHighlightResponse`.
- Storage contract: `internal/storage/storage.go` defines the exact behavior to satisfy across implementations.
- SQL schema: `internal/storage/sqlite/sqlite.go` (CREATE TABLE for `books` and `highlights`) — keep UNIQUE constraints in mind when changing insert/update logic.

When editing files
- Run `./build.sh -t` locally after code changes to run tests and catch race detector issues.
- Follow existing style: small, focused functions; avoid adding global state. Keep `config` and `storage` as injectable dependencies for handlers.
- If touching SQL, mirror the `ON CONFLICT` expressions used in sqlite implementation; tests expect `RETURNING` rows to be available.

Working with the user
- If unsure about intended behavior for edge-cases (e.g., duplicate highlights, retention strategy), propose a change in a single small PR and run tests.
- Ask the maintainer whether to persist data by default or keep in-memory for quick iteration; `DROP_TABLE` env var controls table reset on startup.

Files to review first
- cmd/unwise/main.go
- internal/web/web.go
- internal/web/app.go
- internal/web/handlers.go
- internal/web/types.go
- internal/config/config.go
- internal/storage/storage.go
- internal/storage/sqlite/sqlite.go
- build.sh
