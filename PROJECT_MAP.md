# PROJECT_MAP — Noted (LLM mental model)

Read this file first in a new session. Then open only the files listed for the task.
Do not scan the whole repo unless the map is stale.

**Product**: Local-first note app. User reads a note, captures questions/annotations in context, and comes back later via Active Questions. One local user. No auth. Go API is the source of truth. SvelteKit is UI only.

**Goal (current product direction)**: Make reading interactive.
1. Notes must appear on the Notes page.
2. Do not show raw `{{question:<uuid>}}` tokens to the user.
3. Highlight a passage → right-click → add a **question** or **annotation**.
4. Click the highlight → open a **card** (question/annotation).
5. Unanswered **questions** still appear under Active Questions.
6. Annotations are comments on a passage; they should not clutter Active Questions.

**Shipped follow-on behavior (Phases 13–20):** questions open beside their source passage; notes have an ordered open-question rail and Next unanswered navigation; `/next` groups the daily queue; deferred questions require a resume date; answered highlights are resolved and can optionally insert their answer; rendered selections use Markdown-aware matching; reading view can link/attach existing questions; and Q/A/L/Escape keyboard capture is available without stealing typing keys.

**Remaining planned work:** Phase 23 still needs canonical question-text editing in the lifecycle surfaces. Phase 24 has an approved optional `consequenceText` design, but its contract, persistence, search, portability, and UI work is not implemented. Phase 22's delete flow is already wired in the current note page; its dedicated regression/evidence tasks remain in the task list.

---

## Run locally

Two processes:

| Process | URL | How |
|---|---|---|
| Go API + SQLite | `http://127.0.0.1:8080` | `go -C server run ./cmd/noted` with `NOTED_DATA_DIR=.local/data` |
| Vite / SvelteKit | `http://127.0.0.1:5555` | `npm --prefix web run dev -- --host 127.0.0.1` |

Vite proxies `/api`, `/healthz`, `/readyz` to `:8080` (`web/vite.config.ts`).
DB file: `.local/data/notes.db`.

Useful commands (repo root):

```
make test-server    # Go tests
make test-web       # Vitest
make test-e2e       # Playwright
```

Frontend tests live in `web/tests/`. E2E in `tests/e2e/`.

---

## Architecture (do not break)

```
Browser (SvelteKit static app, web/)
    │  fetch /api/v1/*
    ▼
Go HTTP API (server/)  ← authoritative validation, lifecycle, persistence
    ▼
SQLite + FTS5  (.local/data/notes.db)
```

- Frontend validation is advisory only.
- Do not cache/replay API writes in the service worker.
- IDs: lowercase UUIDs. Optimistic `version` on updates. Errors: `{ error: { code, message, ... } }`.
- Notes body is Markdown. Question placement is a directive in `bodyMarkdown` plus a `note_questions` row. They must match 1:1, same order.
- OpenAPI is the transport contract: `specs/001-track-note-questions/contracts/openapi.yaml`
  Copy to `tests/fixtures/contracts/` when it changes.
- JSON decoder **rejects unknown fields** (`server/internal/api/jsoncodec/json.go`). New request fields must exist on the Go struct **and** in OpenAPI.

---

## Known bugs / in-progress work (current)

The Phase 11 and Phase 13–20 behavior is shipped; do not treat those completed gaps as current bugs.

| Symptom / gap | Current reality | Follow-up |
|---|---|---|
| Canonical question text is not editable from the user-facing question surfaces | The API can update `questionText`, but `QuestionLifecycle.svelte` still edits answer/status/schedule only, and highlight-card/context propagation for text edits is not implemented. | Phase 23, T186–T189b. |
| Optional “Why it matters if unanswered” context is absent | The approved `consequenceText` design is documented in Phase 24, but no contract, `Question` field, database column, FTS update, import allow-list, or UI exists yet. | Phase 24, T191–T195a. |
| Phase 22 has no dedicated regression evidence | The existing saved-note page already calls `notesApi.previewDeletion`, opens `DeleteNoteReview`, calls `deleteReviewed`, and navigates to `/notes` after success. The remaining work is the explicitly listed component/E2E coverage and any cache-invalidation proof. | Phase 22, T181–T185. |

The only currently recorded product-test failure is the strict locator in `tests/e2e/us1-notes-highlight.spec.ts`: the selected sentence appears both in the reader and in the highlight card. Scope that assertion to the intended surface rather than changing the reader. Firefox/WebKit executable availability is an environment prerequisite, not an application bug; see `issues.md`.

`NoteTree.svelte` renders the Notes index hierarchy; keep it aligned with the notes list API.

---

## Intended capture model (implement toward this)

**Reading is the default** for an existing note. Edit mode is a separate toggle (new notes start in edit so `/notes/new` keeps labelled Title + Note controls — Playwright `tests/e2e/accessibility-notes.spec.ts` requires `getByLabel('Title')` and `getByLabel('Note')`).

1. User selects text in the **reading view**.
2. Right-click (and a selection toolbar, because mobile has no right-click) → **Ask a question** or **Add annotation**.
3. Composer dialog: passage quote + text field.
4. Create question/annotation via `POST /api/v1/questions`.
5. Wrap the selected passage in Markdown (do not leave a visible token):

   ```
   {{question:<uuid>}}selected passage{{/question}}
   ```

   Bare `{{question:<uuid>}}` remains valid (legacy / unanchored). Parser already finds the opening token. Closing `{{/question}}` is not a second directive.

6. Save the note with `questionLinks` rebuilt from directive order (`directiveIds(body)` must equal `questionLinks[].questionId` in order — backend rejects mismatch).
7. Click highlight → card (dialog). Questions: answer/status via `QuestionLifecycle.svelte`. Annotations: edit the comment only.
8. Active Questions lists **kind=question** only, statuses unanswered / in_progress / deferred.

**Kind**:

- `question` (default) → Active and Answered Questions
- `annotation` → highlight + card only
- Implemented by migration `002_question_kind.sql`, the Go domain/store/API contract, import allow-list, and frontend question filters. Kind is immutable after creation.

**Anchor limitation**: `findSelectionInMarkdown()` in `web/src/lib/editor/directives.ts` first tries an exact match, then whitespace collapse, then skips Markdown emphasis/code markers while mapping back to the source span. It still fails safely when the rendered selection cannot be mapped, and overlapping highlights remain out of scope.

Legacy display modes (`expanded` / `collapsed` / `link`) can stay on `note_questions` but the reader should ignore them and always use highlight + card.

---

## Directory map — what to edit

### Specs (update when product behavior changes)

| File | Role |
|---|---|
| `specs/001-track-note-questions/spec.md` | User stories, FRs, scope, and the Phase 11 highlight/annotation behavior. |
| `specs/001-track-note-questions/plan.md` | Architecture decisions, including wrapped directives and the reading view. |
| `specs/001-track-note-questions/data-model.md` | Tables, constraints, directive/link integrity. |
| `specs/001-track-note-questions/contracts/note-content.md` | Markdown grammar. Copy to `tests/fixtures/contracts/note-content.md`. |
| `specs/001-track-note-questions/contracts/openapi.yaml` | API. Copy to `tests/fixtures/contracts/openapi.yaml`. |
| `specs/001-track-note-questions/tasks.md` | Task list. Add gap-remediation tasks rather than renumbering. |
| `specs/001-track-note-questions/research.md` | Decision log (note format §7). |

### Frontend — notes / questions UX (most likely)

| File | Role |
|---|---|
| `web/src/routes/notes/+page.svelte` | Notes index; loads the selected workspace's notes and renders `NoteTree`. |
| `web/src/routes/notes/new/+page.svelte` | Create note. Must not POST until a workspace is selected. |
| `web/src/routes/notes/[noteId]/+page.svelte` | Open/edit/delete one note. Hosts `NoteEditor`. |
| `web/src/lib/editor/NoteEditor.svelte` | Title/textarea edit mode for new notes; reading mode, rail navigation, selection capture, link/attach, and answer insertion for existing notes. |
| `web/src/lib/editor/NoteReader.svelte` | Rendered reading surface, selection toolbar/composer, link picker, highlight cards, and Q/A/L/Escape keyboard shortcuts. |
| `web/src/lib/editor/directives.ts` | Directive tokenization, `insertDirective`, `directiveIds`, Markdown-aware `findSelectionInMarkdown`, wrapping, and idempotent answer insertion. |
| `web/src/lib/editor/markdown.ts` | `marked` + DOMPurify; `renderNoteHtml()` turns wraps into `<mark data-annotation-id>`. |
| `web/src/lib/components/InlineQuestion.svelte` | Unused in the body today. Card UI can replace or wrap this. |
| `web/src/lib/components/NoteTree.svelte` | Hierarchy navigation used by the Notes index. |
| `web/src/lib/components/QuestionLifecycle.svelte` | Answer, status, due-date, and priority save. Reuse inside the highlight card; canonical question-text editing remains Phase 23 work. |
| `web/src/lib/components/QuestionContext.svelte` | Question detail surface: source-note excerpt beside lifecycle/schedule controls, note switching, and optional answer insertion. |
| `web/src/lib/components/NoteExcerpt.svelte` | Sanitized source-note rendering for a question, with status-aware highlight styling and scroll-to-highlight. |
| `web/src/lib/components/NoteQuestionRail.svelte` | Ordered open-question rail for a note; excludes annotations and answered questions. |
| `web/src/lib/components/QuestionList.svelte` | Active/Answered list cards. |
| `web/src/lib/components/QuestionPicker.svelte` | Link an existing question (edit-mode and reading view). |
| `web/src/lib/components/NextQueue.svelte` | Presentational sections and links for the `/next` queue. |
| `web/src/lib/api/notes.ts` | `list/create/get/update/previewDeletion/deleteReviewed` |
| `web/src/lib/api/questions.ts` | `list/get/create/update`, including `kind` query/body support |
| `web/src/lib/types/note.ts` | `Note`, `NoteWrite`, `NoteQuestionLink` |
| `web/src/lib/types/question.ts` | Status, priority, displayMode, and `QuestionKind` types. |
| `web/src/routes/questions/+page.svelte` | Active Questions. Filter `kind=question`. |
| `web/src/routes/answered/+page.svelte` | Answered view. Same kind filter. |
| `web/src/routes/questions/[questionId]/+page.svelte` | Full question page; loads the selected linked note and renders `QuestionContext`. |
| `web/src/routes/next/+page.svelte` | Workspace-scoped daily Next queue; loads up to five question pages and handles retry/partial results. |
| `web/src/lib/questions/nextQueue.ts` | Pure client-side grouping into Overdue, Due today, In progress, Deferred ready, and High priority. |
| `web/src/lib/stores/workspace.ts` | Current workspace; hydrate before note mutations. |
| `web/src/lib/components/AppNavigation.svelte` | Notes / Active Questions / Next / Answered / Search / Data. |
| `web/src/routes/+layout.svelte` | Shell, workspace hydration, toasts |
| `web/src/routes/+page.svelte` | Home / first-run workspace CTA |

Highlight support files already present:

- `web/src/lib/components/AnnotationCard.svelte` — dialog/card opened from a highlight; annotation editing and optional answer insertion.
- `web/src/lib/components/QuestionContext.svelte` — answer-in-context detail surface.
- `web/src/lib/components/NoteExcerpt.svelte` — source-note excerpt used by the context surface.
- `web/src/lib/components/NoteQuestionRail.svelte` — note-local open-question navigation.
- `web/src/lib/questions/nextQueue.ts` and `web/src/routes/next/+page.svelte` — the client-side Next queue and route.

Tests to keep green when changing the editor:

- `web/tests/component/NoteEditor.test.ts` — Title, Note, Save note
- `web/tests/unit/markdown-security.test.ts` — raw HTML rejected; directive IDs parsed
- `tests/e2e/accessibility-notes.spec.ts` — `/notes/new` has Title + Note labels
- `tests/e2e/us1-workspace-bootstrap.spec.ts` — no note request before workspace

### Backend — only if kind / directive grammar / list filters change

| File | Role |
|---|---|
| `server/internal/domain/note.go` | `ParseQuestionDirectives`, note validation. Opening `{{question:uuid}}` already matches wrapped form. |
| `server/internal/domain/question.go` | Question entity, status, displayMode, Validate() |
| `server/internal/domain/prioritization.go` | `QuestionQuery` + `NormalizeQuestionQuery` |
| `server/internal/store/migrations/` | Embedded schema migrations, including `002_question_kind.sql`. |
| `server/internal/store/sqlite/capture.go` | Workspace/note/question CRUD. `CreateNote`/`UpdateNote` require directives ≡ links. Question INSERT/SELECT/UPDATE column lists. `noteLinks` join for summaries. |
| `server/internal/store/sqlite/question_query.go` | Filtered question list (Active Questions API) |
| `server/internal/store/sqlite/lifecycle.go` | Status/answer updates — must preserve new fields (e.g. kind) |
| `server/internal/api/handlers/capture.go` | `GET/POST /notes`, `POST/GET/PUT /questions` |
| `server/internal/api/handlers/questions_query.go` | Parses list query params |
| `server/internal/importexport/apply.go` | **Allow-list of JSON fields per table.** New columns must be added or import drops them. Export is `SELECT *` so it follows the schema. |
| `server/cmd/noted/main.go` | Process wiring |

Do not hand-edit `server/webdist/` (built frontend).

---

## Data model (minimum)

```
Workspace 1—* Note
Workspace 1—* Question
Note *—* Question  via note_questions (display_mode, position)
```

**Note**: `title`, `body_markdown`, `parent_note_id`, `topic_id`, `version`
**Question**: `question_text`, `kind` (`question|annotation`), `answer_markdown`, `status` (`unanswered|in_progress|deferred|answered`), `priority`, `due_date`, `version`
**note_questions**: PK `(note_id, question_id)`, unique `(note_id, position)`, `display_mode` (`expanded|collapsed|link`)

Integrity: every `{{question:id}}` in the body ↔ exactly one link, same order, same workspace, at most once per note.
Answered status requires non-empty answer.
Note delete is two-step: preview token, then explicit decisions for singly-linked questions.

---

## API cheat sheet

| Method | Path | Use |
|---|---|---|
| GET | `/api/v1/notes?workspaceId=` | Notes list used by the Notes index |
| POST/PUT | `/api/v1/notes` `/api/v1/notes/{id}` | Body must include `questionLinks` matching directives |
| GET | `/api/v1/questions?workspaceId=&status=&kind=` | Active/answered lists |
| POST | `/api/v1/questions` | Create before wrapping the passage (need the id) |
| PUT | `/api/v1/questions/{id}` | Includes `version`. Answer/status. |
| POST | `/api/v1/notes/{id}/deletion-preview` then `/delete` | Safe delete |

Frontend clients: `web/src/lib/api/{notes,questions,workspaces,tags,search,dataPortability}.ts`

---

## Common tasks → files

**Notes missing from Notes page**
→ Resolved in `web/src/routes/notes/+page.svelte`; the backend list and `NoteTree.svelte` navigation are wired.

**Hide `{{question:uuid}}` / highlight + card**
→ Resolved by `directives.ts`, `markdown.ts`, `NoteEditor.svelte`, `NoteReader.svelte`, `AnnotationCard.svelte`, and the note `[noteId]` page. Keep edit textarea for `/notes/new` labels.

**Questions vs annotations in Active Questions**
→ `kind` on questions: migration, `domain/question.go`, `capture.go` SQL, `question_query.go`, `questions_query.go`, OpenAPI + fixtures, `importexport/apply.go` allow-list, `web/src/lib/types/question.ts`, `questions.ts`, questions/answered pages.

**Change Markdown directive grammar**
→ `domain/note.go` + `note_test.go`, `directives.ts`, `contracts/note-content.md` (both copies). Wrapped form does **not** require a parser change for ID extraction.

**Answer in context / note-local navigation**
→ `web/src/routes/questions/[questionId]/+page.svelte`, `QuestionContext.svelte`, `NoteExcerpt.svelte`, `NoteQuestionRail.svelte`, `NoteReader.svelte`, and `NoteEditor.svelte`. Keep the rail question-only and ordered by `directiveIds(bodyMarkdown)`.

**Next queue**
→ `web/src/lib/questions/nextQueue.ts` for pure grouping, `web/src/lib/components/NextQueue.svelte` for presentation, and `web/src/routes/next/+page.svelte` / `AppNavigation.svelte` for workspace loading and navigation. Do not add a queue API or saved filters.

**Rendered selection matching / keyboard capture**
→ `web/src/lib/editor/directives.ts` (`findSelectionInMarkdown`, `wrapSelection`) and `web/src/lib/editor/NoteReader.svelte` (`onWindowKeydown`). Keep matching mutation-free on failure and do not capture Q/A/L while typing.

**New API field**
→ OpenAPI + fixtures, Go struct, store SQL, import allow-list if persisted, frontend type + client. Unknown JSON fields 400. The approved but unimplemented Phase 24 `consequenceText` field must follow this path before any request includes it.

---

## Constraints / landmines

- Svelte 5, but most components still use `export let` (not runes). Match existing style.
- `/notes/new` must keep accessible Title + Note labels.
- Workspace must exist before any note create (`web/src/stores/workspace.ts`, new-note guard).
- `questionLinks[i].position` must be `0..n-1` contiguous and match directive order.
- Do not log note/question/answer bodies.
- Highlights were originally out of MVP in `spec.md` scope — update the spec if you ship them.
- `tests/e2e/us1-capture-track.spec.ts` is currently a shallow home-page smoke test, not a full capture journey.

---

## Remaining work and product gaps (do not pretend these exist)

Shipped: Notes index, highlight/annotation capture, kind filtering, existing-note edit, capture dismissal, tag search, and all Phase 13–20 outcomes (answer-in-context, note rail, `/next`, deferred dates, resolved highlights/answer insertion, Markdown-aware matching, link/attach, and keyboard capture).

| Phase | Current gap or remaining evidence | Start files |
|---|---|---|
| 22 | Delete-note behavior is already wired through preview/review/confirmed delete; dedicated component/E2E regression coverage and cache-invalidation evidence remain. | `web/src/routes/notes/[noteId]/+page.svelte`, `web/src/lib/components/DeleteNoteReview.svelte`, `web/src/lib/api/notes.ts` |
| 23 | Canonical question-text editing is not exposed in `QuestionLifecycle` or propagated through highlight cards. | `web/src/lib/components/QuestionLifecycle.svelte`, `AnnotationCard.svelte`, `NoteReader.svelte`, `QuestionContext.svelte` |
| 24 | Optional question-only `consequenceText` is approved but not implemented in contracts, model, schema, FTS, portability, or UI. | `specs/001-track-note-questions/tasks.md` Phase 24; follow T192–T194 in order |

Still out of scope:

- No dedicated Annotations inbox
- No overlapping highlights
- Edit mode still shows raw Markdown if you toggle it
- Images / PDFs / stable offset anchors (Phase 18 is match-on-wrap only, not a new schema)
- Cross-workspace moves not in MVP
- Saved filters, AI, flashcards, Feynman modes
- Display-mode controls (expanded/collapsed/link) are legacy once cards shipped
