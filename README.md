# Noted

Noted is a SvelteKit frontend backed by a Go HTTP API and a local SQLite database.
This README explains how to run both services locally and how to run the test suites.

## Prerequisites

Install the following tools:

- Node.js 24.x and npm
- Go 1.26 or newer
- Git
- Docker Desktop (optional, only needed for the Docker workflow)

The repository contains two npm projects:

- The repository root contains the Playwright end-to-end test runner.
- [`web/`](web/) contains the Svelte frontend.

## Initial setup

From the repository root, install the dependencies once:

### PowerShell, macOS, or Linux

```sh
npm ci
npm --prefix web ci
go -C server mod download
```

Install the Playwright browsers if you intend to run end-to-end tests:

```sh
npx playwright install
```

No `.env` file is required for local development. The Go service reads its configuration from environment variables and has usable defaults.

## Run the application locally

The local development setup uses two processes:

- Go backend: `http://127.0.0.1:8080`
- Vite frontend: `http://127.0.0.1:5555`

The Vite development server proxies `/api`, `/healthz`, and `/readyz` to the backend. This proxy is configured in [`web/vite.config.ts`](web/vite.config.ts).

### 1. Start the backend

Open a terminal at the repository root.

#### PowerShell

```powershell
New-Item -ItemType Directory -Force .local\data | Out-Null

$env:NOTED_ADDR = "127.0.0.1:8080"
$env:NOTED_DATA_DIR = (Resolve-Path .local\data).Path
$env:NOTED_DB_FILE = "notes.db"

go -C server run ./cmd/noted

#### macOS/Linux or Git Bash

```sh
mkdir -p .local/data
NOTED_ADDR=127.0.0.1:8080 \
NOTED_DATA_DIR="$PWD/.local/data" \
NOTED_DB_FILE=notes.db \
go -C server run ./cmd/noted
```

The database is created at `.local/data/notes.db`. Keep this terminal running.

Verify the backend from another terminal:

```sh
curl http://127.0.0.1:8080/healthz
curl http://127.0.0.1:8080/readyz
```

On PowerShell, the equivalent commands are:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/healthz
Invoke-RestMethod http://127.0.0.1:8080/readyz
```

### 2. Start the frontend

Open a second terminal at the repository root:

```sh
npm --prefix web run dev -- --host 127.0.0.1
```

Open the application in a browser:

```text
http://127.0.0.1:5555
```

Stop either process with `Ctrl+C`.

### Using Make instead

If `make` is installed, the Makefile provides equivalent commands. Run each command in a separate terminal from the repository root:

```sh
make dev-server
make dev-web
```

`make` is generally easiest to use from macOS/Linux, Git Bash, or WSL. The explicit commands above are recommended for PowerShell.

## Run with Docker

Docker builds the frontend and embeds it into the Go server, so both parts run as one service:

```sh
docker compose up --build
```

Open:

```text
http://localhost:8080
```

The database is stored in the Docker volume `noted-data`. Stop the service with:

```sh
docker compose down
```

To also delete the Docker database and all stored notes:

```sh
docker compose down -v
```

## Run tests

Run frontend unit tests:

```sh
npm --prefix web test
```

Run backend tests:

```sh
go -C server test ./...
```

Run frontend type checking and linting:

```sh
npm --prefix web run check
npm --prefix web run lint
```

Run the full Playwright end-to-end suite:

```sh
npm run test:e2e
```

The Playwright configuration in [`playwright.config.ts`](playwright.config.ts) starts its own Go backend and Vite frontend. Stop manually started development servers before running this command. The E2E frontend uses port `5174` so it does not conflict with the normal development port `5555`.

The E2E HTML report is written to `playwright-report/`.

## Useful database commands

The Go service includes database diagnostics. If using PowerShell, set the same data directory environment variable before running a command:

```powershell
$env:NOTED_DATA_DIR = (Resolve-Path .local\data).Path
go -C server run ./cmd/noted db-status
go -C server run ./cmd/noted db-check
go -C server run ./cmd/noted db-reset
```

`db-reset` deletes the local SQLite database. The next backend start recreates it through the migrations.

## Troubleshooting

### Port 8080 is already in use

Stop the process using port 8080, or start the backend on another port and update the Vite proxy in [`web/vite.config.ts`](web/vite.config.ts). For example:

```powershell
$env:NOTED_ADDR = "127.0.0.1:8081"
go -C server run ./cmd/noted
```

Then change the proxy target from `http://127.0.0.1:8080` to `http://127.0.0.1:8081`.

### The frontend loads but API calls fail

Make sure the Go backend is still running and that these endpoints return successfully:

```text
http://127.0.0.1:8080/healthz
http://127.0.0.1:8080/readyz
```

Also confirm that the frontend is running on port `5555` and that its requests use paths beginning with `/api`.
