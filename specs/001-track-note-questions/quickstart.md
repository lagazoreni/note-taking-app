# Quickstart: Interactive Note Questions

This guide describes the planned development and validation workflow. Source directories are created during implementation.

## Prerequisites

- Go 1.26.x
- Node.js 24.x and npm
- A container engine compatible with OCI images and Compose
- `make`
- A supported desktop browser: one of the latest two stable Chrome, Edge, or Firefox major versions, or Safari 18 or newer

Internet access is needed only to install dependencies and build images. The running application’s core workflows require no internet connection.

## Configuration

Copy the committed example:

```bash
cp .env.example .env
```

Planned keys:

```dotenv
NOTED_ADDR=127.0.0.1:8080
NOTED_DATA_DIR=./.local/data
NOTED_DB_FILE=notes.db
NOTED_LOG_LEVEL=info
NOTED_IMPORT_MAX_BYTES=268435456
NOTED_IMPORT_TTL=24h
```

Production/container values use the same names. `NOTED_DATA_DIR` is `/data` in the container. No secrets or authentication values are required for the MVP.

## Start development

Terminal 1 — backend:

```bash
make dev-server
```

Terminal 2 — frontend with API proxy to the local Go service:

```bash
make dev-web
```

Open the URL printed by SvelteKit. The development proxy is development-only; production uses one origin.

## Build independently

```bash
make build-web
make build-server
```

The web build produces static assets. The server build remains testable without those assets by using an empty/test asset provider. The production build copies the completed web assets into the server artifact.

## Run checks

```bash
make fmt
make lint
make test-server
make test-web
make test-contract
make test-e2e
```

Run all required checks:

```bash
make check
```

The contract check validates `contracts/openapi.yaml`, exercises backend request/response examples, and runs frontend API-client tests against shared fixtures.

## Build and run the container

```bash
make image
mkdir -p .local/container-data
NOTED_DATA_DIR="$(pwd)/.local/container-data" docker compose up --build
```

Open `http://127.0.0.1:8080`.

Health checks:

```bash
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

The first endpoint verifies the process. The second verifies that migrations completed and SQLite is usable.

## Verify offline behavior

1. Start the local Go service and open the application once.
2. Keep the same-origin localhost application/API reachable and block only external network requests. Do not use a browser-wide Offline mode that also blocks localhost. The automated Playwright setup enforces this condition through request interception.
3. Refresh the application and verify the cached shell loads without any external request succeeding.
4. Create a workspace, note, and two questions.
5. Edit one shared question from another linked note and verify all views update.
6. Search, sort, and filter while internet access remains disabled.
7. Stop the Go process and attempt an edit. Verify the editor preserves unsaved input and provides Retry rather than reporting a false save.
8. Restart Go and retry successfully.

Automated equivalent, executed in Chromium, Firefox, and WebKit projects:

```bash
make test-offline
```

## Verify persistence

```bash
docker compose down
docker compose up -d
```

Confirm previously created notes remain. Then run the export/import browser scenario:

```bash
make test-import-export
```

It verifies a full round trip, invalid archive rejection, conflict review, cancellation, and rollback.

## Database diagnostics

Development-only commands:

```bash
make db-status
make db-check
make db-reset       # destructive; refuses outside development data paths
```

`db-check` runs SQLite integrity checks, foreign-key checks, and verifies FTS5 availability. Migration files are append-only after release.

## Contract locations

- HTTP API: `specs/001-track-note-questions/contracts/openapi.yaml`
- Export archive: `specs/001-track-note-questions/contracts/export.schema.json`
- Markdown question directive: `specs/001-track-note-questions/contracts/note-content.md`

## Release smoke test

```bash
make smoke-image
```

The smoke test must verify:

- clean image startup as non-root;
- health and readiness;
- automatic migrations;
- static frontend and `/api/v1` on one origin;
- workspace → note → question → answer workflow;
- restart persistence through `/data`;
- structured logs without note, question, or answer content.
