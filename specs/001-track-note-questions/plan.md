# Implementation Plan: Interactive Note Questions

**Branch**: `001-track-note-questions` | **Date**: 2026-08-09 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-track-note-questions/spec.md`

## Summary

Build an offline-capable, local-first web application with a statically built SvelteKit frontend and a Go HTTP backend backed by one SQLite database. The Go service is authoritative for question lifecycle rules, note/question relationships, deletion review, filtering, search, and atomic import/export. The SvelteKit application provides responsive note editing, question views, workspace navigation, and a cached application shell. Production uses one container and one origin: the Go process serves the versioned API and the built frontend assets, while `/data` is the declared persistent volume.

Notes use portable Markdown for paragraphs and lists, plus a documented question-reference directive for inline shared questions. SQLite relational tables model shared questions and tags; FTS5 indexes notes, questions, and answers. HTTP interfaces are documented with OpenAPI 3.1, use stable error envelopes and optimistic versions, and are tested on both frontend and backend.

## Technical Context

**Language/Version**: Go 1.26; TypeScript with Svelte 5 and the current SvelteKit release compatible with the committed lockfile

**Primary Dependencies**: SvelteKit, `@sveltejs/adapter-static`, Vite, a service worker, `marked` plus DOMPurify for safe Markdown preview, Go standard `net/http` and `log/slog`, `modernc.org/sqlite`, and `kin-openapi` for contract tests only

**Storage**: SQLite with foreign keys, WAL mode, busy timeout, transactional SQL migrations, and FTS5; database and temporary import files live under the configured data directory

**Testing**: Go `testing`/`httptest` with temporary SQLite databases; Vitest and Svelte Testing Library for frontend units/components; Playwright Chromium, Firefox, and WebKit projects for browser workflows and no-internet tests; shared API examples plus OpenAPI validation for contract tests; container smoke test

**Target Platform**: The latest two stable major versions available at release time of desktop Chrome, Edge, and Firefox, plus Safari 18 or newer; local Linux container runtime for delivery; desktop and Android clients are future targets

**Project Type**: Local-first web application with independently buildable frontend and backend and a combined production image

**Performance Goals**: Show local question edits in all active views within 1 second; return 95% of searches over 10,000 notes and 50,000 questions within 2 seconds; filter/sort 50,000 questions within 2 seconds; create or update common records within 500 ms under normal local use

**Constraints**: Core use must require no internet connection while the local Go service remains running and reachable; no authentication; one local user; backend is authoritative; no silent data loss; imports and destructive note operations are atomic; cross-workspace note/question moves are outside MVP scope; frontend and API remain usable on narrow and wide viewports; production runs non-root with one declared persistent volume

**Scale/Scope**: One local user, multiple workspaces, up to 10,000 notes and 50,000 questions, seven MVP user journeys, approximately 25 application routes/components and a versioned HTTP API

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle / constraint | Pre-design gate | Design evidence | Post-design gate |
|---|---|---|---|
| I. Lightweight by Default | PASS | One Go process, one SQLite file, standard Go routing/logging, static SvelteKit output, no cache/queue/auth service; each non-standard dependency has a specific persistence, rendering, or test role. | PASS |
| II. Notes Stay Fast and Focused | PASS | Note editing is a primary route; optimistic UI is reconciled with authoritative API responses; unsaved-change and recoverable-error states are required; FTS5 and indexed filters address scale goals. | PASS |
| III. Explicit SvelteKit-Go Boundary | PASS | SvelteKit owns browser UI/routing; Go owns all validation, business rules, and persistence; `contracts/openapi.yaml` defines `/api/v1`; frontend validation is advisory only. | PASS |
| IV. Container-Ready Delivery | PASS | Multi-stage reproducible build, pinned base images, non-root runtime, `/healthz`, structured standard-output logs, environment configuration, and `/data` volume are specified. | PASS |
| V. Tested, Observable, and Secure | PASS | Unit, integration, contract, browser, offline, and container smoke tests are planned; request IDs and redacted structured logs are required; Markdown is sanitized; all input is validated. | PASS |
| Frontend and backend independently buildable | PASS | `web/` and `server/` have independent build and test commands; production combines only their build artifacts. | PASS |
| Same configuration and entrypoint model across environments | PASS | documented environment keys are used locally and in containers; no environment values enter source control. | PASS |

No constitution violations or exceptions require Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/001-track-note-questions/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── openapi.yaml
│   ├── export.schema.json
│   └── note-content.md
└── tasks.md                 # created by /speckit-tasks, not this plan
```

### Source Code (repository root)

```text
server/
├── cmd/noted/
│   └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── contracttest/
│   ├── domain/
│   ├── service/
│   ├── store/
│   │   ├── sqlite/
│   │   └── migrations/
│   ├── importexport/
│   └── search/
├── webdist/                 # generated frontend artifact; not hand-edited
├── go.mod
└── go.sum

web/
├── src/
│   ├── lib/
│   │   ├── api/
│   │   ├── components/
│   │   ├── editor/
│   │   ├── stores/
│   │   └── types/
│   ├── routes/
│   │   ├── +layout.svelte
│   │   ├── notes/
│   │   ├── questions/
│   │   ├── answered/
│   │   ├── search/
│   │   ├── workspaces/
│   │   └── settings/data/
│   ├── service-worker.ts
│   └── app.html
├── static/
├── tests/
│   ├── unit/
│   ├── component/
│   └── contract/
├── package.json
└── svelte.config.js

tests/
├── e2e/
├── fixtures/
└── smoke/

containers/
└── Containerfile

compose.yaml
Makefile
.env.example
```

**Structure Decision**: Use a two-project web layout because the constitution requires an explicit SvelteKit-Go boundary. `web/` produces static browser assets and remains independently testable. `server/` owns domain logic, SQLite, migrations, and `/api/v1`. A multi-stage production build copies the web artifact into the Go runtime artifact so one local process and one origin are sufficient; this avoids an extra reverse proxy while preserving build boundaries.

## Design Decisions

### Runtime and offline boundary

- “Offline” means no internet connection is needed after installation and first load, while the local Go service remains running and reachable because it owns persistence and business rules. Automated tests block external network requests but keep the same-origin localhost service available; they do not use a browser-wide offline mode that also disables localhost.
- The service worker caches only versioned frontend shell/assets and a fallback page. It does not cache or replay mutation API calls, avoiding two competing sources of truth.
- The UI retains unsaved editor text during transient API failures, displays a recoverable error, and lets the user retry. Successful mutations replace local view state with the API response and invalidate affected query stores.
- Reminder checks occur while the application is open. With permission, browser notifications are used; otherwise reminders remain visible in-app. Missed reminders are surfaced on the next open. Reliable closed-browser scheduled alerts are explicitly outside MVP browser capabilities.

### Persistence and consistency

- SQLite runs with `PRAGMA foreign_keys=ON`, WAL journal mode, a busy timeout, and explicit transactions for multi-record operations.
- Every mutable aggregate has an integer `version`. Update and delete requests include the last observed version; stale writes return `409 VERSION_CONFLICT` with current metadata.
- Shared questions are stored once. `note_questions` records links and presentation. The backend validates note/question workspace ownership and Markdown question directives in one transaction. A topic filter matches a question when at least one linked note is assigned to that topic. Cross-workspace note and question moves are rejected in the MVP.
- Note deletion is two-step: preview impact, then execute with the preview token, note version, and explicit decision for each singly linked question. Any mismatch invalidates the preview and requires review again.
- Imports are validate/preview/apply. Apply uses one transaction and never mutates current data on validation failure, cancellation, or interruption before commit.

### Editing and content

- `bodyMarkdown` is the canonical note body. MVP syntax includes paragraphs and ordered/unordered lists.
- Inline questions use `{{question:<id>}}`; presentation is stored on the corresponding note-question link, not encoded in visible prose.
- The backend rejects malformed, duplicated, foreign-workspace, or unlinked directives. The frontend inserts/removes directives through editor controls so users do not need to type IDs.
- Preview HTML is rendered from Markdown and sanitized before insertion into the document. Raw HTML in Markdown is disabled.

### API and errors

- Same-origin JSON API under `/api/v1`; export download and import upload are the only non-JSON payloads.
- IDs are lowercase UUID strings generated from cryptographically secure randomness. Dates are RFC 3339 UTC strings; calendar-only due dates are represented separately where applicable.
- Errors use `{ "error": { "code", "message", "fieldErrors", "requestId", "details" } }`. Messages are user-safe and logs exclude note/question/answer bodies.
- List endpoints use opaque cursor pagination and deterministic tie-breaking by ID. Filter and sort options are allow-listed.

### Testing strategy

- Domain tests cover status transitions, workspace boundaries, hierarchy cycles, tag access, and deletion decisions.
- SQLite integration tests cover migrations, constraints, FTS synchronization, transactions, import rollback, and query performance fixtures.
- Handler contract tests validate requests/responses against OpenAPI and shared examples. Frontend API tests consume the same examples.
- Playwright Chromium, Firefox, and WebKit projects cover every spec user journey, responsive layouts, external-network blocking with localhost available, refresh recovery, and unsaved-edit errors. Release checks cover the latest two stable Chrome, Edge, and Firefox major versions and Safari 18 or newer.
- Container smoke tests verify non-root startup, health, migration, persistent-volume restart, static shell, API, and log redaction.

## Implementation Phases

### Phase 0 - Research

Research decisions are captured in [research.md](./research.md), including offline boundaries, static SvelteKit delivery, SQLite driver/configuration, Markdown editing, reminders, search, contracts, and import/export safety. All technical unknowns are resolved.

### Phase 1 - Foundation and contracts

1. Establish independent `web/` and `server/` builds plus shared developer commands.
2. Commit OpenAPI, export schema, and note-content grammar before handlers and clients.
3. Implement SQLite migration runner and schema from [data-model.md](./data-model.md).
4. Establish API envelope, request IDs, structured redacted logging, health/readiness, and contract-test harness.
5. Establish SvelteKit shell, service worker, API client boundary, error handling, and responsive navigation.
6. Add reproducible combined container and persistent-volume smoke test.

### Phase 2 - Feature delivery sequence

1. Workspaces, topics, tags, and note hierarchy.
2. Markdown note editor and shared inline question links.
3. Question lifecycle, answer rules, active/answered views, due dates, reminders, filters, and sorting.
4. Cross-scope FTS search and navigation.
5. Two-step safe note deletion.
6. Atomic export, import validation, conflict review, and restore.
7. Performance, offline, accessibility, security, and release verification against success criteria.

## Complexity Tracking

No constitution violations require justification.
