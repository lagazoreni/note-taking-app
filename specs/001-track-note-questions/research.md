# Phase 0 Research: Interactive Note Questions

**Date**: 2026-08-09  
**Feature**: `001-track-note-questions`

All technical unknowns from the plan are resolved below.

## 1. SvelteKit delivery model

**Decision**: Build the SvelteKit application as static client-rendered assets with `@sveltejs/adapter-static`. Cache versioned shell/assets using SvelteKit’s service-worker support. In production, copy the built assets into the Go service artifact and serve them from the same origin as `/api/v1`.

**Rationale**: The application has no public SEO or server-rendering requirement, and Go must own backend behavior. Static output keeps the runtime to one local process, removes CORS and a frontend server, supports offline shell loading, and still leaves frontend and backend independently buildable. SvelteKit documents both static adaptation and integrated service workers.

**Alternatives considered**:

- SvelteKit Node server plus Go API: rejected because it adds a second runtime process without current user value.
- SvelteKit server endpoints: rejected because they would duplicate or blur the constitutionally required Go business boundary.
- A pure SPA outside SvelteKit: rejected because the user selected SvelteKit and its routing/service-worker conventions are useful.

**References**: https://svelte.dev/docs/kit/adapter-static, https://svelte.dev/docs/kit/service-workers

## 2. Meaning of offline for a Go-backed web MVP

**Decision**: Offline means independent of internet access, with the local Go application running. The service worker caches the UI shell, but API writes always go directly to the authoritative local backend. Failed writes remain visibly unsaved in the editor and can be retried; they are not silently queued by the service worker.

**Rationale**: A browser cannot persist through the selected Go/SQLite authority when the local service is absent. Adding a second browser database and replication protocol would contradict the explicit SvelteKit-Go boundary, expand conflict handling beyond MVP scope, and duplicate business rules. Localhost remains available without internet.

**Alternatives considered**:

- IndexedDB write queue synchronized to Go: rejected for MVP because it creates dual persistence, replay semantics, and conflicts.
- SQLite in WebAssembly: rejected because the user selected Go as backend and this would move persistence/business behavior into the frontend.
- Network-only frontend with no service worker: rejected because shell reloads would fail without internet if assets were not already available locally.

## 3. Go baseline and HTTP stack

**Decision**: Use Go 1.26, standard `net/http` routing/middleware, `encoding/json`, and `log/slog`. Use explicit domain, service, transport, and SQLite store packages.

**Rationale**: Go 1.26 is the current stable line at planning time. Modern `net/http` supports method-aware routing, so an external router is not justified. Standard structured logging satisfies observability while minimizing dependencies.

**Alternatives considered**:

- Gin, Echo, or Fiber: rejected because MVP routing needs do not outweigh an extra framework and conventions.
- A monolithic handler/store package: rejected because business rules and persistence must remain independently testable under the constitution.
- Microservices: rejected because one local user and one database require no distributed architecture.

**Reference**: https://go.dev/doc/devel/release

## 4. SQLite driver and operating mode

**Decision**: Use `database/sql` with `modernc.org/sqlite`. Enable foreign keys for every connection, WAL mode, a bounded busy timeout, and explicit transactions. Limit connection concurrency appropriately for one local database and run embedded, numbered SQL migrations at startup before readiness.

**Rationale**: The pure-Go driver avoids CGO toolchain and runtime-library complexity in reproducible containers while exposing SQLite through the standard library interface. SQLite fits hierarchical notes, many-to-many question links/tags, atomic destructive operations, and local operation. FTS5 is available in the selected SQLite engine and must be verified by startup/integration tests.

**Alternatives considered**:

- `mattn/go-sqlite3`: mature and capable, but rejected to avoid CGO and a larger production build toolchain.
- PostgreSQL: rejected because it requires a separate service and provides no benefit for one local user.
- IndexedDB or a document database: rejected because the backend is Go and the domain relies heavily on relational constraints and joins.
- An ORM: rejected because explicit SQL keeps queries, constraints, and performance visible and avoids an unnecessary abstraction.

**References**: https://pkg.go.dev/modernc.org/sqlite, https://sqlite.org/foreignkeys.html, https://sqlite.org/fts5.html, https://sqlite.org/pragma.html

## 5. SQLite durability and concurrency

**Decision**: Use WAL for responsive concurrent reads, a busy timeout for short writer contention, foreign keys, and transactional aggregate updates. Run `PRAGMA optimize` periodically after meaningful write volume and on orderly shutdown when practical. Keep the database and WAL sidecars in one declared persistent directory.

**Rationale**: The workload is read-heavy with one local user, but autosave/search can overlap. WAL allows readers during writes. Transactions are required for links, deletion review execution, migrations, and imports. SQLite recommends foreign keys be explicitly enabled and recommends `PRAGMA optimize` as the modern statistics maintenance path.

**Alternatives considered**:

- Rollback journal: simpler, but gives poorer read/write overlap.
- Multiple database files per workspace: rejected because cross-workspace search, shared tag visibility, export, and migrations become harder.
- In-memory caching: rejected until measured need exists.

**References**: https://sqlite.org/wal.html, https://sqlite.org/transactional.html, https://sqlite.org/pragma.html

## 6. IDs and optimistic concurrency

**Decision**: Use globally unique UUID-form IDs generated from cryptographically secure randomness and an integer `version` on mutable records. Mutations carry the version last read; stale changes receive HTTP 409 with current record metadata.

**Rationale**: Globally unique IDs make portable export/import and future synchronization safer. Optimistic versions prevent one tab from silently overwriting another and are more reliable than timestamp comparison. UUID generation can be implemented with Go’s cryptographic standard library, so no runtime dependency is required.

**Alternatives considered**:

- Auto-increment IDs: rejected because imports and future multi-device data can collide.
- Last-write-wins timestamps: rejected because clock differences and silent loss violate requirements.
- Pessimistic edit locks: rejected because they complicate tab crashes and local workflows.

## 7. Note format and inline question references

**Decision**: Store canonical note text as Markdown. Support paragraphs, ordered lists, and unordered lists in MVP. Represent inline questions with the documented directive `{{question:<uuid>}}`; store display mode and referential integrity in `note_questions`. Render preview using `marked`, disable raw HTML, and sanitize output with DOMPurify.

**Rationale**: Markdown is fast to type, human-readable, portable, and naturally supports the MVP’s structured text. A stable directive places a shared question in context without duplicating it. The relation table preserves database integrity and display metadata. A mature parser plus sanitization is safer than custom HTML parsing.

**Alternatives considered**:

- Plain text: rejected because lists and future formatting are weaker.
- Rich-text JSON/ProseMirror: rejected for MVP due to editor complexity, format lock-in, and larger dependency surface.
- Copying question text into Markdown: rejected because edits would diverge.
- Raw embedded HTML: rejected due to portability and injection risk.

## 8. Search and filtering

**Decision**: Use SQLite FTS5 external-content indexes for note title/body and question text/answer, maintained transactionally with triggers. Use ordinary indexed columns and joins for status, workspace, topic, tag, dates, answer/link presence, priority, and sorting. Search snippets are escaped before display.

**Rationale**: FTS5 provides ranking, tokenization, and scalable text search without another service. Relational predicates are clearer and faster for structured filters. Stable `(sort_column, id)` ordering supports deterministic cursor pagination.

**Alternatives considered**:

- SQL `LIKE` scans: rejected against the 10,000-note/50,000-question target.
- Browser-side search: rejected because it would require transferring all data and duplicate backend logic.
- Elasticsearch/Meilisearch: rejected because an extra service violates lightweight local delivery.

**Reference**: https://sqlite.org/fts5.html

## 9. Reminder behavior in browsers

**Decision**: Evaluate pending reminders while the application is open and on startup/resume. Show in-app alerts in all cases; use the Notifications API only after explicit permission and when supported. Mark past undelivered reminders as missed/overdue on return. Do not promise reliable scheduled notifications while every browser context is closed.

**Rationale**: Web Notifications require permission, secure contexts, and have uneven support. Standard web APIs do not provide universally reliable offline scheduled local notifications after the browser closes. The decision meets the specification’s platform-limit language without adding a remote push service.

**Alternatives considered**:

- Push notifications: rejected because they require online infrastructure and permissions beyond MVP.
- Polling from the Go backend to a desktop notification service: rejected because the MVP is web-only and this creates platform coupling.
- Claiming closed-browser scheduling: rejected because it is not portable or reliably testable across target browsers.

**Reference**: https://developer.mozilla.org/en-US/docs/Web/API/Notifications_API

## 10. Import/export format and safety

**Decision**: Export one ZIP archive containing `manifest.json`, normalized JSON data files, and a checksum file. Version the format and publish `export.schema.json`. Import follows upload → validate/preview → conflict decisions → atomic apply. Staged files live under the configured data directory and expire after 24 hours; cancellation removes them.

**Rationale**: ZIP allows future attachments without changing the outer format, while JSON is inspectable and portable. A manifest supports schema evolution. Preview plus one SQLite transaction prevents partial restoration and silent replacement.

**Alternatives considered**:

- Raw SQLite file export: rejected because it is not a stable interchange contract across schema versions and is unsafe while active.
- One large JSON response: rejected because ZIP has cleaner integrity/versioning and future attachment support.
- Applying records during validation: rejected because failures could partially mutate data.

**Reference**: https://sqlite.org/backup.html

## 11. API contract and testing

**Decision**: Define a same-origin `/api/v1` JSON API using OpenAPI 3.1. Use `kin-openapi` only in Go contract tests to validate the document and handler traffic. Keep shared request/response examples under contract tests and consume them in frontend API client tests.

**Rationale**: OpenAPI documents the constitution-required stable boundary. Backend schema validation and frontend fixture consumption catch provider/consumer drift without generating runtime coupling. Same-origin deployment removes CORS configuration from MVP.

**Alternatives considered**:

- GraphQL: rejected because the domain does not need client-defined graph queries and it adds schema/runtime machinery.
- SvelteKit server actions: rejected because Go must own APIs and business rules.
- Code generation as a required build step: deferred because it adds tooling; handwritten small clients with contract tests are sufficient initially.

## 12. Container topology

**Decision**: Use a multi-stage build with pinned Node and Go builders and a minimal pinned non-root runtime image. The final Go process serves static web assets and API traffic on one port, exposes `/healthz` and `/readyz`, logs to stdout/stderr, and stores mutable files only under `/data`.

**Rationale**: One process is the simplest operable local deployment and supports offline use. Multi-stage builds retain independent frontend/backend build stages while excluding build toolchains from production.

**Alternatives considered**:

- Separate frontend reverse-proxy container: rejected because it adds a service and same-origin routing configuration.
- Host-installed Go/Node runtime: supported for development but rejected as the reproducible delivery artifact.
- Embedding the SQLite database in the image: rejected because updates would destroy user data and violate the constitution.
