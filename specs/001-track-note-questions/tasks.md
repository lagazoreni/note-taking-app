# Tasks: Interactive Note Questions

**Input**: Design documents from `/specs/001-track-note-questions/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**New-session start**: Read [`PROJECT_MAP.md`](../../PROJECT_MAP.md) at the repository root before opening other files. For highlight/annotation work, use Phase 11. For answer-in-context, Next queue, deferred dates, resolved highlights, sturdier wrapping, reading-view link/attach, and keyboard capture, use **Phase 13–20 in this file** — read only that phase, then only the files it names. Do not scan the whole repo unless the map is stale.

**Tests**: Automated tests are included because the project constitution requires unit, integration, API contract, browser workflow, and container smoke coverage. Within each story, create the listed tests first and verify they fail for the expected missing behavior before implementation.

**Organization**: Tasks are grouped by user story. Every story ends with an independently executable browser scenario and checkpoint.

**Workspace prerequisite**: A fresh database can contain zero workspaces. The first P1 capture flow must create and select a workspace before exposing note creation; workspace bootstrap, empty-state guidance, persisted-selection recovery, and the workspace-required note guard are explicit deliverables rather than assumptions.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can be implemented in parallel with adjacent tasks because it changes different files and has no dependency on an unfinished adjacent task.
- **[Story]**: Maps the task to a user story in `spec.md`.
- Paths are relative to the repository root.
- `a`-suffixed IDs are gap-remediation tasks inserted without renumbering later tasks.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize independently buildable Go and SvelteKit projects, shared commands, and delivery scaffolding.

- [X] T001 Create the Go module, `cmd/noted` entrypoint stub, and planned internal package directories in `server/go.mod` and `server/cmd/noted/main.go`
- [X] T002 [P] Scaffold the SvelteKit TypeScript application with static adapter and locked dependencies in `web/package.json`, `web/package-lock.json`, and `web/svelte.config.js`
- [X] T003 Configure Svelte, TypeScript, Vitest, and formatting checks in `web/tsconfig.json`, `web/vite.config.ts`, `web/eslint.config.js`, and `web/.prettierrc`
- [X] T004 Add repository build, development, formatting, linting, test, database, image, and smoke targets in `Makefile`
- [X] T005 [P] Add local artifact exclusions and documented runtime defaults in `.gitignore` and `.env.example`
- [X] T006 Add the pinned multi-stage non-root build and one-origin runtime skeleton in `containers/Containerfile` and `compose.yaml`
- [X] T007 [P] Configure Chromium, Firefox, and WebKit Playwright projects, local web server startup, external-network blocking with localhost allowed, and browser support documentation in `playwright.config.ts` and `tests/e2e/README.md`
- [X] T008 [P] Add shared TypeScript test setup and DOM matchers in `web/tests/setup.ts` and `web/vitest.config.ts`
- [X] T009 Copy contract examples into version-controlled provider/consumer fixtures and document regeneration rules in `tests/fixtures/contracts/README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish runtime configuration, persistence, migrations, API conventions, frontend shell, offline shell caching, and test harnesses required by every story.

**⚠️ CRITICAL**: No user story implementation begins until this phase passes its checkpoint.

- [X] T010 [P] Implement environment parsing and validation for address, data directory, database, log level, import limit, and import TTL in `server/internal/config/config.go` and `server/internal/config/config_test.go`
- [X] T011 [P] Implement lowercase UUID generation, RFC 3339 UTC timestamps, and injectable clock helpers in `server/internal/platform/identity.go`, `server/internal/platform/clock.go`, and `server/internal/platform/platform_test.go`
- [X] T012 Implement SQLite opening with foreign keys, WAL, busy timeout, bounded connections, close optimization, and FTS5 capability checks in `server/internal/store/sqlite/database.go`
- [X] T013 Implement embedded checksummed transactional migration execution and readiness failure behavior in `server/internal/store/sqlite/migrate.go` and `server/internal/store/migrations/embed.go`
- [X] T014 Create the complete constrained schema, required indexes, FTS5 tables/triggers, import staging tables, and schema migration table from `data-model.md` in `server/internal/store/migrations/001_initial.sql`
- [X] T015 [P] Implement the stable API error envelope, field errors, status mapping, and JSON codecs from OpenAPI in `server/internal/api/problem/problem.go` and `server/internal/api/jsoncodec/json.go`
- [X] T016 [P] Implement request IDs, panic recovery, structured request logging, body-size limits, and content-redacting middleware in `server/internal/api/middleware/middleware.go`
- [X] T017 Wire configuration, migrations, `/healthz`, `/readyz`, `/api/v1`, graceful shutdown, and static asset fallback in `server/cmd/noted/main.go` and `server/internal/api/router.go`
- [X] T018 [P] Implement the typed frontend fetch wrapper, error-envelope decoding, request cancellation, and version-conflict handling in `web/src/lib/api/client.ts` and `web/src/lib/api/errors.ts`
- [X] T019 [P] Create current-workspace, query invalidation, toast, unsaved-change, and connectivity stores in `web/src/lib/stores/workspace.ts`, `web/src/lib/stores/query.ts`, `web/src/lib/stores/toast.ts`, and `web/src/lib/stores/connectivity.ts`
- [X] T020 [P] Build the responsive application shell, loading/error boundaries, navigation placeholders, and narrow-view layout in `web/src/routes/+layout.svelte`, `web/src/routes/+error.svelte`, and `web/src/lib/components/AppNavigation.svelte`
- [X] T021 [P] Implement versioned application-shell caching with no API mutation caching or replay in `web/src/service-worker.ts` and `web/static/offline.html`
- [X] T022 Implement OpenAPI parsing and provider request/response validation helpers using `kin-openapi` in `server/internal/api/contracttest/openapi.go` and `server/internal/api/contracttest/openapi_test.go`
- [X] T023 Implement temporary SQLite stores, HTTP test servers, deterministic clocks/IDs, and database integrity assertions in `server/internal/testsupport/testsupport.go`

**Checkpoint**: Both projects build independently; migrations complete; health/readiness and the cached shell work; contract and foundation tests pass. A fresh database with zero workspaces loads a usable setup state and does not leave workspace-scoped routes attempting note or question mutations.

---

## Phase 3: User Story 1 - Capture and Track Questions (Priority: P1) 🎯 MVP

**Goal**: Create the first workspace, then create a note with multiple inline questions and see each canonical question once in the workspace’s active Questions view while offline.

**Independent Test**: Starting with an empty database, open the app, create a named workspace from the setup state, verify it becomes the current workspace, create a note, add two inline questions with different display modes, save, reload, verify both appear in the note and active Questions view, and navigate from a question back to the source note without internet access while the local Go service remains available.

### Tests for User Story 1

- [X] T024 [P] [US1] Add domain tests for required note/question text, default unanswered status, display modes, directive/link integrity, and same-workspace rules in `server/internal/domain/note_test.go` and `server/internal/domain/question_test.go`
- [X] T025 [P] [US1] Add SQLite integration tests for workspace, note, question, and note-question CRUD plus foreign-key and FTS trigger behavior in `server/internal/store/sqlite/capture_test.go`
- [X] T026 [P] [US1] Add OpenAPI handler tests for an empty workspace list, workspace creation, note create/get/update, question create/get/list, validation errors for missing or unknown workspaces, and cursors in `server/internal/api/handlers/capture_contract_test.go`
- [X] T027 [P] [US1] Add editor and active-question component tests for multiple directives and expanded/collapsed/link presentation in `web/tests/component/NoteEditor.test.ts` and `web/tests/component/ActiveQuestions.test.ts`
- [X] T027a [P] [US1] Add workspace setup and empty-state component tests for creating a workspace, selecting the created workspace, recovering from a stale selection, and exposing a workspace-required note CTA in `web/tests/component/WorkspaceSetup.test.ts`
- [X] T028 [P] [US1] Add the offline capture-and-track Playwright journey after a current workspace has been established in `tests/e2e/us1-capture-track.spec.ts`
- [X] T028a [P] [US1] Add a fresh-database workspace bootstrap Playwright journey that opens the setup state, creates a workspace, verifies automatic selection, reaches New note, and confirms no note request is attempted before selection in `tests/e2e/us1-workspace-bootstrap.spec.ts`

### Implementation for User Story 1

- [X] T029 [P] [US1] Define Workspace, Note, Question, NoteQuestion, status, priority, and display-mode domain types and validation in `server/internal/domain/workspace.go`, `server/internal/domain/note.go`, and `server/internal/domain/question.go`
- [X] T030 [P] [US1] Define frontend API and view types for workspaces, notes, links, question summaries, and cursor pages in `web/src/lib/types/workspace.ts`, `web/src/lib/types/note.ts`, and `web/src/lib/types/question.ts`
- [X] T031 [US1] Implement transactional workspace, note, question, and note-question repository operations with optimistic versions, including explicit workspace-existence checks for workspace-owned records, in `server/internal/store/sqlite/capture.go`
- [X] T032 [US1] Implement capture services that validate workspace ownership and Markdown directives, atomically save links, return stable validation and conflict errors for frontend retry handling, and list active questions in `server/internal/service/capture.go`
- [X] T033 [US1] Implement workspace list/create plus note and question create/get/update/list handlers and route registration from `contracts/openapi.yaml` in `server/internal/api/handlers/capture.go` and `server/internal/api/router.go`
- [X] T034 [P] [US1] Implement safe Markdown parsing, raw-HTML rejection, DOMPurify sanitization, and question-directive tokenization in `web/src/lib/editor/markdown.ts` and `web/src/lib/editor/directives.ts`
- [X] T035 [US1] Build the note editor with paragraph/list controls, inline question creation, display-mode controls, retryable save state, and question cards in `web/src/lib/editor/NoteEditor.svelte` and `web/src/lib/components/InlineQuestion.svelte`
- [X] T036 [P] [US1] Implement the workspace creation form, workspace list, selection behavior, and initial empty state in `web/src/routes/workspaces/+page.svelte` and `web/src/lib/api/workspaces.ts`
- [X] T036a [US1] Wire first-run workspace bootstrap and workspace-required route guarding: hydrate the workspace list before enabling workspace-scoped links, select and persist the workspace returned by create, clear stale stored IDs, show a create-workspace CTA when none exists, and prevent `New note` from issuing a request until a valid workspace is selected in `web/src/routes/+layout.svelte`, `web/src/routes/+page.svelte`, `web/src/lib/stores/workspace.ts`, `web/src/lib/components/AppNavigation.svelte`, and `web/src/routes/notes/new/+page.svelte`
- [X] T037 [US1] Implement note creation/editing/loading and linked-question navigation, with an actionable empty/stale-workspace guard that redirects to workspace setup instead of attempting a workspace-less save, in `web/src/routes/notes/new/+page.svelte`, `web/src/routes/notes/[noteId]/+page.svelte`, and `web/src/lib/api/notes.ts`
- [X] T038 [US1] Implement the default active Questions page with canonical deduplication and source-note links in `web/src/routes/questions/+page.svelte`, `web/src/lib/components/QuestionList.svelte`, and `web/src/lib/api/questions.ts`

**Checkpoint**: User Story 1 passes independently and delivers the first usable MVP slice. From an empty database, the user can create and select a workspace before creating a note; a missing or stale selection produces setup guidance and no invalid note mutation.

---

## Phase 4: User Story 2 - Share a Question Across Notes (Priority: P1)

**Goal**: Link an existing question to multiple notes and edit the one canonical question from any note or central view.

**Independent Test**: Link one existing question to two notes, edit its question text from the second note, and verify the first note and central view show the same update without duplicate question records.

### Tests for User Story 2

- [X] T039 [P] [US2] Add service tests for linking an existing question, duplicate-link rejection, cross-workspace rejection, and atomic unlink/reorder behavior in `server/internal/service/shared_question_test.go`
- [X] T040 [P] [US2] Add repository concurrency tests proving one canonical question is returned through multiple note links and stale edits return conflicts in `server/internal/store/sqlite/shared_question_test.go`
- [X] T041 [P] [US2] Add question-picker and linked-note-context component tests in `web/tests/component/QuestionPicker.test.ts` and `web/tests/component/LinkedNotes.test.ts`
- [X] T042 [P] [US2] Add the cross-note linking and synchronized question-text edit Playwright journey in `tests/e2e/us2-shared-question.spec.ts`

### Implementation for User Story 2

- [X] T043 [US2] Extend the note/question repository with existing-question lookup, atomic multi-note link updates, linked-note summaries, and conflict metadata in `server/internal/store/sqlite/shared_question.go`
- [X] T044 [US2] Implement canonical shared-question linking, unlinking, reordering, and edit reconciliation rules in `server/internal/service/shared_question.go`
- [X] T045 [US2] Extend note and question API responses and handlers with linked-note context and stale-version responses in `server/internal/api/handlers/shared_question.go`
- [X] T046 [P] [US2] Build the searchable existing-question picker with enough note/workspace context to distinguish duplicates in `web/src/lib/components/QuestionPicker.svelte`
- [X] T047 [P] [US2] Build linked-note context, navigation, and display-mode controls in `web/src/lib/components/LinkedNotes.svelte`
- [X] T048 [US2] Integrate existing-question insertion and canonical mutation invalidation into note and central question views in `web/src/lib/editor/NoteEditor.svelte` and `web/src/lib/stores/query.ts`

**Checkpoint**: User Story 2 passes independently on top of the foundation and US1 capture primitives.

---

## Phase 5: User Story 3 - Answer and Reopen Questions (Priority: P1)

**Goal**: Enforce question lifecycle rules, separate active and answered filtered views, and preserve answers when reopening.

**Independent Test**: Reject answering without answer text, answer a question, verify it moves from active to Answered view, reopen it into In Progress with its answer preserved, edit it, and prevent an empty answer from remaining Answered.

### Tests for User Story 3

- [X] T049 [P] [US3] Add table-driven domain tests for every allowed and rejected question status/answer transition in `server/internal/domain/question_lifecycle_test.go`
- [X] T050 [P] [US3] Add transactional repository and handler tests for answer updates, active/answered filters, reopen defaults, and stale versions in `server/internal/store/sqlite/lifecycle_test.go` and `server/internal/api/handlers/lifecycle_contract_test.go`
- [X] T051 [P] [US3] Add question answer editor, status control, and answered-list component tests in `web/tests/component/QuestionLifecycle.test.ts`
- [X] T052 [P] [US3] Add answer editing from linked-note and central views, cross-view answer consistency, mark-answered, reopen, revise, and empty-answer rejection Playwright coverage in `tests/e2e/us3-question-lifecycle.spec.ts`

### Implementation for User Story 3

- [X] T053 [US3] Implement authoritative transition validation, answer normalization, and default reopen-to-In-Progress behavior in `server/internal/domain/question_lifecycle.go`
- [X] T054 [US3] Implement atomic lifecycle updates and status-filtered cursor queries in `server/internal/store/sqlite/lifecycle.go`
- [X] T055 [US3] Implement answer/status update behavior and validation error mapping through the existing question endpoint in `server/internal/service/lifecycle.go` and `server/internal/api/handlers/lifecycle.go`
- [X] T056 [P] [US3] Build reusable answer editor, status selector, answer-required feedback, and reopen action in `web/src/lib/components/QuestionLifecycle.svelte`
- [X] T057 [US3] Integrate lifecycle editing into inline, expanded, and central question components in `web/src/lib/components/InlineQuestion.svelte` and `web/src/lib/components/QuestionList.svelte`
- [X] T058 [P] [US3] Implement the Answered Questions filtered route in `web/src/routes/answered/+page.svelte`
- [X] T059 [US3] Ensure status changes invalidate active and answered queries while preserving local answer drafts on failed saves in `web/src/lib/stores/query.ts` and `web/src/lib/stores/unsaved.ts`

**Checkpoint**: All three P1 stories work together; captured questions can be shared, answered, and reopened.

---

## Phase 6: User Story 4 - Prioritize and Find Work (Priority: P2)

**Goal**: Add due dates, reminders, priorities, mandatory filters/sorts, and scoped full-text search that all work without internet access.

**Independent Test**: Seed varied questions and notes, including one question linked to notes with different topics; verify it matches either linked topic without duplication, set due dates/reminders/priorities, exercise every required filter/sort, search notes/questions/answers at current/all-workspace scope, and navigate from each result.

### Tests for User Story 4

- [X] T060 [P] [US4] Add reminder state, due-date, priority, filter validation, and deterministic sort unit tests in `server/internal/domain/prioritization_test.go`
- [X] T061 [P] [US4] Add SQLite FTS synchronization, scoped search, `EXISTS`-based linked-note topic matching without duplicate questions, structured filter, query-plan, cursor, and 50,000-question fixture tests in `server/internal/store/sqlite/search_test.go` and `server/internal/store/sqlite/filter_test.go`
- [X] T062 [P] [US4] Add OpenAPI tests for every question query parameter and `/api/v1/search` scope/result/error shape in `server/internal/api/handlers/search_contract_test.go`
- [X] T063 [P] [US4] Add filter bar, sort control, reminder permission/fallback, overdue state, and search result component tests in `web/tests/component/QuestionDiscovery.test.ts`
- [X] T064 [P] [US4] Add prioritization, reminder, filtering, sorting, scoped search, and result navigation Playwright coverage in `tests/e2e/us4-prioritize-find.spec.ts`

### Implementation for User Story 4

- [X] T065 [P] [US4] Define Reminder plus filter, sort, search-scope, and search-result domain types in `server/internal/domain/reminder.go` and `server/internal/domain/search.go`
- [X] T066 [US4] Implement due date, priority, reminder state, and pending/missed reminder repository operations in `server/internal/store/sqlite/prioritization.go`
- [X] T067 [US4] Implement allow-listed structured question filters, including match-on-any-linked-note topic semantics without duplicate results, deterministic sorting, cursor encoding, and indexed queries in `server/internal/store/sqlite/question_query.go`
- [X] T068 [US4] Implement FTS5 note/question/answer search with workspace/content scopes, escaped snippets, rank, and destinations in `server/internal/search/search.go` and `server/internal/store/sqlite/search.go`
- [X] T069 [US4] Implement prioritization/reminder services and search/filter handlers from OpenAPI in `server/internal/service/prioritization.go`, `server/internal/api/handlers/questions_query.go`, and `server/internal/api/handlers/search.go`
- [X] T070 [P] [US4] Build due-date, priority, overdue, and reminder controls with in-app fallback in `web/src/lib/components/QuestionSchedule.svelte` and `web/src/lib/reminders/reminders.ts`
- [X] T071 [P] [US4] Build all required question filters, sort controls, active-filter summaries, and URL state in `web/src/lib/components/QuestionFilters.svelte` and `web/src/lib/stores/questionQuery.ts`
- [X] T072 [US4] Integrate schedule/filter/sort behavior and cursor pagination into active and answered views in `web/src/lib/components/QuestionList.svelte`
- [X] T073 [P] [US4] Implement scoped search client, search form, result type/workspace/context display, and navigation in `web/src/lib/api/search.ts` and `web/src/routes/search/+page.svelte`
- [X] T074 [US4] Evaluate reminders on startup/resume, request notification permission only from a user action, and surface missed reminders in `web/src/routes/+layout.svelte` and `web/src/lib/reminders/reminders.ts`

**Checkpoint**: User Story 4 independently demonstrates prioritization and discovery at the specified scale without internet access.

---

## Phase 7: User Story 5 - Organize Separate Areas of Life (Priority: P2)

**Goal**: Manage multiple workspaces, note hierarchy, one note topic, multiple tags, and explicit cross-workspace tag availability.

**Independent Test**: Create Work and Studies workspaces, build a valid note tree, assign a topic and tags, verify workspace-isolated default views, explicitly share one tag, and verify unshared tags remain unavailable elsewhere.

### Tests for User Story 5

- [X] T075 [P] [US5] Add hierarchy-cycle, same-workspace parent/topic, tag ownership/access, and unshare-impact domain tests in `server/internal/domain/organization_test.go`
- [X] T076 [P] [US5] Add workspace/topic/tag/hierarchy repository tests including case-insensitive uniqueness and assignment constraints in `server/internal/store/sqlite/organization_test.go`
- [X] T077 [P] [US5] Add OpenAPI tests for workspace/topic/tag CRUD and tag workspace-access replacement in `server/internal/api/handlers/organization_contract_test.go`
- [X] T078 [P] [US5] Add workspace switcher, note tree, topic selector, tag selector, and share confirmation component tests in `web/tests/component/Organization.test.ts`
- [X] T079 [P] [US5] Add multi-workspace, hierarchy, topics, tags, and explicit sharing Playwright coverage in `tests/e2e/us5-workspaces-organization.spec.ts`

### Implementation for User Story 5

- [X] T080 [P] [US5] Define Topic, Tag, TagWorkspaceAccess, NoteTag, and QuestionTag domain types and validation in `server/internal/domain/topic.go` and `server/internal/domain/tag.go`
- [X] T081 [US5] Implement cycle-safe note hierarchy, topic assignment, workspace switching data, tags, assignments, access previews, and unsharing transactions in `server/internal/store/sqlite/organization.go`
- [X] T082 [US5] Implement organization services enforcing workspace boundaries and explicit assignment-removal confirmation in `server/internal/service/organization.go`
- [X] T083 [US5] Implement workspace update, topic CRUD, tag CRUD/access, hierarchy, and organization route handlers from OpenAPI in `server/internal/api/handlers/organization.go`
- [X] T084 [P] [US5] Build the persistent workspace switcher and workspace management page in `web/src/lib/components/WorkspaceSwitcher.svelte` and `web/src/routes/workspaces/+page.svelte`
- [X] T085 [P] [US5] Build accessible note-tree navigation with root/child moves and cycle-error feedback in `web/src/lib/components/NoteTree.svelte`
- [X] T086 [P] [US5] Build topic and tag selectors with owned/shared distinctions in `web/src/lib/components/TopicSelector.svelte` and `web/src/lib/components/TagSelector.svelte`
- [X] T087 [P] [US5] Build tag sharing and unsharing impact review UI in `web/src/lib/components/TagWorkspaceAccess.svelte`
- [X] T088 [US5] Integrate workspace-scoped note trees, topics, tags, questions, and default query resets into `web/src/routes/+layout.svelte`, `web/src/routes/notes/[noteId]/+page.svelte`, and `web/src/lib/stores/workspace.ts`

**Checkpoint**: User Story 5 independently proves context separation, hierarchical organization, and opt-in tag sharing.

---

## Phase 8: User Story 6 - Delete Notes Safely (Priority: P2)

**Goal**: Preview all deletion consequences and atomically preserve, unlink, or explicitly delete affected questions and resolve child notes.

**Independent Test**: Delete a note containing singly and multiply linked questions plus child notes; verify one review shows all effects, chosen singly linked questions remain unlinked or are deleted, multiply linked questions survive, child decisions apply, and cancellation changes nothing.

### Tests for User Story 6

- [X] T089 [P] [US6] Add deletion-impact and command validation tests for every question/child decision and stale preview in `server/internal/service/note_deletion_test.go`
- [X] T090 [P] [US6] Add atomic rollback, cancellation, unlink, keep-unlinked, explicit-delete, and child re-parent integration tests in `server/internal/store/sqlite/note_deletion_test.go`
- [X] T091 [P] [US6] Add deletion preview/execute OpenAPI tests including stale token and version conflicts in `server/internal/api/handlers/note_deletion_contract_test.go`
- [X] T092 [P] [US6] Add single-review dialog component tests for mixed question and child-note consequences in `web/tests/component/DeleteNoteReview.test.ts`
- [X] T093 [P] [US6] Add mixed-link deletion, cancellation, and preservation Playwright coverage in `tests/e2e/us6-safe-delete.spec.ts`

### Implementation for User Story 6

- [X] T094 [US6] Implement cryptographically random opaque short-lived deletion preview tokens backed by server-side preview state, plus complete impact calculation without mutation, in `server/internal/service/note_deletion.go`
- [X] T095 [US6] Implement one-transaction note deletion with version rechecks, explicit per-question decisions, link preservation, and child resolution in `server/internal/store/sqlite/note_deletion.go`
- [X] T096 [US6] Implement deletion preview and execute handlers with stale-preview conflict responses in `server/internal/api/handlers/note_deletion.go`
- [X] T097 [US6] Build and integrate the consolidated deletion review dialog with cancel-safe behavior in `web/src/lib/components/DeleteNoteReview.svelte` and `web/src/routes/notes/[noteId]/+page.svelte`

**Checkpoint**: User Story 6 independently proves no question or child note is silently lost during note deletion.

---

## Phase 9: User Story 7 - Export and Restore All Data (Priority: P3)

**Goal**: Export a complete portable archive and validate, review, and atomically apply an import without silent loss.

**Independent Test**: Export representative nested and unlinked data, import it into an empty state, compare all canonical data/relationships/settings/dates, then verify invalid archives, conflicts, cancellation, expiry, and interrupted apply leave existing data unchanged.

### Tests for User Story 7

- [X] T098 [P] [US7] Add export manifest/schema, deterministic ordering, checksum, ZIP traversal, size/member limit, and malformed archive tests in `server/internal/importexport/archive_test.go`
- [X] T099 [P] [US7] Add full round-trip, duplicate/conflict detection, per-conflict resolution, expiry, rollback, FTS rebuild, and staging cleanup integration tests in `server/internal/store/sqlite/importexport_test.go`
- [X] T100 [P] [US7] Add export, import validate/apply/cancel, payload-limit, and user-safe error OpenAPI tests in `server/internal/api/handlers/importexport_contract_test.go`
- [X] T101 [P] [US7] Add export action, import preview, conflict resolution, invalid archive, and progress component tests in `web/tests/component/DataPortability.test.ts`
- [X] T102 [P] [US7] Add complete round-trip, invalid import, conflict review, cancellation, and rollback Playwright coverage in `tests/e2e/us7-import-export.spec.ts`

### Implementation for User Story 7

- [X] T103 [P] [US7] Define versioned export record DTOs, manifest, checksums, import preview, conflict, and resolution types in `server/internal/importexport/types.go`
- [X] T104 [US7] Implement deterministic streaming ZIP export with normalized JSON records and SHA-256 checksums in `server/internal/importexport/export.go`
- [X] T105 [US7] Implement bounded ZIP validation with normalized-path protection, schema/version/checksum validation, dependency checks, and staging expiry in `server/internal/importexport/validate.go`
- [X] T106 [US7] Implement conflict discovery and one-transaction import apply/rollback with FTS rebuild and session state changes in `server/internal/importexport/apply.go` and `server/internal/store/sqlite/importexport.go`
- [X] T107 [US7] Implement export download and import validate/apply/cancel handlers with cleanup in `server/internal/api/handlers/importexport.go`
- [X] T108 [P] [US7] Implement browser streaming download/upload clients with progress and retry-safe errors in `web/src/lib/api/dataPortability.ts`
- [X] T109 [P] [US7] Build export controls, import validation summary, and conflict decision UI in `web/src/lib/components/ImportReview.svelte` and `web/src/lib/components/ExportData.svelte`
- [X] T110 [US7] Implement the Data settings route and post-import application-state refresh in `web/src/routes/settings/data/+page.svelte` and `web/src/lib/stores/query.ts`

**Checkpoint**: User Story 7 independently proves complete portable round trips and no mutation from failed, cancelled, or unapproved imports.

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Verify performance, accessibility, security, offline reliability, observability, containers, documentation, and all success criteria across completed stories.

- [X] T111 [P] Add keyboard navigation, focus management, accessible names, live error announcements, and contrast fixes for `web/src/lib/editor/NoteEditor.svelte`, `web/src/lib/components/InlineQuestion.svelte`, `web/src/lib/components/QuestionList.svelte`, and `web/src/lib/components/QuestionLifecycle.svelte`, with checks in `tests/e2e/accessibility-notes.spec.ts`
- [X] T112 [P] Add focus trapping/restoration, accessible names, live status announcements, and keyboard behavior for `web/src/lib/components/AppNavigation.svelte`, `web/src/lib/components/WorkspaceSwitcher.svelte`, `web/src/lib/components/NoteTree.svelte`, `web/src/lib/components/DeleteNoteReview.svelte`, and `web/src/lib/components/ImportReview.svelte`, with checks in `tests/e2e/accessibility-organization.spec.ts`
- [X] T113 [P] Add Markdown injection, unsafe directive, oversized input, malformed cursor, and log-redaction security regressions in `server/internal/api/security_test.go` and `web/tests/unit/markdown-security.test.ts`
- [X] T114 [P] Add deterministic 10,000-note/50,000-question data generation and measured search/filter/edit thresholds in `server/internal/testsupport/performance_fixture.go` and `tests/smoke/performance_test.go`
- [X] T115 [P] Add no-internet browser tests that block external requests while retaining localhost API access, plus reload, failed-save draft preservation, retry, and service-worker update checks in `tests/e2e/no-internet.spec.ts`
- [X] T116 [P] Add migration checksum, database integrity, foreign-key, FTS rebuild, restart persistence, and development-only reset checks in `server/internal/store/sqlite/diagnostics_test.go` and `server/cmd/noted/dbcommands.go`
- [X] T117 Add non-root startup, health/readiness, migration, static/API same-origin, persistent-volume restart, and structured-log smoke coverage in `tests/smoke/container.sh`
- [X] T118 [P] Add narrow/wide viewport behavior and empty/loading/error states for `web/src/routes/notes/new/+page.svelte`, `web/src/routes/notes/[noteId]/+page.svelte`, `web/src/routes/questions/+page.svelte`, and `web/src/routes/answered/+page.svelte`, with checks in `tests/e2e/responsive-notes.spec.ts`
- [X] T119 [P] Add narrow/wide viewport behavior and empty/loading/error states for `web/src/routes/workspaces/+page.svelte`, `web/src/routes/search/+page.svelte`, and `web/src/routes/settings/data/+page.svelte`, with checks in `tests/e2e/responsive-secondary.spec.ts`
- [X] T120 Pin and review production/build dependencies, add vulnerability audit commands, and record justified runtime dependencies in `server/go.mod`, `web/package-lock.json`, and `Makefile`
- [X] T121 Create and conduct the usability protocol for SC-001, SC-002, and SC-010, starting participants from an empty database so workspace creation is measured as part of first-attempt note capture, including participant criteria, workflow timing, ease ratings, and anonymized results in `tests/usability/protocol.md` and `tests/usability/results.md`
- [X] T122 Run formatting, static analysis, unit, integration, contract, browser-matrix, performance, and container checks and record results in `specs/001-track-note-questions/validation-results.md`
- [X] T123 Execute every scenario in `quickstart.md`, correct commands and expected outcomes, and update `specs/001-track-note-questions/quickstart.md`
- [X] T124 Trace FR-001 through FR-040 and SC-001 through SC-010 to passing automated or usability evidence and document any approved exception in `specs/001-track-note-questions/traceability.md

**Checkpoint**: All required checks pass from a clean checkout and production-like container with no unresolved constitution exception.

---

## Phase 11: Notes List + Highlight/Annotation Reading (Gap remediation)

**Purpose**: Fix the stub Notes index, stop showing raw `{{question:<uuid>}}` tokens, and make reading interactive: highlight a passage, right-click (or use a selection toolbar) to add a question or annotation, click the highlight to open a card. Unanswered questions still appear under Active Questions; annotations do not.

**New session**: Start with [`PROJECT_MAP.md`](../../PROJECT_MAP.md). Use its file tables, intended capture model, landmines, and “Common tasks → files” instead of rediscovering the tree. Update the map if you change ownership or contracts.

**Why this phase exists**: US1 shipped a textarea editor that appends `{{question:id}}` and a Notes page that never calls `GET /api/v1/notes`. The list API already works. Display modes (expanded/collapsed/link) become legacy once cards ship.

**Independent Test**: Create a workspace and two notes; both appear on `/notes`. Open a note, select a sentence, add a question and an annotation. The passage is highlighted (no raw directive). Clicking it opens a card. The question appears under Active Questions; the annotation does not. `/notes/new` still has labelled Title and Note controls.

### Tests for Phase 11

- [X] T125 [P] [US1] Add unit tests for wrapped highlight directives `{{question:<id>}}passage{{/question}}`, bare legacy tokens, wrap-around-selection, and snippet stripping in `web/tests/unit/directives.test.ts`; keep `web/tests/unit/markdown-security.test.ts` asserting rendered HTML contains no raw `{{question:` token
- [X] T126 [P] [US1] Add NoteReader / AnnotationCard component tests for selection menu actions, opening a card from a highlight click, and question vs annotation chrome in `web/tests/component/NoteReader.test.ts` and `web/tests/component/AnnotationCard.test.ts`
- [X] T127 [P] [US1] Add Notes index tests for loading workspace notes, empty/error states, and navigating to a note in `web/tests/component/NotesIndex.test.ts`
- [X] T128 [P] [US1] Add domain/store tests for question `kind` (`question` default, `annotation`) and Active Questions excluding annotations in `server/internal/domain/question_test.go` and `server/internal/store/sqlite/capture_test.go`
- [X] T129 [P] [US1] Add a Playwright journey: notes appear in the list; highlight capture of one question and one annotation; click highlight opens a card; Active Questions shows only the question — `tests/e2e/us1-notes-highlight.spec.ts`. Keep `tests/e2e/accessibility-notes.spec.ts` Title/Note labels on `/notes/new`.

### Implementation for Phase 11

- [X] T130 [P] [US1] Update product contracts for highlight wraps and `kind`: `specs/001-track-note-questions/spec.md` (US1, FR-010, annotations, notes index), `plan.md` editing section, `data-model.md`, `contracts/note-content.md`, `contracts/openapi.yaml`; copy contracts to `tests/fixtures/contracts/`
- [X] T131 [P] [US1] Add `kind` (`question` | `annotation`, default `question`) via `server/internal/store/migrations/002_question_kind.sql`, `server/internal/domain/question.go`, SQL in `server/internal/store/sqlite/capture.go`, list filter in `server/internal/store/sqlite/question_query.go` and `server/internal/api/handlers/questions_query.go`, preserve kind in `server/internal/store/sqlite/lifecycle.go`, allow-list `kind` in `server/internal/importexport/apply.go`
- [X] T132 [US1] Implement the Notes index: load `notesApi.list` for the current workspace, show notes (reuse or restyle `web/src/lib/components/NoteTree.svelte`), empty/loading/error and workspace CTA in `web/src/routes/notes/+page.svelte`
- [X] T133 [US1] Extend `web/src/lib/editor/directives.ts` and `web/src/lib/editor/markdown.ts` to tokenize wrapped highlights, wrap a selected passage, and render sanitized HTML with `<mark data-annotation-id>` (bare tokens become a marker chip, never visible `{{question:uuid}}`)
- [X] T134 [US1] Build reading-view capture UI: `web/src/lib/editor/NoteReader.svelte` (selection, context menu, selection toolbar) and `web/src/lib/components/AnnotationCard.svelte` (passage quote, question lifecycle or annotation editor). Wire read/edit modes in `web/src/lib/editor/NoteEditor.svelte` and `web/src/routes/notes/[noteId]/+page.svelte`. Existing notes default to read; new notes stay in edit so Title/Note labels remain.
- [X] T135 [US1] Create question/annotation from the composer (`POST /api/v1/questions` with `kind`), wrap the passage, rebuild `questionLinks` from directive order, save the note, open the card. Filter Active and Answered lists to `kind=question` in `web/src/lib/types/question.ts`, `web/src/lib/api/questions.ts`, `web/src/routes/questions/+page.svelte`, and `web/src/routes/answered/+page.svelte`.
- [X] T136 [P] [US1] After the slice ships, refresh [`PROJECT_MAP.md`](../../PROJECT_MAP.md) so the “known bugs / in-progress” table matches reality.

**Checkpoint**: From a selected workspace, created notes appear under Notes. Reading a note never shows raw question directives. A highlighted passage can receive a question or an annotation; the highlight opens a card; unanswered questions remain in Active Questions; annotations do not. `/notes/new` accessibility labels still pass.

### [X] T137 Finish remaining Phase 11 (handoff — completed)

**For a small model.** Do not scan the repo. Do not reread US1–US10. Read `PROJECT_MAP.md` once, then only the files named below. Keep the session small: implement, run the listed tests, stop.

**Implementation status**: The backend kind support, contracts, reader wiring, and map updates described below are landed.

#### Already present — do not rewrite

| Area | Files |
|---|---|
| Tests | `web/tests/unit/directives.test.ts`, `web/tests/unit/markdown-security.test.ts` (no raw `{{question:`), `web/tests/component/NoteReader.test.ts`, `web/tests/component/AnnotationCard.test.ts`, `web/tests/component/NotesIndex.test.ts`, `tests/e2e/us1-notes-highlight.spec.ts`, `TestQuestionKindDefaultsAndValidation` in `server/internal/domain/question_test.go`, kind/list/lifecycle cases in `server/internal/store/sqlite/capture_test.go` |
| Notes list (T132) | `web/src/routes/notes/+page.svelte` already calls `notesApi.list` and renders `NoteTree` |
| Wrap/render (T133) | `web/src/lib/editor/directives.ts` (`wrapSelection`, wrapped tokenize, `stripDirectives`), `web/src/lib/editor/markdown.ts` (`renderNoteHtml` → `<mark data-annotation-id>`) |
| Reader/card UI | `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/AnnotationCard.svelte` |
| Frontend types/client | `web/src/lib/types/question.ts` (`QuestionKind`), `web/src/lib/api/questions.ts` (`kind` query + create body) |
| Migration | `server/internal/store/migrations/002_question_kind.sql` (one `ALTER TABLE`); `embed.go` already embeds `*.sql` |

#### Completed implementation scope

**1. Backend `kind` (T128/T131)** so existing Go tests compile.

- `server/internal/domain/question.go`: add `QuestionKind`, `KindQuestion = "question"`, `KindAnnotation = "annotation"`, `ValidKind`. Add `Kind` to `QuestionSummary`, `Question`, `QuestionWrite` (`json:"kind"`). In `Validate()`, default empty kind to `KindQuestion` and reject anything else.
- `server/internal/domain/prioritization.go`: add `Kind QuestionKind` to `QuestionQuery`. In `NormalizeQuestionQuery`, if `Kind != ""` and not valid, return `invalid question kind`.
- `server/internal/store/sqlite/capture.go`:
  - `CreateQuestion`: copy `write.Kind` onto the domain value; `INSERT` includes `kind`.
  - `getQuestionBase`: `SELECT` + `Scan` `kind`.
  - `UpdateQuestion`: set `q.Kind = old.Kind` (immutable). **Do not** put `kind` in the `UPDATE` SET list (that is how lifecycle preserves it).
  - `noteLinks`: `SELECT q.kind` and scan into `l.Question.Kind`.
- `server/internal/store/sqlite/question_query.go`: if `query.Kind != ""`, add `q.kind = ?`.
- `server/internal/api/handlers/questions_query.go`: `parseQuestionQuery` reads `kind` query param into `query.Kind`.
- `server/internal/store/sqlite/lifecycle.go`: pass `Kind: q.Kind` in the `QuestionWrite` to `UpdateQuestion`.
- `server/internal/importexport/apply.go`: add `"kind": true` to the `questions` allow-list.
- `server/internal/store/sqlite/note_deletion_impact.go`: set `Kind: q.Kind` on the `QuestionSummary` literal (otherwise JSON encodes `"kind":""` and OpenAPI enum fails).

JSON decoder rejects unknown fields. `kind` must exist on the Go structs **and** OpenAPI before any client sends it.

**2. Contracts (T130)** — small surgical edits only.

- `specs/001-track-note-questions/contracts/openapi.yaml`: add `QuestionKind` enum `[question, annotation]`; optional `kind` on `QuestionWrite`; `kind` on `QuestionSummary`; list query param `kind`. Then copy to `tests/fixtures/contracts/openapi.yaml`.
- `specs/001-track-note-questions/contracts/note-content.md`: wrapped form `{{question:<uuid>}}passage{{/question}}`; bare `{{question:<uuid>}}` still valid; reader renders `<mark>`, never visible tokens; `{{/question}}` is not a second directive. Copy to `tests/fixtures/contracts/note-content.md`.
- `spec.md`: US1 = notes index + highlight capture + card; annotations excluded from Active Questions; rewrite **FR-010** to highlight+card (legacy display modes may stay stored); remove “highlights” from Out of Scope / Assumptions deferred list.
- `plan.md` “Editing and content”: wrapped directive + reading view.
- `data-model.md` `questions` table: `kind` text NOT NULL default `question` CHECK (`question`,`annotation`).

**3. Wire read/capture (T134/T135)** — do not rebuild `NoteReader` / `AnnotationCard`.

- `web/src/lib/editor/NoteEditor.svelte`:
  - `mode = existing ? 'read' : 'edit'`.
  - Read: `NoteReader` with `markdown`, `questions`, `onCapture`.
  - Edit: keep current form so `/notes/new` still has labelled **Title** and **Note** (`tests/e2e/accessibility-notes.spec.ts`). Toggle “Edit note” / “Reading view” only when `existing` is set.
  - `captureFromReader(kind, passage, text)`: `questionsApi.create({ ..., kind })`, `wrapSelection(bodyMarkdown, passage, id)`, rebuild `questionLinks` from `directiveIds(body)` (same order, `displayMode` default `collapsed`), save note, return the created question (so the reader can open the card).
  - After save, assign `existing = note` so `version` stays current.
  - Hydrate full questions with `questionsApi.get` for each link when `existing` is set (summaries lack `answerMarkdown` / `version` for `QuestionLifecycle`).
  - Question picker list: pass `kind: 'question'`.
- `web/src/lib/editor/NoteReader.svelte`: if `onCapture` returns a question, set `openQuestion` / `openPassage` after save.
- `web/src/routes/notes/[noteId]/+page.svelte`: already hosts `NoteEditor`; no extra mode flag if the editor defaults existing notes to read.
- `web/src/routes/questions/+page.svelte` and `web/src/routes/answered/+page.svelte`: add `kind: 'question'` to `questionsApi.list`.

**4. Map + checkboxes (T136)**

- `PROJECT_MAP.md`: known-bugs table — Notes list and highlight/card are done; remaining gaps only if tests still fail.
- This file: mark T125–T137 `[X]` only after tests pass.

#### Verify (targeted — do not run the whole matrix first)

```
go -C server test ./internal/domain ./internal/store/sqlite -count=1
npm --prefix web test -- --run tests/unit/directives.test.ts tests/unit/markdown-security.test.ts tests/component/NoteReader.test.ts tests/component/AnnotationCard.test.ts tests/component/NotesIndex.test.ts tests/component/NoteEditor.test.ts
```

Then `tests/e2e/us1-notes-highlight.spec.ts` and `tests/e2e/accessibility-notes.spec.ts`.

#### Landmines

- Do not change `/notes/new` labels.
- `questionLinks[i]` order must match directive order or note save 422s.
- Selection wrap is exact substring of `bodyMarkdown`; e2e uses plain text.
- Kind is immutable after create; empty `kind` on update must not turn annotations into questions.
- Do not log note/question/answer bodies.

---

## Phase 12: Note editing, capture dismissal, status filters, and tag discovery (gap remediation)

**Purpose**: Make existing-note editing reliable, make selection actions dismissible, improve the Active Questions status filter, and add note tagging plus tag-aware note search without regressing inline question capture.

**Independent Test**: Create a note, reopen it, enter edit mode, append text and add an inline question, save and reload, and verify both persist. In reading mode, select text and dismiss the selection actions without a mutation, then select again and capture a highlight. Assign two workspace tags to the note, search notes by each tag, and verify only matching notes appear. Exercise the redesigned Active Questions status filter at narrow and wide viewports.

- [X] T138 [US1] [US4] [US5] Close the note workflow gaps: make existing notes genuinely editable with bound Title/Note controls and current-version saves while retaining inline question creation and picker controls; add Cancel/Escape/outside-click dismissal for selection actions and the composer; replace the ugly native multi-select status control with an accessible compact multi-select/chip or checkbox filter that preserves active query semantics; add a note-tags section using available owned/shared workspace tags and persist `tagIds`; add a tag filter to note search (including OpenAPI/fixture, client, handler, indexed `EXISTS` filtering without duplicate results, URL state, empty/error states, and responsive UI). Add tests first in `web/tests/component/NoteEditor.test.ts`, `web/tests/component/NoteReader.test.ts`, `web/tests/component/QuestionDiscovery.test.ts`, new `web/tests/component/TagSelector.test.ts` and search coverage, `server/internal/store/sqlite/search_test.go`, `server/internal/api/handlers/search_contract_test.go`, and `tests/e2e/t138-note-edit-tags-filter.spec.ts`; preserve the `/notes/new` Title/Note accessibility journey in `tests/e2e/accessibility-notes.spec.ts`. Implement across `web/src/lib/editor/NoteEditor.svelte`, `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/QuestionFilters.svelte`, `web/src/lib/components/TagSelector.svelte`, `web/src/lib/api/tags.ts`, `web/src/lib/api/search.ts`, `web/src/routes/search/+page.svelte`, `server/internal/domain/search.go`, `server/internal/store/sqlite/search.go`, `server/internal/api/handlers/search.go`, `specs/001-track-note-questions/contracts/openapi.yaml`, and `tests/fixtures/contracts/openapi.yaml`.

**T138 acceptance checkpoint**: Existing notes can be edited and saved/reloaded; edit mode still exposes inline question creation and maintains directive/link order; selection actions and the composer can be cancelled with no write; Active Questions status filtering is readable, keyboard accessible, and responsive; notes can be tagged from the editor; searching notes with a tag returns matching notes only, does not duplicate notes, and supports clearing the filter. No raw directives or workspace-boundary violations are introduced.

---

## Phase 13: Answer in context (US8) — gap remediation

**Purpose**: Opening a question shows the source note and highlighted passage on the same page as answer/status controls.

**For a small model.** Do not scan the repo. Do not reread US1–US7 or Phase 11 internals. Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement **one task**, run its verify command, stop.

**Independent Test**: Capture a question on a sentence. From Active Questions open it. The question page shows the note, the sentence is highlighted and in view, and the answer box is on the same page. An unlinked question still opens and says it is unlinked.

**Do not**: change capture, kind, lifecycle rules, NoteReader toolbar, `/notes/new` labels, or add APIs.

### Tests for Phase 13 (write first; they must fail)

- [X] T139 [P] [US8] Add `web/tests/component/NoteExcerpt.test.ts`: render markdown with a wrapped directive; with `questionId` set, the `mark[data-annotation-id]` exists; no “Ask a question” / “Add annotation” buttons. Add `web/tests/component/QuestionContext.test.ts`: linked question shows note title and `QuestionLifecycle`; unlinked question shows text “currently unlinked” and still shows `QuestionLifecycle`; two `linkedNotes` render a control to switch notes.
- [X] T140 [P] [US8] Add Playwright `tests/e2e/us8-answer-in-context.spec.ts`: create workspace + note + one highlight question; visit Active Questions; click the question; assert the source passage text is visible and an Answer field is visible on the same page. Keep `tests/e2e/accessibility-notes.spec.ts` passing.

### Implementation for Phase 13

- [ ] T141 [P] [US8] Create `web/src/lib/components/NoteExcerpt.svelte`.
  **Read:** `web/src/lib/editor/markdown.ts` (`renderNoteHtml`).
  **Do:** `export let markdown = ''`; `export let questionId = ''`. Render sanitized HTML in `<article aria-label="Source note">`. After markdown/questionId change, `document.querySelector('mark[data-annotation-id="' + questionId + '"]')?.scrollIntoView({ block: 'center' })`.
  **Do not:** import `NoteReader.svelte`; no toolbar, composer, or capture.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteExcerpt.test.ts`

- [ ] T142 [US8] Create `web/src/lib/components/QuestionContext.svelte`.
  **Read:** `web/src/lib/types/question.ts`, `web/src/lib/types/note.ts`, `web/src/lib/components/QuestionLifecycle.svelte`, `web/src/lib/components/QuestionSchedule.svelte`, `web/src/lib/api/notes.ts`.
  **Do:** Props: `question: Question`, optional `note: Note | null`, `selectedNoteId: string`, `onSelectNote: (id: string) => void`, `onSave: (q: Question) => void`. Layout: note pane (or “This question is currently unlinked.”) + `QuestionLifecycle` + `QuestionSchedule`. If `question.linkedNotes.length > 1`, render a `<label for="context-note">Source note</label>` `<select id="context-note">` of linked notes. Wide: two columns. Narrow (`max-width: 800px`): stack note above controls. Use `export let` (not runes).
  **Do not:** fetch inside this component; the page fetches.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionContext.test.ts`

- [ ] T143 [US8] Wire `web/src/routes/questions/[questionId]/+page.svelte`.
  **Read:** that page, `web/src/lib/api/questions.ts`, `web/src/lib/api/notes.ts`, `$app/state` `page`.
  **Do:** Load question with `questionsApi.get`. `noteId` = `page.url.searchParams.get('noteId')` if it is in `question.linkedNotes`, else `question.linkedNotes[0]?.id`. If `noteId`, load `notesApi.get`. Render `QuestionContext`. Switching notes reloads that note. Keep loading/error/retry patterns already on the page.
  **Do not:** add new API fields; do not log question/note bodies.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionContext.test.ts tests/component/NoteExcerpt.test.ts`

- [ ] T144 [US8] Point list cards at context.
  **Read:** `web/src/lib/components/QuestionList.svelte`.
  **Do:** Question title link becomes `/questions/{id}?noteId={firstLinkedNoteId}` when `linkedNotes[0]` exists, else `/questions/{id}`. Accessible name stays the question text.
  **Verify:** `npx playwright test tests/e2e/us8-answer-in-context.spec.ts tests/e2e/us3-question-lifecycle.spec.ts`

**Landmines:** Svelte 5 but match `export let`. JSON decoder rejects unknown fields — send no new keys. `questionLinks` order is irrelevant here (do not save notes). Do not use `NoteReader` on this page (it would offer capture).

**Checkpoint**: Active Questions → question page shows passage + answer controls together. Unlinked questions still open.

---

## Phase 14: Note-local open questions + next gap (US8) — gap remediation

**Purpose**: While reading a note, see that note’s open questions and jump to the next unanswered highlight.

**For a small model.** Read this phase only. One task at a time. Do not rebuild `NoteReader` capture.

**Independent Test**: Two questions on different sentences in one note. Reading view lists both. Next unanswered opens the first card and scrolls to it; again opens the second; again announces none remain. Annotations and answered questions are absent from the rail.

**Do not**: put annotations or answered items in the rail; change Active Questions; add APIs.

### Tests for Phase 14 (write first; they must fail)

- [ ] T145 [P] [US8] Add `web/tests/component/NoteQuestionRail.test.ts`: given mixed questions, the rail lists only `kind=question` with status `unanswered` / `in_progress` / `deferred`, in the provided `orderedIds` order; clicking an item calls `onSelect(id)`; Next unanswered calls `onNext`; when `remaining === 0` the next control is disabled or the rail text is `No open questions in this note`.
- [ ] T146 [P] [US8] Extend `web/tests/component/NoteReader.test.ts`: when `focusQuestionId` is set to a wrapped id, the card opens. Extend `tests/e2e/us8-answer-in-context.spec.ts` (or add `tests/e2e/us8-note-rail.spec.ts`) with two highlight questions, Next unanswered twice, then the empty message.

### Implementation for Phase 14

- [ ] T147 [P] [US8] Create `web/src/lib/components/NoteQuestionRail.svelte`.
  **Do:** `export let questions: Question[] = []`; `export let orderedIds: string[] = []`; `export let onSelect: ((id: string) => void) | undefined`; `export let onNext: (() => void) | undefined`. Filter `kind !== 'annotation'` (missing kind counts as question) and `status !== 'answered'`. Sort by `orderedIds`. Render `<aside aria-label="Open questions in this note">` with a button per item (question text + status) and a “Next unanswered” button. Empty: `No open questions in this note`.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteQuestionRail.test.ts`

- [ ] T148 [US8] Open a highlight from the outside.
  **Read:** `web/src/lib/editor/NoteReader.svelte`.
  **Do:** Add `export let focusQuestionId: string | null = null`. When `focusQuestionId` changes to a non-null id present in `questions`, set `openQuestion` / `openPassage` the same way highlight click does (use `tokenizeDirectives` snippet). Do not clear `focusQuestionId` yourself if it is a prop; parent may reset it.
  **Do not:** remove Cancel / Escape / outside-click dismissal.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [ ] T149 [US8] Wire the rail in `web/src/lib/editor/NoteEditor.svelte` read mode.
  **Read:** `NoteEditor.svelte` around the `NoteReader` usage; `web/src/lib/editor/directives.ts` `directiveIds`.
  **Do:** Beside `NoteReader` (not inside the article), render `NoteQuestionRail` with `questions` and `orderedIds={directiveIds(bodyMarkdown)}`. Keep `focusId` in the editor. Rail `onSelect` / `onNext` set `focusId` to the chosen / next open id in directive order (skip annotation + answered). After the last item, set `focusId = null` (rail empty message handles UX). Pass `focusQuestionId={focusId}` to `NoteReader`. Narrow: stack rail below the reader.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts tests/component/NoteQuestionRail.test.ts tests/component/NoteReader.test.ts` then the Phase 14 Playwright file.

**Landmines:** Directive order, not createdAt. Annotations never in the rail. `/notes/new` stays edit mode with Title/Note labels — do not show the rail on new unsaved notes.

**Checkpoint**: Reading a note, open questions are listed in passage order and Next unanswered walks them.

---

## Phase 15: Today / Next queue (US9) — gap remediation

**Purpose**: One opinionated workspace queue. Not saved filters.

**For a small model.** Read this phase only. Pure helper first, then page, then nav. No new backend endpoint.

**Independent Test**: Seed overdue, due today, in progress, deferred with due ≤ today, high priority unanswered, deferred with future due, annotation, answered. Next shows five sections in spec order; omits the last three kinds of items; opening a row lands on the US8 question page.

**Do not**: add saved filters, query builders, new API routes, or new DB columns.

### Tests for Phase 15 (write first; they must fail)

- [ ] T150 [P] [US9] Add `web/tests/unit/nextQueue.test.ts` covering `buildNextQueue` in `web/src/lib/questions/nextQueue.ts` with a fixed `today` of `2026-04-01`:
  - overdue: `dueDate < today` and status not `deferred`
  - due today: `dueDate === today` and status not `deferred`
  - in progress: `in_progress` not already in overdue/due today
  - deferred ready: `deferred` and `dueDate <= today`
  - high priority: `high` or `urgent`, status `unanswered` or `in_progress`, not already listed
  - omit: `answered`, `kind=annotation`, `deferred` with `dueDate > today`, `deferred` with null due date
  - hide empty sections; never duplicate an id across sections
- [ ] T151 [P] [US9] Add `web/tests/component/NextQueue.test.ts` for empty guidance and section headings. Add `tests/e2e/us9-next-queue.spec.ts`: one overdue question appears under Overdue; an answered question does not; clicking opens `/questions/{id}`.

### Implementation for Phase 15

- [ ] T152 [P] [US9] Create `web/src/lib/questions/nextQueue.ts`.
  **Do:** Export `buildNextQueue(questions: Question[], today: string): { id: string; title: string; items: Question[] }[]`. `today` is `YYYY-MM-DD`. Filter to active questions only (`kind` missing or `question`; status `unanswered` | `in_progress` | `deferred`). Apply FR-047. Return only non-empty sections with titles `Overdue`, `Due today`, `In progress`, `Deferred ready`, `High priority`.
  **Verify:** `npm --prefix web test -- --run tests/unit/nextQueue.test.ts`

- [ ] T153 [US9] Create `web/src/routes/next/+page.svelte`.
  **Read:** `web/src/routes/questions/+page.svelte` for workspace loading/error/retry; `web/src/lib/api/questions.ts`; `web/src/lib/stores/workspace.ts`.
  **Do:** Require current workspace. `questionsApi.list({ workspaceId, status: ['unanswered','in_progress','deferred'], kind: 'question', pageSize: 200 })`. Follow `nextCursor` up to 5 pages. `today = new Date().toISOString().slice(0, 10)`. Render sections from `buildNextQueue`. Each item links like Phase 13 (`/questions/{id}?noteId=...`). Empty page: `Nothing in Next. Capture a question from a note, or set a due date.` If `nextCursor` remains after 5 pages, show `Showing the first loaded questions.` Loading/error/retry required. No filter widgets.
  **Verify:** `npm --prefix web test -- --run tests/unit/nextQueue.test.ts tests/component/NextQueue.test.ts`

- [ ] T154 [US9] Add Next to `web/src/lib/components/AppNavigation.svelte` after Active Questions: `href={`/next${workspace}`}` labelled `Next`, `class:active={page.url.pathname === '/next'}` (do not use `startsWith` if it would clash). Only when `$workspaceReady`.
  **Verify:** `npx playwright test tests/e2e/us9-next-queue.spec.ts tests/e2e/us8-answer-in-context.spec.ts`

**Landmines:** Client-side grouping only. Do not add `queue=` to OpenAPI. Annotations never appear (`kind=question` on the list call). Do not log question text.

**Checkpoint**: `/next` is a usable daily queue that opens answer-in-context.

---

## Phase 16: Deferred requires a resume date (US10) — gap remediation

**Purpose**: Status `deferred` requires a due date. No new status. Legacy null due dates still load.

**For a small model.** Domain rule first, then every `ValidateTransition` caller, then UI. Do not migrate old rows.

**Independent Test**: Save as Deferred with no date → rejected. With a date → saved. Next omits future-deferred and includes deferred-ready. Opening a legacy deferred question with null due date still works.

### Tests for Phase 16 (write first; they must fail)

- [ ] T155 [P] [US10] Extend `server/internal/domain/question_lifecycle_test.go`: `deferred` with empty/nil due date fails; `deferred` with `YYYY-MM-DD` succeeds; `answered` still requires an answer and does not require a due date; `unanswered` / `in_progress` still allow null due date.
- [ ] T156 [P] [US10] Extend `web/tests/component/QuestionLifecycle.test.ts`: choosing Deferred with no date shows `A resume date is required to defer a question.` and does not call update; with a date, save sends `status: 'deferred'` and that `dueDate`. Extend `tests/e2e/us3-question-lifecycle.spec.ts` or add `tests/e2e/us10-defer-date.spec.ts` for the rejection + success path.

### Implementation for Phase 16

- [ ] T157 [US10] Update `server/internal/domain/question_lifecycle.go`.
  **Read:** `ValidateTransition` and every caller (`rg ValidateTransition server`).
  **Do:** `ValidateTransition(from, to, answer, dueDate *string)`. If `to == StatusDeferred`, trimmed due date must be non-empty `YYYY-MM-DD`; else return `a resume date is required to defer a question`. Do **not** reject existing rows on read. Update all callers so compile succeeds. Keep answer-required-for-answered behavior.
  **Do not:** add a column; do not rewrite import of legacy deferred rows to fail.
  **Also edit:** `specs/001-track-note-questions/data-model.md` lifecycle table if the deferred-due-date row is missing.
  **Verify:** `go -C server test ./internal/domain -count=1`

- [ ] T158 [US10] Map the domain error through `server/internal/service/lifecycle.go` and `server/internal/api/handlers/lifecycle.go` (or the existing question update handler) to a 422 field error on `dueDate` / `status`. Preserve optimistic `version` conflicts.
  **Verify:** `go -C server test ./internal/domain ./internal/store/sqlite ./internal/api/handlers -count=1`

- [ ] T159 [US10] Update `web/src/lib/components/QuestionLifecycle.svelte`.
  **Read:** that file and `web/src/lib/components/QuestionSchedule.svelte` (due date already exists).
  **Do:** When `status === 'deferred'`, show `<label for="resume-date">Resume date</label> <input id="resume-date" type="date">` bound to a local due date (initialize from `question.dueDate`). On save, if deferred and no date, set the error string above and return without calling the API. Include `dueDate` on `questionsApi.update`. Do not apply this to annotations (`question.kind === 'annotation'` skips the rule).
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts` then the Phase 16 Playwright file.

**Landmines:** Read of null `due_date` must still work. Do not auto-write a date on load. JSON unknown fields 400. Do not log answer text.

**Checkpoint**: Deferred without a date cannot be saved; Next can wake deferred items by date.

---

## Phase 17: Resolved highlights + insert answer into note (US11) — gap remediation

**Purpose**: Answered wraps look resolved. Optional, idempotent “Insert answer into note”. Answering does not mutate the note by itself.

**For a small model.** Helper first, then render attributes, then a button on the card. Do not auto-insert on save.

**Independent Test**: Answer a question → highlight style changes, note body unchanged. Insert → blockquote after wrap. Insert again → body unchanged. Reopen → active highlight style; blockquote remains.

### Tests for Phase 17 (write first; they must fail)

- [ ] T160 [P] [US11] Add `web/tests/unit/insertAnswer.test.ts` for `insertAnswerAfterDirective` in `web/src/lib/editor/directives.ts`: wrapped id inserts `\n\n` + blockquote of the answer immediately after `{{/question}}`; multiline answer prefixes each line with `> `; second call returns the same markdown; unknown id returns markdown unchanged; empty answer throws and does not mutate.
- [ ] T161 [P] [US11] Extend `web/tests/unit/markdown-security.test.ts` or add `web/tests/unit/markdown-status.test.ts`: `renderNoteHtml(md, questions)` sets `data-status` and `data-kind` on the mark; raw `{{question:` still absent. Extend `web/tests/component/AnnotationCard.test.ts`: answered question shows Insert answer into note; clicking calls `onInsertAnswer` once.
- [ ] T162 [P] [US11] Add `tests/e2e/us11-resolved-highlight.spec.ts`: answer, reload note, mark has resolved styling or `data-status="answered"`; insert once; reload; blockquote visible; answer save without insert does not add a blockquote.

### Implementation for Phase 17

- [ ] T163 [P] [US11] Add `insertAnswerAfterDirective(markdown: string, id: string, answer: string): string` to `web/src/lib/editor/directives.ts`.
  **Do:** Trim answer; throw if empty. Find the wrapped directive with that lowercase id. If the text after that wrap (skip one run of whitespace) already starts with a `>` line that contains the trimmed answer, return markdown unchanged. Else splice `\n\n` + answer lines each prefixed with `> ` + `\n` immediately after `{{/question}}`. Bare (unwrapped) tokens: insert the same blockquote immediately after the opening token. Reconstruct via `tokenizeDirectives` or index math; do not invent a second wrap.
  **Verify:** `npm --prefix web test -- --run tests/unit/insertAnswer.test.ts`

- [ ] T164 [US11] Teach `renderNoteHtml` statuses.
  **Read:** `web/src/lib/editor/markdown.ts`.
  **Do:** `renderNoteHtml(markdown: string, questions: { id: string; status?: string; kind?: string }[] = [])`. On `<mark>` add `data-status` (default `unanswered`) and `data-kind` (default `question`). Add `data-status` and `data-kind` to DOMPurify `ADD_ATTR`. Existing callers still work with one argument.
  **CSS:** In `NoteReader.svelte` and `NoteExcerpt.svelte`: `mark[data-status="answered"]` background `#d1fae5`, border-bottom `#047857`. Leave active marks as they are.
  **Do:** Pass `questions` into `renderNoteHtml` from `NoteReader` (it already has `questions`). Pass question status into `NoteExcerpt` (`export let status = ''` or a questions array).
  **Verify:** `npm --prefix web test -- --run tests/unit/markdown-security.test.ts tests/component/NoteReader.test.ts tests/component/NoteExcerpt.test.ts`

- [ ] T165 [US11] Insert from the card, user-initiated.
  **Read:** `web/src/lib/components/AnnotationCard.svelte`, `web/src/lib/editor/NoteEditor.svelte`, `web/src/lib/components/QuestionLifecycle.svelte`.
  **Do:** If `question.kind !== 'annotation'` and `question.status === 'answered'` and `question.answerMarkdown`, show button `Insert answer into note`. `export let onInsertAnswer: (() => Promise<void> | void) | undefined`. In `NoteEditor` (read mode owns the body): implement `insertAnswer(question)` using `insertAnswerAfterDirective`, rebuild `questionLinks` from `directiveIds` (order must still match), `save()`. Do not change question status. Wire the callback through `NoteReader` → `AnnotationCard`. If insert is clicked on the question context page, apply it to the loaded note via `notesApi.update` with current `version` and rebuilt links from that note’s body; skip if the question is unlinked.
  **Do not:** insert inside `QuestionLifecycle` save.
  **Verify:** component tests then `npx playwright test tests/e2e/us11-resolved-highlight.spec.ts tests/e2e/us8-answer-in-context.spec.ts`

**Landmines:** `questionLinks` must match directive order after insert (blockquote is not a directive). Idempotent insert. Never log the answer. `/notes/new` labels unchanged.

**Checkpoint**: Answered passages look resolved; insert is optional and safe to click twice.

---

## Phase 18: Sturdier passage matching (US12) — gap remediation

**Purpose**: Wrap the source Markdown span that produced the rendered selection. Fail with no mutation when it cannot be mapped.

**For a small model.** Pure function in `directives.ts` only, then switch `wrapSelection` to it, then show the error in the composer. No schema.

**Independent Test**: Note body `This is **bold** text.` Select rendered `bold`. Capture succeeds and wrap includes `**bold**`. Select `zzz` → error, no new question, body unchanged.

### Tests for Phase 18 (write first; they must fail)

- [ ] T166 [P] [US12] Extend `web/tests/unit/directives.test.ts`:
  - exact match still wraps the first occurrence
  - markdown `**bold**` + selection `bold` wraps `**bold**`
  - markdown `hello   world` + selection `hello world` wraps the original spaced span
  - no match throws `Could not find that passage in the note. Try selecting plain text.`
  - empty selection throws the existing empty error
  - wrap still produces `{{question:id}}…{{/question}}` and `directiveIds` order is preserved when a wrap is added among existing directives
- [ ] T167 [P] [US12] Extend `web/tests/component/NoteReader.test.ts`: when `onCapture` rejects with that error, composer shows it and no success path runs. Playwright: `tests/e2e/us12-passage-match.spec.ts` for the bold case (use a note whose body is exactly `This is **bold** text.`).

### Implementation for Phase 18

- [ ] T168 [US12] Add `findSelectionInMarkdown(markdown: string, selectedText: string): { index: number; length: number }` in `web/src/lib/editor/directives.ts`.
  **Do, in order:**
  1. Trimmed selection empty → throw `selection is empty`.
  2. `markdown.indexOf(selectedText)` ≥ 0 → `{ index, length: selectedText.length }`.
  3. Collapse whitespace (runs of whitespace → one space) on both strings with an index map back to markdown; if the collapsed selection occurs, return the mapped `[start, end)` in the original markdown.
  4. Build a readable string from markdown by skipping `*`, `_`, and `` ` `` and collapsing whitespace, mapping each readable index to a markdown index; find collapsed selection there; return mapped `[start, end)` so markers stay inside the wrap (`**bold**` not `bold`).
  5. Else throw `Could not find that passage in the note. Try selecting plain text.`
  First match wins. Then change `wrapSelection` to slice `markdown[index, index+length]` as the wrapped span.
  **Do not:** parse links, headings, or HTML. Do not add DB columns.
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts`

- [ ] T169 [US12] Surface the error in `web/src/lib/editor/NoteEditor.svelte` `captureFromReader` and `web/src/lib/editor/NoteReader.svelte` `submitComposer` (composerError already exists). If wrap throws, do not `questionsApi.create` first — **reorder capture** so wrap is resolved **before** create, or delete the created question is NOT allowed; instead compute wrap on a copy first, then create, then save. Preferred order: `findSelectionInMarkdown` → `questionsApi.create` → wrap with returned id → save. If save fails, keep existing failed-save draft behavior.
  **Verify:** unit + `npx playwright test tests/e2e/us12-passage-match.spec.ts tests/e2e/us1-notes-highlight.spec.ts`

**Landmines:** Create-after-failed-wrap would orphan questions — match before create. `questionLinks` order. Do not log the passage.

**Checkpoint**: Formatted-text selections wrap source Markdown; unmappable selections mutate nothing.

---

## Phase 19: Link existing + attach later from reading (US12) — gap remediation

**Purpose**: Reading toolbar can link an existing same-workspace question. Active Questions can create an unlinked question. Reading can attach a passage later.

**For a small model.** Reuse `QuestionPicker.svelte`. POST `/api/v1/questions` already creates without a note. Do not add endpoints.

**Independent Test**: Question A on note 1. From note 2 reading, link A onto a sentence — one question, two linked notes. From Active Questions create unlinked B. From a note, attach B to a passage — B is no longer unlinked.

**Do not:** link a question already on this note; link annotations; change kind; allow cross-workspace links.

### Tests for Phase 19 (write first; they must fail)

- [ ] T170 [P] [US12] Extend `web/tests/component/NoteReader.test.ts`: toolbar has `Link existing question`; it opens the picker; choosing an item calls `onLinkExisting(passage, question)`. Extend `web/tests/component/QuestionPicker.test.ts` if needed. Add `web/tests/component/ActiveQuestions.test.ts` (or extend it): `New question` without a passage posts only question fields (no note wrap).
- [ ] T171 [P] [US12] Add `tests/e2e/us12-link-attach.spec.ts`: (1) link existing question onto a second note from reading; both notes listed on the question page; (2) create unlinked from Active Questions; it appears there; attach to a passage from reading; highlight opens the same question.

### Implementation for Phase 19

- [ ] T172 [US12] Reading toolbar + picker.
  **Read:** `NoteReader.svelte`, `NoteEditor.svelte` `captureFromReader`, `QuestionPicker.svelte`, `questionsApi.list`.
  **Do:** Add toolbar button `Link existing question`. Opens picker (reuse `QuestionPicker`) with workspace questions `kind: 'question'` excluding ids already in `directiveIds(markdown)`. Add `export let onLinkExisting: ((passage: string, question: Question) => Promise<Question | void>) | undefined`. In `NoteEditor`, `linkExistingFromReader`: skip `questionsApi.create`; `findSelectionInMarkdown` / `wrapSelection` first; rebuild `questionLinks` from `directiveIds`; save; return the same question. Duplicate link: throw `That question is already in this note.` and do not save. Keep Cancel/Escape/outside-click dismissal for the picker.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/NoteEditor.test.ts tests/component/QuestionPicker.test.ts`

- [ ] T173 [US12] Unlinked create on `web/src/routes/questions/+page.svelte`.
  **Do:** Button `New question` labelled. Short form: question text required. `questionsApi.create({ workspaceId, questionText, kind: 'question', status: 'unanswered', priority: 'none', tagIds: [] })`. No note update. Reload the list. Empty text does not POST.
  **Do not:** set status answered; do not attach a fake directive.
  **Verify:** component test + existing US1/US3 e2e still pass.

- [ ] T174 [US12] Attach is the same as link-existing (T172). Ensure unlinked questions appear in the picker (`QuestionPicker` already shows `Unlinked`). After attach, Active Questions still shows one row for that id.
  **Verify:** `npx playwright test tests/e2e/us12-link-attach.spec.ts tests/e2e/us2-shared-question.spec.ts tests/e2e/us1-notes-highlight.spec.ts`

**Landmines:** Same workspace only (list is already workspace-scoped). `questionLinks` 1:1 with directives. Kind immutable. Do not wrap with an annotation id from this picker.

**Checkpoint**: Shared questions and later-attached passages work from reading view.

---

## Phase 20: Keyboard-first capture (US12) — gap remediation

**Purpose**: Selection shortcuts on the reading view; do not steal keys while typing.

**For a small model.** Edit `NoteReader` keydown only, plus document shortcuts. Do not add a keymap framework.

**Independent Test**: Select text, press Q → question composer; Escape → dismiss; select, A → annotation composer; select, L → link picker. Type the letter Q inside the composer — it inserts Q. `/notes/new` labels unchanged.

### Tests for Phase 20 (write first; they must fail)

- [ ] T175 [P] [US12] Extend `web/tests/component/NoteReader.test.ts`: with a selection and no composer, `keydown` Q opens question composer, A annotation, L picker (if `onLinkExisting` provided), Escape closes toolbar. With composer open or target `input`/`textarea`, Q does not toggle kind. Add `tests/e2e/us12-keyboard-capture.spec.ts` for Q then Escape on a selected sentence.

### Implementation for Phase 20

- [ ] T176 [US12] Extend `onWindowKeydown` in `web/src/lib/editor/NoteReader.svelte`.
  **Do:** If `event.defaultPrevented`, return. If target is `input, textarea, select, [contenteditable]`, return (Escape may still close composer — keep current Escape behavior when composer is open). If composer open, ignore Q/A/L. If no `selectedPassage` and no `showToolbar`, ignore Q/A/L. Otherwise: `q`/`Q` → `openComposer('question')`; `a`/`A` → `openComposer('annotation')`; `l`/`L` → open link picker if `onLinkExisting` exists. `preventDefault` only when handling those keys. Do not handle shortcuts with Ctrl/Meta/Alt.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts` then `npx playwright test tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts`

- [ ] T177 [P] [US12] Add a one-line hint on the reading toolbar: `Q ask · A annotate · L link · Esc cancel`. No new route.
  **Verify:** NoteReader component test asserts the hint exists when the toolbar is visible.

**Landmines:** Do not steal keys in Title/Note fields. Do not change `/notes/new`. Keep existing Escape/outside-click dismissal.

**Checkpoint**: Capture is keyboard-complete without breaking typing or accessibility labels.

---

## Phase 21: Map, spec trace, and regression (after US8–US12)

**Purpose**: Keep the mental model honest after the new slices.

- [ ] T178 [P] Refresh [`PROJECT_MAP.md`](../../PROJECT_MAP.md): remaining-gaps list must match reality; add file pointers for `QuestionContext.svelte`, `NoteExcerpt.svelte`, `NoteQuestionRail.svelte`, `nextQueue.ts`, `/next`, `findSelectionInMarkdown`, and keyboard shortcuts. Known-bugs table: only real failures.
- [ ] T179 Trace FR-041 through FR-055 and SC-011 through SC-014 to tests in `specs/001-track-note-questions/traceability.md`. Do not invent passing evidence.
- [ ] T180 Run the new Playwright files plus `tests/e2e/us1-notes-highlight.spec.ts`, `tests/e2e/us3-question-lifecycle.spec.ts`, `tests/e2e/us2-shared-question.spec.ts`, `tests/e2e/accessibility-notes.spec.ts`. Record results in `specs/001-track-note-questions/validation-results.md` only for runs you actually executed.

**Checkpoint**: A new session can implement leftover tasks from the map without rediscovering the tree.

---

## Dependencies & Execution Order

### Phase dependencies

1. **Phase 1 — Setup** has no dependencies.
2. **Phase 2 — Foundational** depends on Phase 1 and blocks all stories.
3. **Phase 3 — US1** depends on Phase 2 and establishes first-run workspace bootstrap/selection before basic note/question capture.
4. **Phase 4 — US2** depends on US1’s canonical note/question primitives.
5. **Phase 5 — US3** depends on US1; it can run in parallel with US2 after US1, though final P1 validation uses both.
6. **Phase 6 — US4** depends on US3 status behavior and the foundational schema. It does not require US2.
7. **Phase 7 — US5** depends on US1 workspace/note primitives. It can run in parallel with US3/US4, but US4’s topic/tag filter acceptance uses US5-created organization data for final integrated validation.
8. **Phase 8 — US6** depends on US2 shared links and US5 note hierarchy.
9. **Phase 9 — US7** depends on all entity-producing stories selected for export; implement after US1–US6 for the complete MVP archive.
10. **Phase 10 — Polish** depends on all stories intended for the release.
11. **Phase 11 — Notes list + highlight/annotation** depends on US1 capture primitives (notes/questions API already exist). Implement after the current app runs. Start the session from `PROJECT_MAP.md`. T132 (notes list) can ship before T131/T133–T135.
12. **Phase 12 — Note editing, capture dismissal, filters, and tag discovery** depends on Phase 11 plus the existing note-tag persistence and search primitives; implement T138 as a coordinated frontend/API slice.
13. **Phase 13 — Answer in context** depends on Phase 11 reading view and US3 lifecycle. No API changes.
14. **Phase 14 — Note-local rail** depends on Phase 13 `focusQuestionId` / card opening, but may start after T148’s NoteReader prop exists.
15. **Phase 15 — Next queue** depends on Phase 13 list links (T144) so rows open in context. Uses existing list API only.
16. **Phase 16 — Deferred resume date** depends on US3 lifecycle. Implement before relying on Next’s Deferred ready section in manual demos; Next grouping (T152) already encodes the rule.
17. **Phase 17 — Resolved highlights / insert answer** depends on Phase 11 wraps and US3 answered status. Insert on the context page needs Phase 13.
18. **Phase 18 — Passage matching** depends on Phase 11 `wrapSelection`. Do this before Phase 19 if possible so link-existing uses the same matcher.
19. **Phase 19 — Link/attach from reading** depends on Phase 18 matcher and US2 picker. Unlinked create uses existing POST `/questions`.
20. **Phase 20 — Keyboard capture** depends on Phase 19 toolbar/picker so L has a target; Q/A/Escape can ship after Phase 11 if L is ignored when `onLinkExisting` is missing.
21. **Phase 21 — Map/trace** depends on the phases you actually shipped.

### User story dependency graph

```text
Setup → Foundation → US1 Workspace bootstrap → US1 Capture
                                             ├──→ US2 Shared Questions ───────┐
                                             ├──→ US3 Lifecycle → US4 Find   ├──→ US6 Safe Delete
                                             └──→ US5 Organization ───────────┘

US1 + US2 + US3 + US4 + US5 + US6 → US7 Import/Export → Polish
                                                      ↘ Phase 11 Notes list + highlights (uses US1 APIs)
                                                         → Phase 13 context → Phase 14 rail
                                                         → Phase 15 Next (after T144)
                                                         → Phase 16 defer date (US3)
                                                         → Phase 17 resolved / insert
                                                         → Phase 18 wrap match → Phase 19 link/attach → Phase 20 keys
```

### Within each user story

1. Add and run story tests first; confirm failures describe missing behavior rather than broken setup.
2. Add domain types/rules before repository logic.
3. Add repository logic before services and handlers.
4. Add handlers before wiring frontend API consumers. For US1, workspace creation and list hydration must be available before note creation is enabled.
5. Add reusable components before route integration; the no-workspace guard must precede the New note route.
6. Run the story’s independent Playwright test and all prior-story regression tests at its checkpoint.

---

## Parallel Execution Examples

### User Story 1

```text
Parallel tests: T024 domain | T025 SQLite | T026 API contract | T027 components | T028 browser
Additional first-run regression: T028a fresh-database workspace bootstrap (must pass before the US1 checkpoint).
Parallel types after failing tests: T029 Go domain | T030 frontend types
Parallel UI after API availability: T034 Markdown engine | T036 workspace screen
Sequential workspace prerequisite: T036a bootstrap/guard → T037 note creation/editing.
```

### User Story 2

```text
Parallel tests: T039 service | T040 SQLite | T041 components | T042 browser
Parallel UI after shared-question API: T046 question picker | T047 linked-note context
```

### User Story 3

```text
Parallel tests: T049 lifecycle domain | T050 persistence/API | T051 components | T052 browser
Parallel UI after lifecycle API: T056 lifecycle component | T058 answered route
```

### User Story 4

```text
Parallel tests: T060 domain | T061 search/filter store | T062 API contract | T063 components | T064 browser
Parallel implementation: T065 domain types | T070 schedule UI | T071 filter UI | T073 search UI after endpoint shape is stable
```

### User Story 5

```text
Parallel tests: T075 domain | T076 store | T077 API contract | T078 components | T079 browser
Parallel UI after organization API: T084 workspace switcher | T085 note tree | T086 selectors | T087 sharing review
```

### User Story 6

```text
Parallel tests: T089 service | T090 store | T091 API contract | T092 component | T093 browser
Sequential implementation for safety: T094 preview → T095 atomic execution → T096 handlers → T097 UI
```

### User Story 7

```text
Parallel tests: T098 archive | T099 database | T100 API contract | T101 components | T102 browser
Parallel implementation after types: T104 export and T108 browser transfer client; T109 UI after preview response is stable
```

### Phase 11 (Notes list + highlight/annotation)

```text
Start: read PROJECT_MAP.md (repo root), then only the files it names.
Parallel tests: T125 directives | T126 reader/card | T127 notes index | T128 kind domain/store | T129 browser
Parallel after failing tests: T130 contracts | T131 kind backend | T132 notes list (can merge first)
Sequential UX: T133 wrap/render → T134 reader/card/editor → T135 create+filter → T136 refresh PROJECT_MAP.md
```

### Phases 13–20 (small-model slices)

```text
One agent = one task. Never start T14x while T13x tests are unwritten.
13: T139–T140 tests → T141 NoteExcerpt → T142 QuestionContext → T143 page → T144 list links
14: T145–T146 tests → T147 rail → T148 NoteReader focus → T149 editor wire
15: T150–T151 tests → T152 helper → T153 /next page → T154 nav
16: T155–T156 tests → T157 domain → T158 API mapping → T159 lifecycle UI
17: T160–T162 tests → T163 insert helper → T164 render/CSS → T165 card + save
18: T166–T167 tests → T168 findSelection → T169 wrap-before-create
19: T170–T171 tests → T172 link existing → T173 unlinked create → T174 attach = link
20: T175 tests → T176 keydown → T177 toolbar hint
21: T178 map → T179 trace → T180 only tests you ran
```

---

## Implementation Strategy

### Smallest demonstrable MVP

1. Complete Setup and Foundational phases.
2. Complete US1 through T038, including T028a and T036a for first-run workspace bootstrap.
3. Run the US1 checkpoint offline from an empty database.
4. Demo creating a workspace first, then a note with two embedded questions, the active Questions view, and source-note navigation.

This is the smallest useful slice, but it is not the complete product MVP described by the specification. The complete specification MVP includes all seven stories.

### P1 release increment

1. Complete US1 capture, including first-run workspace creation/selection and the workspace-required note guard.
2. Add US2 shared questions and validate canonical edits.
3. Add US3 lifecycle and validate answer/reopen behavior.
4. Run all P1 browser and contract tests before starting P2 features.

### Incremental complete MVP

1. **Foundation** → reproducible local shell, API, SQLite, contracts, tests, and a usable zero-workspace state.
2. **US1–US3** → create/select a workspace first, then complete the core question workflow.
3. **US4–US5** → discovery, prioritization, and organization.
4. **US6** → safe destructive operations.
5. **US7** → user-controlled portability and restoration.
6. **Polish** → measured success criteria and release evidence.
7. **Phase 13–14** → answer beside the passage; next gap on the note.
8. **Phase 15–16** → Next queue and deferred resume dates.
9. **Phase 17–20** → resolved highlights, sturdier wrap, link/attach, keyboard.

### Team parallelism

After Foundation and US1:

- Stream A: US2, then US6 after US5 is available.
- Stream B: US3, then US4.
- Stream C: US5.
- Stream D: prepare US7 archive tests and schemas, then implement after entity shapes stabilize.

All streams must update the same OpenAPI contract and shared fixtures in coordinated changes; avoid parallel edits to `server/internal/api/router.go`, `web/src/routes/+layout.svelte`, or `web/src/lib/stores/query.ts` without sequencing.

---

## Notes

- `[P]` means the task is safe to assign concurrently under the dependency notes; it does not remove phase prerequisites.
- Story labels provide traceability to acceptance scenarios in `spec.md`.
- OpenAPI is the authoritative transport contract; backend behavior is authoritative for business rules.
- Do not cache/replay API writes in the service worker or introduce a second browser database.
- Never log note content, question text, or answer content.
- Commit after each task or small logical group and run affected tests at every checkpoint.
- Phase 11 is gap remediation on top of a working US1–US7 tree. A new implementation session must read `PROJECT_MAP.md` first and follow its file map rather than rereading the repository.
- Phases 13–20 are further gap remediation (US8–US12). A new session reads `PROJECT_MAP.md`, **only the current phase in this file**, and only named files. One task per session for small models. Do not start a later phase to “also add” extra product ideas (AI, flashcards, saved filters, Feynman modes).
