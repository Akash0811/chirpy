# Chirpy

Chirpy is a Go HTTP API backed by PostgreSQL. It supports user registration, login, chirp creation and retrieval, token refresh and revocation, a webhook for Chirpy Red upgrades, and a small static app served from `/app`.

Detailed endpoint documentation lives in `docs/api.md`.

## Tech Stack

- Go
- PostgreSQL
- `psql`
- `goose` for schema migrations
- `sqlc` for type-safe database access generation

## Prerequisites

Install the following tools before running the project:

### 1. Install Go

The project declares Go `1.26.2` in `go.mod`. Install a compatible Go `1.26.x` release.

- Verify installation:

```bash
go version
```

### 2. Install goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- Verify installation:

```bash
goose -version
```

### 3. Install sqlc

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

- Verify installation:

```bash
sqlc version
```

### 4. Install PostgreSQL and psql

Install PostgreSQL for your operating system. Make sure the `psql` CLI is also available.

- Verify installation:

```bash
psql --version
```

## Local Setup

### 1. Clone the project and install Go dependencies

```bash
go mod download
```

### 2. Create the database

Example using `psql`:

```bash
createdb chirpy
```

If `createdb` is unavailable, use:

```bash
psql -U postgres -c "CREATE DATABASE chirpy;"
```

If your PostgreSQL installation requires it, enable `pgcrypto` so `gen_random_uuid()` is available:

```bash
psql -U postgres -d chirpy -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
```

### 3. Create an environment file

Create `.env` in the project root with placeholder values only:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=replace-with-a-long-random-secret
POLKA_KEY=replace-with-your-polka-webhook-key
```

Environment variables used by the app:

- `DB_URL`: PostgreSQL connection string
- `PLATFORM`: set to `dev` if you want `/admin/reset` to work locally
- `JWT_SECRET`: secret used to sign access tokens
- `POLKA_KEY`: shared API key for the Polka webhook endpoint

Do not commit real secrets.

### 4. Run database migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 5. Generate database code with sqlc

```bash
sqlc generate
```

### 6. Run the server

```bash
go run .
```

The server listens on `http://localhost:8080`.

## Useful Commands

Run tests:

```bash
go test ./...
```

Regenerate SQL code after changing queries or schema:

```bash
sqlc generate
```

Re-run migrations against the configured database:

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

## Project Structure

- `main.go`: application entrypoint and route registration
- `internal/backend`: HTTP handlers and API behavior
- `internal/auth`: JWT, password, API key, and refresh token helpers
- `internal/database`: generated `sqlc` database access layer
- `sql/schema`: `goose` migrations
- `sql/queries`: SQL queries used by `sqlc`
- `docs/api.md`: endpoint-by-endpoint API reference

## API Documentation

See `docs/api.md` for:

- available resources
- endpoint paths
- supported HTTP methods
- request and response shapes
- authentication requirements
- example `curl` requests
