# Tasks: Interactive Note Questions

**Input**: Design documents from `/specs/001-track-note-questions/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Completed work**: Phases 1–12 (T001–T138), Phase 13 T139–T142, and Phase 24 T190 are archived in [`completed-task.md`](./completed-task.md). Do not re-implement them.

**New-session start**: Read [`PROJECT_MAP.md`](../../PROJECT_MAP.md) at the repository root before opening other files. Then read **only the current phase in this file**, then only the files that phase names. For leftover Phase 13 wiring use **T143, T143a, T144**. For highlight/annotation follow-on work, Next queue, deferred dates, resolved highlights, sturdier wrapping, reading-view link/attach, and keyboard capture, use **Phases 14–20a**. For the emergency capture-textbox focus regression, use **Phase 21a**; for the note-rail card-opening dismissal regression, use **Phase 21b**. For the no-op Delete note button, use **Phase 22**; for canonical question-text editing, use **Phase 23**; for the approved optional consequence/why-it-matters behavior, use **Phase 24**. Do not scan the whole repo unless the map is stale.

**Tests**: Automated tests are included because the project constitution requires unit, integration, API contract, browser workflow, and container smoke coverage. Within each story, create the listed tests first and verify they fail for the expected missing behavior before implementation.

**Organization**: Remaining tasks are grouped by gap-remediation phase. Every phase ends with an independently executable browser scenario and checkpoint.

**Workspace prerequisite**: A fresh database can contain zero workspaces. Workspace bootstrap already shipped in Phases 1–3 (see `completed-task.md`). Remaining work must not regress that guard.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can be implemented in parallel with adjacent tasks because it changes different files and has no dependency on an unfinished adjacent task.
- **[Story]**: Maps the task to a user story in `spec.md`.
- Paths are relative to the repository root.
- `a`-suffixed IDs are gap-remediation tasks inserted without renumbering later tasks. Letter siblings (`T143a`, `T143b`) are sequential slices of one original task. A small model implements **one ID per session**, in order: parent, then `a`, then `b`.

---

## Phase 13: Answer in context (US8) — gap remediation

**Purpose**: Opening a question shows the source note and highlighted passage on the same page as answer/status controls.

**For a small model.** Do not scan the repo. Do not reread US1–US7 or Phase 11 internals. Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement **one task**, run its verify command, stop.

**Independent Test**: Capture a question on a sentence. From Active Questions open it. The question page shows the note, the sentence is highlighted and in view, and the answer box is on the same page. An unlinked question still opens and says it is unlinked.

**Do not**: change capture, kind, lifecycle rules, NoteReader toolbar, `/notes/new` labels, or add APIs.

**Already done (do not rewrite):** T139–T142 — `NoteExcerpt.test.ts`, `QuestionContext.test.ts`, `us8-answer-in-context.spec.ts`, `NoteExcerpt.svelte`, `QuestionContext.svelte`. Details in [`completed-task.md`](./completed-task.md).

### Implementation remaining for Phase 13

Current page already loads `questionsApi.get` in `onMount`, keeps `saved(value)` to replace `question`, and renders `QuestionLifecycle` + `QuestionSchedule` + a source-notes list. `QuestionContext.svelte` already exists with `export let question`, `note`, `selectedNoteId`, `onSelectNote`, `onSave`. Match `export let`. Do not use `NoteReader` here.

- [x] T143 [US8] Load the source note and render `QuestionContext` on `web/src/routes/questions/[questionId]/+page.svelte`.
  **Read:** that page, `web/src/lib/components/QuestionContext.svelte`, `web/src/lib/api/questions.ts`, `web/src/lib/api/notes.ts`, `$app/state` `page`.
  **Do:** Keep `questionsApi.get`. After the question loads, `noteId` = `page.url.searchParams.get('noteId')` if that id is in `question.linkedNotes`, else `question.linkedNotes[0]?.id`. If `noteId`, load `notesApi.get` into `note`; else `note = null`. Replace the lifecycle/schedule/source-notes block with `<QuestionContext {question} {note} selectedNoteId={noteId ?? ''} onSelectNote={...} onSave={saved} />`. `onSelectNote` may only set `selectedNoteId` in this task (reload is T143a). Keep the existing loading and error UI.
  **Do not:** add API fields; do not log question/note bodies; do not change `QuestionContext.svelte` props.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionContext.test.ts tests/component/NoteExcerpt.test.ts`

- [x] T143a [US8] Reload the note when the context switcher changes.
  **Read:** `web/src/routes/questions/[questionId]/+page.svelte` only.
  **Do:** `onSelectNote(id)` sets `selectedNoteId`, loads `notesApi.get(id)`, and assigns `note`. Keep the current question. If the page has no retry control, add the same Retry pattern as `web/src/routes/questions/+page.svelte` (error + button that re-runs the load). Unlinked questions stay on `note = null` and the existing “currently unlinked” copy inside `QuestionContext`.
  **Do not:** save notes; do not add `NoteReader`.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionContext.test.ts`

- [x] T144 [US8] Point list cards at context.
  **Read:** `web/src/lib/components/QuestionList.svelte` only. The title is currently a link to `/questions/{id}` with no query string.
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

- [x] T145 [P] [US8] Add `web/tests/component/NoteQuestionRail.test.ts`: given mixed questions, the rail lists only `kind=question` with status `unanswered` / `in_progress` / `deferred`, in the provided `orderedIds` order; clicking an item calls `onSelect(id)`; Next unanswered calls `onNext`; when `remaining === 0` the next control is disabled or the rail text is `No open questions in this note`.
- [x] T146 [P] [US8] Extend `web/tests/component/NoteReader.test.ts` only: when `focusQuestionId` is set to a wrapped id, the card opens (`openQuestion` / dialog). Do not write Playwright in this task.
- [x] T146a [P] [US8] Add `tests/e2e/us8-note-rail.spec.ts` (do not extend `us8-answer-in-context.spec.ts`): two highlight questions on one note, Next unanswered twice, then the empty message. Annotations and answered items must not appear in the rail.

### Implementation for Phase 14

Read-mode markup today is `{#if mode === 'read' && existing}<NoteReader markdown={bodyMarkdown} {questions} onCapture={captureFromReader} />`. `directiveIds` is already imported in `NoteEditor.svelte`. New unsaved notes stay in edit mode — do not show the rail on `/notes/new`.

- [x] T147 [P] [US8] Create `web/src/lib/components/NoteQuestionRail.svelte`.
  **Do:** `export let questions: Question[] = []`; `export let orderedIds: string[] = []`; `export let onSelect: ((id: string) => void) | undefined`; `export let onNext: (() => void) | undefined`. Filter `kind !== 'annotation'` (missing kind counts as question) and `status !== 'answered'`. Sort by `orderedIds`. Render `<aside aria-label="Open questions in this note">` with a button per item (question text + status) and a “Next unanswered” button. Empty: `No open questions in this note`.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteQuestionRail.test.ts`

- [x] T148 [US8] Open a highlight from the outside.
  **Read:** `web/src/lib/editor/NoteReader.svelte` (`onClick` already finds `mark[data-annotation-id]`, then `questions.find`, then `tokenizeDirectives` for `openPassage`).
  **Do:** Add `export let focusQuestionId: string | null = null`. When `focusQuestionId` changes to a non-null id present in `questions`, set `openQuestion` / `openPassage` the same way that click path does. Do not clear `focusQuestionId` yourself if it is a prop; parent may reset it.
  **Do not:** remove Cancel / Escape / outside-click dismissal; do not change capture/toolbar.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T149 [US8] Mount the rail beside the reader.
  **Read:** `web/src/lib/editor/NoteEditor.svelte` only around the read-mode `NoteReader` block.
  **Do:** Beside `NoteReader` (not inside the article), render `NoteQuestionRail` with `questions` and `orderedIds={directiveIds(bodyMarkdown)}`. Leave `onSelect` / `onNext` unset. Do not add `focusId` yet. Do not show the rail in edit mode.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts tests/component/NoteQuestionRail.test.ts`

- [x] T149a [US8] Walk open questions from the rail.
  **Read:** `web/src/lib/editor/NoteEditor.svelte` read-mode block only.
  **Do:** `let focusId: string | null = null`. Rail `onSelect(id)` sets `focusId = id`. `onNext` sets `focusId` to the next open id in `directiveIds(bodyMarkdown)` order (skip `kind === 'annotation'` and `status === 'answered'`; missing kind counts as question). After the last open item, `focusId = null`. Pass `focusQuestionId={focusId}` to `NoteReader`. Narrow viewport: stack the rail below the reader with CSS, not a second `NoteReader`.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts tests/component/NoteQuestionRail.test.ts tests/component/NoteReader.test.ts` then `npx playwright test tests/e2e/us8-note-rail.spec.ts`

**Landmines:** Directive order, not createdAt. Annotations never in the rail. `/notes/new` stays edit mode with Title/Note labels — do not show the rail on new unsaved notes.

**Checkpoint**: Reading a note, open questions are listed in passage order and Next unanswered walks them.

---

## Phase 15: Today / Next queue (US9) — gap remediation

**Purpose**: One opinionated workspace queue. Not saved filters.

**For a small model.** Read this phase only. Pure helper first, then page, then nav. No new backend endpoint.

**Independent Test**: Seed overdue, due today, in progress, deferred with due ≤ today, high priority unanswered, deferred with future due, annotation, answered. Next shows five sections in spec order; omits the last three kinds of items; opening a row lands on the US8 question page.

**Do not**: add saved filters, query builders, new API routes, or new DB columns.

### Tests for Phase 15 (write first; they must fail)

- [x] T150 [P] [US9] Add `web/tests/unit/nextQueue.test.ts` covering `buildNextQueue` in `web/src/lib/questions/nextQueue.ts` with a fixed `today` of `2026-04-01`:
  - overdue: `dueDate < today` and status not `deferred`
  - due today: `dueDate === today` and status not `deferred`
  - in progress: `in_progress` not already in overdue/due today
  - deferred ready: `deferred` and `dueDate <= today`
  - high priority: `high` or `urgent`, status `unanswered` or `in_progress`, not already listed
  - omit: `answered`, `kind=annotation`, `deferred` with `dueDate > today`, `deferred` with null due date
  - hide empty sections; never duplicate an id across sections
- [x] T151 [P] [US9] Add `web/tests/component/NextQueue.test.ts` only: empty guidance `Nothing in Next. Capture a question from a note, or set a due date.` and the five section headings when those sections have items. Do not write Playwright in this task.
- [x] T151a [P] [US9] Add `tests/e2e/us9-next-queue.spec.ts`: one overdue question appears under Overdue; an answered question does not; clicking opens `/questions/{id}` (query string from T144 is allowed).

### Implementation for Phase 15

`questionsApi.list` already accepts `cursor` and `pageSize` and returns `nextCursor`. Copy workspace loading/error/retry from `web/src/routes/questions/+page.svelte` (`currentWorkspaceId`). There is no `/next` route yet. `AppNavigation.svelte` has Active Questions then Answered inside `{#if $workspaceReady}`.

- [x] T152 [P] [US9] Create `web/src/lib/questions/nextQueue.ts`.
  **Do:** Export `buildNextQueue(questions: Question[], today: string): { id: string; title: string; items: Question[] }[]`. `today` is `YYYY-MM-DD`. Filter to active questions only (`kind` missing or `question`; status `unanswered` | `in_progress` | `deferred`). Apply FR-047. Return only non-empty sections with titles `Overdue`, `Due today`, `In progress`, `Deferred ready`, `High priority`.
  **Verify:** `npm --prefix web test -- --run tests/unit/nextQueue.test.ts`

- [x] T153 [US9] Create presentational `web/src/lib/components/NextQueue.svelte`.
  **Do:** `export let sections: { id: string; title: string; items: Question[] }[] = []`. Render each non-empty section heading and its items. Each item links `/questions/{id}?noteId={linkedNotes[0].id}` when `linkedNotes[0]` exists, else `/questions/{id}`. If `sections` is empty, show `Nothing in Next. Capture a question from a note, or set a due date.` No filter widgets. Match `export let`.
  **Verify:** `npm --prefix web test -- --run tests/component/NextQueue.test.ts`

- [x] T153a [US9] Create `web/src/routes/next/+page.svelte` (first page of questions only).
  **Read:** `web/src/routes/questions/+page.svelte`, `web/src/lib/api/questions.ts`, `web/src/lib/stores/workspace.ts`.
  **Do:** Require `$currentWorkspaceId`. Loading/error/retry required. `questionsApi.list({ workspaceId, status: ['unanswered','in_progress','deferred'], kind: 'question', pageSize: 200 })`. `today = new Date().toISOString().slice(0, 10)`. Pass `buildNextQueue(items, today)` into `NextQueue`. Do not follow `nextCursor` yet.
  **Do not:** add filter widgets; do not log question text.
  **Verify:** `npm --prefix web test -- --run tests/unit/nextQueue.test.ts tests/component/NextQueue.test.ts`

- [x] T153b [US9] Paginate Next.
  **Read:** `web/src/routes/next/+page.svelte` only.
  **Do:** Follow `nextCursor` up to 5 list calls. Concatenate `items`. If `nextCursor` remains after 5 pages, show `Showing the first loaded questions.` above or below the queue.
  **Verify:** `npm --prefix web test -- --run tests/unit/nextQueue.test.ts tests/component/NextQueue.test.ts`

- [x] T154 [US9] Add Next to `web/src/lib/components/AppNavigation.svelte` after the Active Questions `<a>`: `href={`/next${workspace}`}` labelled `Next`, `class:active={page.url.pathname === '/next'}` (do not use `startsWith`). Only inside the existing `{#if $workspaceReady}` branch.
  **Verify:** `npx playwright test tests/e2e/us9-next-queue.spec.ts tests/e2e/us8-answer-in-context.spec.ts`

**Landmines:** Client-side grouping only. Do not add `queue=` to OpenAPI. Annotations never appear (`kind=question` on the list call). Do not log question text.

**Checkpoint**: `/next` is a usable daily queue that opens answer-in-context.

---

## Phase 16: Deferred requires a resume date (US10) — gap remediation

**Purpose**: Status `deferred` requires a due date. No new status. Legacy null due dates still load.

**For a small model.** Domain rule first, then every `ValidateTransition` caller, then UI. Do not migrate old rows.

**Independent Test**: Save as Deferred with no date → rejected. With a date → saved. Next omits future-deferred and includes deferred-ready. Opening a legacy deferred question with null due date still works.

### Tests for Phase 16 (write first; they must fail)

- [x] T155 [P] [US10] Extend `server/internal/domain/question_lifecycle_test.go`: `deferred` with empty/nil due date fails; `deferred` with `YYYY-MM-DD` succeeds; `answered` still requires an answer and does not require a due date; `unanswered` / `in_progress` still allow null due date.
- [x] T156 [P] [US10] Extend `web/tests/component/QuestionLifecycle.test.ts` only: choosing Deferred with no date shows `A resume date is required to defer a question.` and does not call update; with a date, save sends `status: 'deferred'` and that `dueDate`. Do not write Playwright in this task.
- [x] T156a [P] [US10] Add `tests/e2e/us10-defer-date.spec.ts` for the rejection + success path. Do not rewrite `tests/e2e/us3-question-lifecycle.spec.ts`.

### Implementation for Phase 16

`ValidateTransition(from, to QuestionStatus, answer *string)` lives in `server/internal/domain/question_lifecycle.go`. Callers (update both so the build succeeds): `server/internal/store/sqlite/capture.go` `UpdateQuestion` (passes `q.AnswerMarkdown`) and `server/internal/store/sqlite/lifecycle.go` `UpdateLifecycle`. There is **no** `server/internal/service/lifecycle.go`. `server/internal/api/handlers/lifecycle.go` is a comment stub; real PUT is `updateQuestion` in `server/internal/api/handlers/capture.go`. `writeStoreError` already 422s messages containing `required` without field errors. `data-model.md` already has the deferred resume-date row — do not rewrite it. `QuestionSchedule.svelte` already has `#due-date`; do not duplicate that control in T159 — add `#resume-date` on the lifecycle form only.

- [x] T157 [US10] Update `server/internal/domain/question_lifecycle.go` and the two store callers.
  **Read:** `question_lifecycle.go`, `server/internal/store/sqlite/capture.go` `UpdateQuestion`, `server/internal/store/sqlite/lifecycle.go` `UpdateLifecycle`.
  **Do:** `ValidateTransition(from, to, answer, dueDate *string)`. If `to == StatusDeferred`, trimmed due date must be non-empty `YYYY-MM-DD`; else return `a resume date is required to defer a question`. Keep answer-required-for-answered. Pass `q.DueDate` from `UpdateQuestion` and the due date already on the loaded question from `UpdateLifecycle` (that helper does not accept a new due date — still pass `q.DueDate` so deferred-without-date fails if the stored value is empty). Do **not** reject existing rows on read.
  **Do not:** add a column; do not fail import of legacy deferred rows; do not edit `data-model.md`.
  **Verify:** `go -C server test ./internal/domain ./internal/store/sqlite -count=1`

- [x] T158 [US10] Map the resume-date error to a 422 field error.
  **Read:** `server/internal/api/handlers/capture.go` `writeStoreError` and `server/internal/api/problem/problem.go` `Validation` / `FieldError`.
  **Do:** When the error message contains `resume date`, write `problem.Validation` with `FieldError{Field: "dueDate", Message: ...}` (status 422). Leave `sqlite.ErrConflict` as the existing version-conflict 409. Do not create a service package. Do not edit the comment-only `handlers/lifecycle.go`.
  **Verify:** `go -C server test ./internal/domain ./internal/store/sqlite ./internal/api/handlers -count=1`

- [x] T159 [US10] Client-side reject deferred with no date in `web/src/lib/components/QuestionLifecycle.svelte`.
  **Read:** that file only (`save()` already sends `dueDate: question.dueDate`).
  **Do:** If `question.kind === 'annotation'`, skip this rule. If `status === 'deferred'` and there is no trimmed due date (`question.dueDate` in this task), set error `A resume date is required to defer a question.` and `return` without calling `questionsApi.update`.
  **Do not:** edit `QuestionSchedule.svelte`.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts`

- [x] T159a [US10] Resume-date input on the lifecycle form.
  **Read:** `web/src/lib/components/QuestionLifecycle.svelte` only.
  **Do:** `let dueDate = question.dueDate ?? ''`. When `status === 'deferred'`, show `<label for="resume-date">Resume date</label> <input id="resume-date" type="date" bind:value={dueDate}>`. On save, if deferred and `!dueDate.trim()`, use the T159 error and return. Include `dueDate: dueDate.trim() || null` on the existing `questionsApi.update` payload (keep `questionText`, answer, status, priority, tags, `version`).
  **Do not:** edit `QuestionSchedule.svelte`; do not apply the rule to annotations.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts` then `npx playwright test tests/e2e/us10-defer-date.spec.ts`

**Landmines:** Read of null `due_date` must still work. Do not auto-write a date on load. JSON unknown fields 400. Do not log answer text.

**Checkpoint**: Deferred without a date cannot be saved; Next can wake deferred items by date.

---

## Phase 17: Resolved highlights + insert answer into note (US11) — gap remediation

**Purpose**: Answered wraps look resolved. Optional, idempotent “Insert answer into note”. Answering does not mutate the note by itself.

**For a small model.** Helper first, then render attributes, then a button on the card. Do not auto-insert on save.

**Independent Test**: Answer a question → highlight style changes, note body unchanged. Insert → blockquote after wrap. Insert again → body unchanged. Reopen → active highlight style; blockquote remains.

### Tests for Phase 17 (write first; they must fail)

- [x] T160 [P] [US11] Add `web/tests/unit/insertAnswer.test.ts` for `insertAnswerAfterDirective` in `web/src/lib/editor/directives.ts`: wrapped id inserts `\n\n` + blockquote of the answer immediately after `{{/question}}`; multiline answer prefixes each line with `> `; second call returns the same markdown; unknown id returns markdown unchanged; empty answer throws and does not mutate.
- [x] T161 [P] [US11] Add `web/tests/unit/markdown-status.test.ts`: `renderNoteHtml(md, questions)` sets `data-status` and `data-kind` on the mark; raw `{{question:` still absent. Keep `markdown-security.test.ts` unchanged in this task.
- [x] T161a [P] [US11] Extend `web/tests/component/AnnotationCard.test.ts` only: answered question shows Insert answer into note; clicking calls `onInsertAnswer` once. Annotations and unanswered questions do not show the button.
- [x] T162 [P] [US11] Add `tests/e2e/us11-resolved-highlight.spec.ts`: answer, reload note, mark has resolved styling or `data-status="answered"`; insert once; reload; blockquote visible; answer save without insert does not add a blockquote.

### Implementation for Phase 17

`renderNoteHtml(markdown: string)` currently builds `<mark data-annotation-id>` and sanitizes with `ADD_ATTR: ['data-annotation-id']`. `NoteReader` uses `$: html = renderNoteHtml(markdown)` and already has `questions`. `NoteExcerpt` uses `$: html = renderNoteHtml(markdown)` with no status prop. `AnnotationCard` renders `QuestionLifecycle` without `onSave` and has no insert button. `NoteEditor.save()` already rebuilds `questionLinks` from `directiveIds(bodyMarkdown)`. Do not insert inside `QuestionLifecycle.save()`.

- [x] T163 [P] [US11] Add `insertAnswerAfterDirective(markdown: string, id: string, answer: string): string` to `web/src/lib/editor/directives.ts`.
  **Do:** Trim answer; throw if empty. Find the wrapped directive with that lowercase id. If the text after that wrap (skip one run of whitespace) already starts with a `>` line that contains the trimmed answer, return markdown unchanged. Else splice `\n\n` + answer lines each prefixed with `> ` + `\n` immediately after `{{/question}}`. Bare (unwrapped) tokens: insert the same blockquote immediately after the opening token. Reconstruct via `tokenizeDirectives` or index math; do not invent a second wrap.
  **Verify:** `npm --prefix web test -- --run tests/unit/insertAnswer.test.ts`

- [x] T164 [US11] Teach `renderNoteHtml` statuses in `web/src/lib/editor/markdown.ts` only.
  **Do:** `renderNoteHtml(markdown: string, questions: { id: string; status?: string; kind?: string }[] = [])`. On `<mark>` add `data-status` (default `unanswered`) and `data-kind` (default `question`). Add `data-status` and `data-kind` to DOMPurify `ADD_ATTR`. One-argument callers must still compile.
  **Do not:** edit Svelte files in this task.
  **Verify:** `npm --prefix web test -- --run tests/unit/markdown-status.test.ts tests/unit/markdown-security.test.ts`

- [x] T164a [US11] Resolved-highlight CSS.
  **Read:** style blocks in `web/src/lib/editor/NoteReader.svelte` and `web/src/lib/components/NoteExcerpt.svelte` (both already style `mark[data-annotation-id]`).
  **Do:** Add `mark[data-status="answered"]` background `#d1fae5`, border-bottom `#047857`. Leave active marks as they are.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/NoteExcerpt.test.ts`

- [x] T164b [US11] Pass status into render.
  **Read:** `NoteReader.svelte` (`$: html = renderNoteHtml(markdown)`), `NoteExcerpt.svelte`, `QuestionContext.svelte` (`<NoteExcerpt markdown={note.bodyMarkdown} questionId={question.id} ... />`).
  **Do:** `NoteReader`: `$: html = renderNoteHtml(markdown, questions)`. `NoteExcerpt`: `export let status = ''` and pass a one-item questions array (or equivalent) into `renderNoteHtml`. `QuestionContext`: pass `status={question.status}` into `NoteExcerpt`.
  **Verify:** `npm --prefix web test -- --run tests/unit/markdown-status.test.ts tests/component/NoteReader.test.ts tests/component/NoteExcerpt.test.ts tests/component/QuestionContext.test.ts`

- [x] T165 [US11] Insert button on `web/src/lib/components/AnnotationCard.svelte` only.
  **Do:** If `question.kind !== 'annotation'` and `question.status === 'answered'` and `question.answerMarkdown`, show button `Insert answer into note`. `export let onInsertAnswer: (() => Promise<void> | void) | undefined`. Click calls `onInsertAnswer` once. Do not call APIs from the card.
  **Verify:** `npm --prefix web test -- --run tests/component/AnnotationCard.test.ts`

- [x] T165a [US11] Implement insert in `web/src/lib/editor/NoteEditor.svelte`.
  **Read:** `save()` and `captureFromReader` in that file; `insertAnswerAfterDirective` from T163.
  **Do:** `async function insertAnswer(question: Question)`: run `insertAnswerAfterDirective` on `bodyMarkdown`; assign the result; rebuild `links` from `directiveIds` the same way `save()` already does; `await save()`. Do not change question status. Do not wrap this into `QuestionLifecycle`.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts tests/unit/insertAnswer.test.ts`

- [x] T165b [US11] Wire insert through `NoteReader`.
  **Read:** `NoteReader.svelte` (`<AnnotationCard question={openQuestion} passage={openPassage} onClose={...} />`) and the read-mode `NoteReader` tag in `NoteEditor.svelte`.
  **Do:** `NoteReader`: `export let onInsertAnswer: ((question: Question) => Promise<void> | void) | undefined` and pass a thunk into `AnnotationCard`. `NoteEditor`: `onInsertAnswer={insertAnswer}` (or a wrapper that passes `openQuestion`).
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/AnnotationCard.test.ts tests/component/NoteEditor.test.ts`

- [x] T165c [US11] Insert from the question context page.
  **Read:** `web/src/lib/components/QuestionContext.svelte` and `web/src/routes/questions/[questionId]/+page.svelte`.
  **Do:** If the loaded question is answered with an answer and `note` is non-null, show `Insert answer into note` on `QuestionContext` (not inside `QuestionLifecycle.save`). Click: `insertAnswerAfterDirective` on `note.bodyMarkdown`, rebuild `questionLinks` from `directiveIds`, `notesApi.update` with current `note.version`. Skip if unlinked (`note` is null). Do not change question status.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionContext.test.ts` then `npx playwright test tests/e2e/us11-resolved-highlight.spec.ts tests/e2e/us8-answer-in-context.spec.ts`

**Landmines:** `questionLinks` must match directive order after insert (blockquote is not a directive). Idempotent insert. Never log the answer. `/notes/new` labels unchanged.

**Checkpoint**: Answered passages look resolved; insert is optional and safe to click twice.

---

## Phase 18: Sturdier passage matching (US12) — gap remediation

**Purpose**: Wrap the source Markdown span that produced the rendered selection. Fail with no mutation when it cannot be mapped.

**For a small model.** Pure function in `directives.ts` only, then switch `wrapSelection` to it, then show the error in the composer. No schema.

**Independent Test**: Note body `This is **bold** text.` Select rendered `bold`. Capture succeeds and wrap includes `**bold**`. Select `zzz` → error, no new question, body unchanged.

### Tests for Phase 18 (write first; they must fail)

Split the unit file so each matcher slice can pass on its own. `wrapSelection` today does `markdown.indexOf(selectedText)` and throws `selection is empty` / `selection not found in note`.

- [x] T166 [P] [US12] Extend `web/tests/unit/directives.test.ts` with cases that can pass after T168: empty selection throws `selection is empty`; exact substring still wraps the first occurrence; wrap still produces `{{question:id}}…{{/question}}`; `directiveIds` order is preserved when a wrap is added among existing directives.
- [x] T166a [P] [US12] Add cases: markdown `hello   world` + selection `hello world` wraps the original spaced span.
- [x] T166b [P] [US12] Add cases: markdown `**bold**` + selection `bold` wraps `**bold**`; no match throws `Could not find that passage in the note. Try selecting plain text.`
- [x] T167 [P] [US12] Extend `web/tests/component/NoteReader.test.ts` only: when `onCapture` rejects with that error, composer shows it (`composerError` / role=alert) and the success path does not run (composer stays open). Do not write Playwright in this task.
- [x] T167a [P] [US12] Add `tests/e2e/us12-passage-match.spec.ts` for the bold case (note body exactly `This is **bold** text.`).

### Implementation for Phase 18

`captureFromReader` currently `questionsApi.create` then `wrapSelection(bodyMarkdown, passage, question.id)` — that order orphans questions if wrap fails. `submitComposer` already assigns `composerError` from a thrown `Error` and does not close the composer. Do not parse links, headings, or HTML. No DB columns.

- [x] T168 [US12] Add `findSelectionInMarkdown` steps 1–2 in `web/src/lib/editor/directives.ts`.
  **Do:** `findSelectionInMarkdown(markdown: string, selectedText: string): { index: number; length: number }`. (1) Trimmed selection empty → throw `selection is empty`. (2) `markdown.indexOf(selectedText)` ≥ 0 → `{ index, length: selectedText.length }`. Else throw `Could not find that passage in the note. Try selecting plain text.` Do not change `wrapSelection` yet.
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts` (T166 cases only need to keep passing; T166a/T166b may still fail).

- [x] T168a [US12] Whitespace collapse in `findSelectionInMarkdown`.
  **Do:** After exact `indexOf` fails, collapse runs of whitespace to one space on both strings with an index map back to markdown; if the collapsed selection occurs, return mapped `{ index, length }` in the original markdown. First match wins.
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts` (T166 + T166a).

- [x] T168b [US12] Skip markdown markers in `findSelectionInMarkdown`.
  **Do:** After whitespace collapse fails, build a readable string by skipping `*`, `_`, and `` ` `` and collapsing whitespace, mapping each readable index to a markdown index; find the collapsed selection there; return mapped `{ index, length }` so markers stay inside the wrap (`**bold**` not `bold`). Else throw `Could not find that passage in the note. Try selecting plain text.`
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts`

- [x] T168c [US12] Point `wrapSelection` at the finder.
  **Read:** `wrapSelection` in `web/src/lib/editor/directives.ts` only.
  **Do:** Use `findSelectionInMarkdown` and slice `markdown[index, index+length]` as the wrapped span. Keep the same `{{question:id}}…{{/question}}` output. Empty error stays `selection is empty`.
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts`

- [x] T169 [US12] Reorder capture so match happens before create.
  **Read:** `captureFromReader` in `web/src/lib/editor/NoteEditor.svelte` only.
  **Do:** Preferred order: `findSelectionInMarkdown(bodyMarkdown, passage)` → `questionsApi.create` → `wrapSelection` with the returned id → `save()`. If match/wrap throws, do not create. If save fails, keep existing failed-save draft behavior. Do not rewrite `submitComposer` — it already surfaces `cause.message` as `composerError`.
  **Verify:** `npm --prefix web test -- --run tests/unit/directives.test.ts tests/component/NoteReader.test.ts tests/component/NoteEditor.test.ts` then `npx playwright test tests/e2e/us12-passage-match.spec.ts tests/e2e/us1-notes-highlight.spec.ts`

**Landmines:** Create-after-failed-wrap would orphan questions — match before create. `questionLinks` order. Do not log the passage.

**Checkpoint**: Formatted-text selections wrap source Markdown; unmappable selections mutate nothing.

---

## Phase 19: Link existing + attach later from reading (US12) — gap remediation

**Purpose**: Reading toolbar can link an existing same-workspace question. Active Questions can create an unlinked question. Reading can attach a passage later.

**For a small model.** Reuse `QuestionPicker.svelte`. POST `/api/v1/questions` already creates without a note. Do not add endpoints.

**Independent Test**: Question A on note 1. From note 2 reading, link A onto a sentence — one question, two linked notes. From Active Questions create unlinked B. From a note, attach B to a passage — B is no longer unlinked.

**Do not:** link a question already on this note; link annotations; change kind; allow cross-workspace links.

### Tests for Phase 19 (write first; they must fail)

Reading toolbar today is Ask a question / Add annotation / Cancel. `QuestionPicker` already takes `questions` + `onSelect` and shows `Unlinked` when `linkedNotes` is empty. Edit-mode `selectExisting` in `NoteEditor` inserts a **bare** token at the end — do not reuse that for reading wrap. `web/src/routes/questions/+page.svelte` has no New question control. `ActiveQuestions.svelte` is only a `QuestionList` wrapper.

- [x] T170 [P] [US12] Extend `web/tests/component/NoteReader.test.ts` only: toolbar has `Link existing question`; it opens the picker; choosing an item calls `onLinkExisting(passage, question)`.
- [x] T170a [P] [US12] Extend `web/tests/component/QuestionPicker.test.ts` only if T170 needs picker behavior that is not already covered. Skip this task when existing picker tests already prove search + `onSelect` + Unlinked copy.
- [x] T170b [P] [US12] Add or extend `web/tests/component/ActiveQuestions.test.ts` (test the Active Questions page component if easier): `New question` without a passage posts only question fields (no note wrap).
- [x] T171 [P] [US12] Add `tests/e2e/us12-link-attach.spec.ts`: (1) link existing question onto a second note from reading; both notes listed on the question page; (2) create unlinked from Active Questions; it appears there; attach to a passage from reading; highlight opens the same question.

### Implementation for Phase 19

- [x] T172 [US12] Reading toolbar + picker UI in `web/src/lib/editor/NoteReader.svelte` only.
  **Read:** that file and `web/src/lib/components/QuestionPicker.svelte`.
  **Do:** `export let onLinkExisting: ((passage: string, question: Question) => Promise<Question | void>) | undefined`. `export let linkableQuestions: Question[] = []`. Toolbar button `Link existing question` (only useful when `onLinkExisting` is set). Opens `QuestionPicker` with `linkableQuestions`. Choosing an item calls `onLinkExisting(selectedPassage, question)`. Keep Cancel/Escape/outside-click dismissal for the picker the same way the composer is dismissed. Do not create or wrap here.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/QuestionPicker.test.ts`

- [x] T172a [US12] Wrap-without-create in `web/src/lib/editor/NoteEditor.svelte`.
  **Read:** `captureFromReader`, `selectExisting` (edit-mode bare insert — do not call it), `availableQuestions`, `directiveIds`.
  **Do:** Pass `linkableQuestions` = workspace questions with `kind: 'question'` excluding ids already in `directiveIds(bodyMarkdown)`. `linkExistingFromReader(passage, question)`: if `directiveIds` already contains `question.id`, throw `That question is already in this note.` and do not save. Else `findSelectionInMarkdown` / `wrapSelection` first (no `questionsApi.create`); rebuild `links` from `directiveIds`; `save()`; return the same question. Pass `onLinkExisting={linkExistingFromReader}`.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/NoteEditor.test.ts tests/component/QuestionPicker.test.ts`

- [x] T173 [US12] Unlinked create on `web/src/routes/questions/+page.svelte`.
  **Do:** Button labelled `New question`. Short form: question text required. `questionsApi.create({ workspaceId, questionText, kind: 'question', status: 'unanswered', priority: 'none', tagIds: [] })`. No note update. Reload the list. Empty text does not POST.
  **Do not:** set status answered; do not attach a fake directive; do not edit `NoteEditor`.
  **Verify:** `npm --prefix web test -- --run tests/component/ActiveQuestions.test.ts` (or the page test added in T170b)

- [x] T174 [US12] Keep unlinked questions in the reading picker.
  **Read:** `QuestionPicker.svelte` (already renders `Unlinked`) and the `linkableQuestions` filter in `NoteEditor.svelte`.
  **Do:** Do not filter out questions with empty `linkedNotes`. After attach, Active Questions still shows one row for that id (no second create). If T172a already does this, only run verify.
  **Verify:** `npx playwright test tests/e2e/us12-link-attach.spec.ts tests/e2e/us2-shared-question.spec.ts tests/e2e/us1-notes-highlight.spec.ts`

**Landmines:** Same workspace only (list is already workspace-scoped). `questionLinks` 1:1 with directives. Kind immutable. Do not wrap with an annotation id from this picker.

**Checkpoint**: Shared questions and later-attached passages work from reading view.

---

## Phase 20: Keyboard-first capture (US12) — gap remediation

**Purpose**: Selection shortcuts on the reading view; do not steal keys while typing.

**For a small model.** Edit `NoteReader` keydown only, plus document shortcuts. Do not add a keymap framework.

**Independent Test**: Select text, press Q → question composer; Escape → dismiss; select, A → annotation composer; select, L → link picker. Type the letter Q inside the composer — it inserts Q. `/notes/new` labels unchanged.

### Tests for Phase 20 (write first; they must fail)

`onWindowKeydown` today handles **Escape only** (close composer or toolbar).

- [x] T175 [P] [US12] Extend `web/tests/component/NoteReader.test.ts` only: with a selection and no composer, `keydown` Q opens question composer, A annotation, L picker (if `onLinkExisting` provided), Escape closes toolbar. With composer open or target `input`/`textarea`, Q does not toggle kind.
- [x] T175a [P] [US12] Add `tests/e2e/us12-keyboard-capture.spec.ts` for Q then Escape on a selected sentence.

### Implementation for Phase 20

- [x] T176 [US12] Extend `onWindowKeydown` in `web/src/lib/editor/NoteReader.svelte` only.
  **Do:** Keep the existing Escape branch. Then: if `event.defaultPrevented`, return. If target is `input, textarea, select, [contenteditable]`, return (Escape may still close composer). If composer open, ignore Q/A/L. If no `selectedPassage` and no `showToolbar`, ignore Q/A/L. Otherwise: `q`/`Q` → `openComposer('question')`; `a`/`A` → `openComposer('annotation')`; `l`/`L` → open the T172 picker if `onLinkExisting` exists. `preventDefault` only when handling those keys. Do not handle shortcuts with Ctrl/Meta/Alt.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts` then `npx playwright test tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts`

- [x] T177 [P] [US12] Add a one-line hint on the reading toolbar in `NoteReader.svelte`: `Q ask · A annotate · L link · Esc cancel`. No new route.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts` — assert the hint exists when the toolbar is visible.

**Landmines:** Do not steal keys in Title/Note fields. Do not change `/notes/new`. Keep existing Escape/outside-click dismissal.

**Checkpoint**: Capture is keyboard-complete without breaking typing or accessibility labels.

---

## Phase 20a: In-context floating capture toolbar and mobile-friendly highlight positioning (UX gap remediation)

**Purpose**: Address user experience feedback on long notes and mobile devices:
1. When notes are long, the user currently has to scroll down past the note body just to see the selection toolbar and action buttons. The action controls must be easily accessible without keyboard shortcuts.
2. The selection toolbar must automatically show up immediately adjacent to the highlighted/selected text (floating or anchored to the selection range) so mobile users and mouse users see the capture actions in their current viewport.
3. The question/annotation composer and highlight card must not be stranded at the bottom of a long document below the article; they should appear in-context right next to or anchored beside the highlighted passage (with responsive fallback on small/mobile screens so they do not clip off-screen).

**For a small model.** Read `PROJECT_MAP.md` once, this phase only, then only `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/AnnotationCard.svelte`, and their test files. Implement **one task at a time**, run its verify command, and stop.

**Independent Test**:
1. Open a long note (several screens of text).
2. Select text near the top or middle of the note without scrolling to the bottom. The selection toolbar ("Ask a question", "Add annotation", "Link existing question", "Cancel") automatically appears floating directly adjacent to (above or below) the highlighted passage in view.
3. Tap or click "Ask a question" or "Add annotation". The composer opens adjacent to the highlighted selection (not appended at the bottom of the long article).
4. Save the question. Click the newly created highlight. The question card/dialog opens next to the clicked highlight in view (or as a bottom sheet / centered sheet on mobile viewports), avoiding any need to scroll down to view or edit the card.
5. On mobile/touch devices or narrow viewports, the floating toolbar and composer remain fully visible within the viewport (clamped to viewport boundaries) without horizontal overflow or off-screen clipping.

**Do not**:
- Change backend APIs, domain validation, SQLite schemas, or OpenAPI contracts.
- Alter Markdown directives grammar (`{{question:uuid}}...{{/question}}`) or `findSelectionInMarkdown`.
- Break existing Q/A/L/Escape keyboard shortcuts or tests.
- Change `/notes/new` Title/Note labels or layout.

### Tests for Phase 20a (write first; they must fail)

- [x] T177a [P] [US12] Extend `web/tests/component/NoteReader.test.ts`:
  - When text is selected in the reader, the toolbar element receives floating/positioned styling (e.g., inline coordinates or positioning classes based on selection bounds) rather than rendering statically at the bottom of the article.
  - When the composer is opened from a selection, the composer container retains the contextual positioning adjacent to the selection.
  - Viewport boundary safety: clamped coordinates keep toolbar and composer inside visible container bounds.
- [x] T177b [P] [US12] Extend `web/tests/component/AnnotationCard.test.ts`:
  - Card supports positioning anchored near the target highlight element or range when coordinates or target rect are provided, while retaining its responsive modal/fixed overlay behavior when positioned near edges or on mobile screens.
- [x] T177c [P] [US12] Add Playwright test `tests/e2e/us12-floating-selection-ux.spec.ts`:
  - Create a long note (multiple paragraphs requiring scrolling).
  - Select text near the top of the viewport; verify the selection toolbar is immediately visible within the viewport without scrolling down.
  - Click "Ask a question"; verify the composer is visible within the viewport next to the selection.
  - Create the question; scroll to a highlight; click it; verify the question card is visible in the viewport near the highlight.
  - Test on mobile viewport preset (`viewport: { width: 375, height: 667 }`) to ensure touch-friendly access and no clipping.

### Implementation for Phase 20a

- [x] T177d [US12] Implement selection range coordinate calculation in `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/editor/NoteReader.svelte` selection handling (`onMouseUp`, `onContextMenu`, `selectedText`).
  **Do:**
  - In `onMouseUp` / `onContextMenu`, obtain `window.getSelection()?.getRangeAt(0)?.getBoundingClientRect()`.
  - Calculate relative top and left coordinates relative to the reader container (or viewport-clamped fixed coordinates) for the selection.
  - Store `toolbarPosition = { top: number, left: number }`.
  - Add `selectionchange` / `touchend` listeners on the reader element so mobile touch highlights automatically capture the selection bounds and show the toolbar without requiring mouse events or right-click.
  - Clamp coordinates so the toolbar never renders off-screen or outside the reader container.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T177e [US12] Float toolbar and composer next to selection in `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/editor/NoteReader.svelte` markup and `.toolbar`, `.composer`, `.picker-layer` styles.
  **Do:**
  - Apply contextual positioning (absolute or fixed anchored to `toolbarPosition`) to `.toolbar`, `.composer`, and `.picker-layer`.
  - Style the toolbar as a floating popover/bubble above (or below) the selection with a subtle shadow and pointer arrow or backdrop, ensuring touch targets are minimum 44px for mobile friendliness.
  - When switching from toolbar to composer or picker, retain the anchor position so the composer opens in place next to the selected passage instead of displacing to document end.
  - Add `@media (max-width: 640px)` mobile overrides: if the selection is near edges, dock the toolbar/composer gracefully (e.g. pinned bottom bar or clamped floating sheet) so it remains 100% accessible on small screens.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T177f [US12] Anchor highlight card next to clicked highlight in `web/src/lib/components/AnnotationCard.svelte` and `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/components/AnnotationCard.svelte` style block (`position: fixed; inset: auto 1rem 1rem auto;`) and `NoteReader.svelte` `onClick`.
  **Do:**
  - In `NoteReader.svelte` `onClick(event)`, record the clicked mark's bounding rect or relative position (`openAnchorRect` or `cardPosition`).
  - Pass anchor position or style prop to `AnnotationCard`.
  - In `AnnotationCard.svelte`: on wide viewports (desktop/tablet), position the card anchored adjacent to the clicked highlight mark (or in the adjacent margin/gutter aligned with the mark's Y-coordinate), clamped to avoid overflowing the viewport.
  - On narrow/mobile screens (`max-width: 640px`), fallback to a bottom sheet overlay so it never clips off the side of small screens.
  **Verify:** `npm --prefix web test -- --run tests/component/AnnotationCard.test.ts tests/component/NoteReader.test.ts`

- [x] T177g [US12] Verify e2e and accessibility regressions.
  **Read:** `tests/e2e/us12-floating-selection-ux.spec.ts`, `tests/e2e/us1-notes-highlight.spec.ts`, `tests/e2e/accessibility-notes.spec.ts`.
  **Do:** Run all note reader Playwright tests on desktop and mobile viewports. Ensure existing click targets, dialog accessibility roles, and escape dismissal remain intact.
  **Verify:** `npx playwright test tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/accessibility-notes.spec.ts`
  **Result:** Chromium passed all five tests serially (`--project=chromium --workers=1`). Firefox and WebKit were unavailable because their Playwright browser executables are not installed in this environment.

### UX Remediation: Floating Toolbar Occlusion and Selection Persistence

**Problem**:
1. When selecting text to highlight, the floating toolbar renders immediately over/adjacent to the selected text, obstructing reading and further text selection while dragging or highlighting.
2. The toolbar does not persist predictably: live `selectionchange` triggers mid-drag, clicks, or pointer micro-movements prematurely dismiss or clear selection state before user action.

**Tests for Toolbar Occlusion & Persistence (write first; they must fail)**:

- [x] T177h [P] [US12] Add unit & component tests for toolbar non-occlusion and persistence in `web/tests/component/NoteReader.test.ts`:
  - Verify toolbar is positioned outside the active line/selection boundary (preferring above the selection start or below selection end with clearance, rather than overlaying the highlighted text).
  - Verify toolbar does not trigger while mouse drag selection is still active (`mousedown` without `mouseup`), avoiding obstructing text mid-selection.
  - Verify toolbar and selection state persist when clicking inside the reader or moving focus until explicitly dismissed (via Escape, Cancel button, or starting a new distinct selection).
  - Verify viewport boundary handling: when selection is at the top of the viewport, toolbar flips below the selection with safe clearance without clipping or obscuring the selected line.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T177i [P] [US12] Add Playwright test in `tests/e2e/us12-floating-selection-ux.spec.ts`:
  - Select a multi-line passage and verify the floating toolbar bounding rect does not overlap the bounding rect of the selected text range.
  - Verify the toolbar persists after pointer release without flickering or dismissing prematurely until an explicit action or dismiss occurs.
  **Verify:** `npx playwright test tests/e2e/us12-floating-selection-ux.spec.ts`

**Implementation for Toolbar Occlusion & Persistence**:

- [x] T177j [US12] Prevent mid-drag toolbar popup and premature clearing in `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/editor/NoteReader.svelte` (`updateSelectionPosition`, `onSelectionChange`, `onMouseUp`, `dismissFromOutside`).
  **Do:**
  - Track pointer selection drag state (e.g., `isSelecting` flag set on reader `mousedown` and cleared on `mouseup`).
  - Suppress `onSelectionChange` toolbar updates while pointer drag is active so the toolbar does not flash or obstruct while the user is actively dragging.
  - Refine `dismissFromOutside` and selection change listeners so existing active selection and toolbar do not prematurely collapse on minor mouse jitter or non-dismissing clicks.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T177k [US12] Position floating toolbar with non-occluding clearance in `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/editor/NoteReader.svelte` `calculateToolbarPosition` and positioning styles.
  **Do:**
  - Update `calculateToolbarPosition` to prefer placing the toolbar above the selection range (`rect.top - layerHeight - clearanceGap`) when space permits, so the selected passage remains unobstructed.
  - If placing below (when near top of viewport), anchor below `rect.bottom + clearanceGap` to prevent overlapping the text lines being read.
  - Horizontally align toolbar relative to selection start or center, clamped within reader/viewport padding so it doesn't overlap text awkwardly or clip viewport edges.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T177l [US12] Ensure selection persistence and dismiss UX in `web/src/lib/editor/NoteReader.svelte`.
  **Read:** `web/src/lib/editor/NoteReader.svelte` selection lifecycle and toolbar actions.
  **Do:**
  - Preserve the browser DOM selection or visual highlight indicator while the toolbar or composer is active.
  - Allow explicit dismissal via Escape key, "Cancel" button, or clicking outside without accidental dismissals when interacting with nearby reader controls.
  **Verify:** `npx playwright test tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/accessibility-notes.spec.ts`

**Landmines:**
- Selection coordinates can collapse or return 0 when clicking outside; verify `range.getBoundingClientRect()` is valid before updating position.
- Do not let fixed or absolute popovers cause horizontal scrollbars or overflow on mobile devices.
- Outside-click dismissal (`dismissFromOutside`) must continue to ignore clicks inside the floating toolbar/composer elements.
- Keep accessibility roles (`role="toolbar"`, `role="dialog"`) and keyboard trap / escape dismissal functional.

**Checkpoint**: Users reading long notes see the action toolbar, composer, and highlight card appear right beside their highlighted text on desktop and touch-friendly on mobile without scrolling down.

---

## Phase 21: Map, spec trace, and regression (after US8–US12)

**Purpose**: Keep the mental model honest after the new slices.

- [x] T178 [P] Refresh [`PROJECT_MAP.md`](../../PROJECT_MAP.md): remaining-gaps list must match reality; add file pointers for `QuestionContext.svelte`, `NoteExcerpt.svelte`, `NoteQuestionRail.svelte`, `nextQueue.ts`, `/next`, `findSelectionInMarkdown`, and keyboard shortcuts. Known-bugs table: only real failures.
- [x] T179 Trace FR-041 through FR-047 to tests in `specs/001-track-note-questions/traceability.md`. Do not invent passing evidence.
- [x] T179a Trace FR-048 through FR-055 in the same file. Do not invent passing evidence.
- [x] T179b Trace SC-011 through SC-014 in the same file. Do not invent passing evidence.
- [x] T180 Run the new Playwright files plus `tests/e2e/us1-notes-highlight.spec.ts`, `tests/e2e/us3-question-lifecycle.spec.ts`, `tests/e2e/us2-shared-question.spec.ts`, `tests/e2e/accessibility-notes.spec.ts`. Record results in `specs/001-track-note-questions/validation-results.md` only for runs you actually executed.

**Checkpoint**: A new session can implement leftover tasks from the map without rediscovering the tree.

---

## Phase 21a: Capture composer focus hotfix (US12) — emergency UX regression remediation

**Purpose**: Fix the regression introduced by Phase 20a. The floating toolbar/composer currently reacts to pointer and selection events generated while a user clicks the question input or annotation textarea. The action layer can be dismissed, rerendered, or have its selection restored on that same interaction, so browser focus immediately leaves the textbox. Restore normal native focus and typing while retaining the floating anchor and explicit dismissal behavior.

**For a small model.** Read `PROJECT_MAP.md` once, this phase only, then only `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/AnnotationCard.svelte`, and the test files named below. Implement **one task at a time**, run its verify command, and stop. This is an emergency client-only fix that must land before Phase 22.

**Independent Test**:
1. Open a long note and select text. Click "Ask a question". Click inside the `Question text` input and type several characters; the input keeps focus, the composer stays open and anchored, and the typed value remains visible.
2. Dismiss that composer, select another passage, choose "Add annotation", and repeat with the `Annotation` textarea. Saving after typing submits the exact text instead of cancelling or submitting an empty value.
3. Open an existing highlight card and interact with any editable answer/annotation control it contains; focus and typed content remain stable. Escape and an intentional outside click still dismiss the appropriate layer.
4. Repeat the capture interactions at a 375×667 viewport. No overlay clips the controls or causes horizontal scrolling.

**Do not**:
- Change backend APIs, domain validation, SQLite schemas, OpenAPI contracts, directive grammar, or `findSelectionInMarkdown`.
- Remove floating positioning, selection persistence, Q/A/L/Escape shortcuts, or outside-click dismissal.
- Fix the issue by globally disabling dismissal, preventing default on every `mousedown`, or adding a second composer/card.
- Change question/annotation validation or save semantics.

### Tests for Phase 21a (write first; they must fail)

- [x] T180a [P] [US12] Extend `web/tests/component/NoteReader.test.ts` with the capture-composer focus regression: use a real user pointer sequence after selecting a passage, click "Ask a question", click/fill the `Question text` input, and assert that it remains `document.activeElement`, the composer remains mounted, and the value survives window `mousedown`/`click` and `selectionchange` events. Repeat for "Add annotation" and its textarea. On save, assert `onCapture` receives the entered text. Do not rely on programmatic `.focus()` alone.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T180b [P] [US12] Extend `web/tests/component/NoteReader.test.ts` only: open a positioned question highlight card and an annotation highlight card, then click/fill the card's editable Answer or Annotation control. Assert that focus and typed content remain stable, `onClose` is not invoked by the internal interaction, and an intentional outside click still closes it. Preserve the existing anchored and mobile-safe assertions.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T180c [P] [US12] Add `tests/e2e/us12-capture-focus.spec.ts`: on desktop and a 375×667 mobile viewport, select a passage, open the question composer, fill the question textbox, and save; repeat with the annotation textarea. Assert the composer does not disappear while typing, the saved values are used, and Escape/cancel plus an intentional outside click still work. Keep the existing floating-toolbar and accessibility specs unchanged.
  **Verify:** `npx playwright test tests/e2e/us12-capture-focus.spec.ts`

### Implementation for Phase 21a

- [x] T180d [US12] Make `web/src/lib/editor/NoteReader.svelte` focus-safe for internal pointer events.
  **Read:** `dismissFromOutside`, `onWindowMouseDown`, `onWindowClick`, `onSelectionChange`, `updateSelectionPosition`, and the toolbar/composer/picker markup.
  **Do:** Add one event-path/active-element guard that recognizes descendants of `.toolbar`, `.composer`, `.picker-layer`, and the anchored `.card` as internal interactive content. Check it before the reader/outside branches. Pointerdown, click, and selectionchange events from those controls must not call `clearSelection`, recompute the anchor, or call `restorePreservedSelection` in a way that moves focus; do not call `preventDefault` for normal control interaction. While a composer, picker, or card control owns focus, keep `selectedPassage`, `preservedSelectionRange`, and `toolbarPosition` stable. Keep Escape and genuine outside clicks cancel-safe.
  **Do not:** remove the existing drag-selection guard, change capture ordering, or disable outside dismissal for the rest of the document.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [x] T180e [US12] Stabilize the toolbar-to-composer/picker transition in `web/src/lib/editor/NoteReader.svelte` and, only if needed for the interaction boundary, `web/src/lib/components/AnnotationCard.svelte`.
  **Do:** Do not rely on a transient shared `bind:this` value while the toolbar is replaced by the composer or picker; use a stable event-path marker or equivalent refs. Preserve the selected passage, floating coordinates, and preserved range when an action button opens a layer. Ensure the card and nested lifecycle/annotation controls receive normal pointer events and are treated as inside the interactive layer without adding a focus trap or repeatedly forcing focus. Keep the existing desktop anchor and mobile bottom-sheet fallback.
  **Do not:** move the controls back to the end of the article or add a second overlay.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/AnnotationCard.test.ts`

- [x] T180f [US12] Verify the emergency focus fix and regressions.
  **Read:** `tests/e2e/us12-capture-focus.spec.ts`, `tests/e2e/us12-floating-selection-ux.spec.ts`, `tests/e2e/us12-keyboard-capture.spec.ts`, and `tests/e2e/accessibility-notes.spec.ts`.
  **Do:** Run the new focus spec on desktop and mobile plus the existing floating-toolbar, keyboard-capture, highlight, and accessibility specs. Confirm typing in the question/annotation controls, Q/A/L/Escape behavior, dialog roles, outside-click dismissal, and mobile viewport safety. Record only results from tests actually executed.
  **Verify:** `npx playwright test tests/e2e/us12-capture-focus.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts`

**Landmines**:
- A `svelte:window` `mousedown` handler runs during the browser's focus sequence; never prevent the default focus action for a textbox/button inside the floating layers.
- `selectionchange` may fire when an input or textarea receives focus. Do not interpret that event as a new empty selection while a capture layer or card control is active.
- The toolbar, composer, picker, and card can replace one another in the same event turn. Use the event path or a stable marker rather than only a mutable `bind:this` reference.
- Keep the existing native selection/visual anchor, but do not repeatedly remove and re-add ranges while the user is typing.
- Preserve minimum touch targets, fixed-position clamping, dialog accessibility, and explicit Escape/outside-click dismissal.

**Checkpoint**: Selecting a passage, opening either capture action, and clicking or typing in its textbox works on desktop and mobile. The floating layer stays anchored until save or explicit dismissal, while outside-click and keyboard regressions remain covered.

---

## Phase 21b: Note-rail card-opening dismissal hotfix (US8) — emergency UX regression remediation

**Purpose**: Fix a regression exposed after Phase 21a. `NoteQuestionRail` is a sibling of `NoteReader`, not a descendant of its reader/layer DOM. Clicking a question on the right rail runs `selectQuestion`, which sets `focusQuestionId`; `NoteReader` opens the card, but the same pointer/click path is then handled by its `svelte:window` outside-dismissal logic. Because the rail target is outside `readerElement` and is not marked as an internal layer, the newly opened `openQuestion` can be dismissed immediately. The result is that questions on the right side of the note appear not to open. Treat intentional rail navigation as an internal event without weakening genuine outside-click dismissal.

**For a small model.** Read `PROJECT_MAP.md` once, this phase only, then only `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/NoteQuestionRail.svelte`, and the test files named below. Implement one task at a time, run its verify command, and stop.

**Independent Test**:
1. Open an existing note with two open questions. Click the first question button in the right-side rail with a real pointer sequence; its question card remains open after the complete click event and its answer controls are usable.
2. Click the second rail question and then `Next unanswered`; each intentional rail navigation opens the requested card and does not immediately dismiss it.
3. Clicking a highlight in the article still opens its card, and an intentional click outside the reader/card still dismisses it. Toolbar/composer/picker focus behavior from Phase 21a remains unchanged.
4. Repeat at a 375×667 viewport without clipping or horizontal scrolling.

**Do not**:
- Change backend APIs, domain validation, SQLite schemas, OpenAPI contracts, directive grammar, or question lifecycle behavior.
- Remove the Phase 21a internal-layer, drag-selection, keyboard, or native-focus guards.
- Globally ignore outside clicks, add a timeout/race-prone deferred dismissal, or close/reopen cards from the rail by creating a second overlay.
- Treat every click in the document as internal; genuine outside clicks must remain cancel-safe.

### Tests for Phase 21b (write first; they must fail)

- [ ] T180g [P] [US8] Extend `web/tests/component/NoteReader.test.ts` with the rail-navigation dismissal regression: use an external rail-like button/control to change `focusQuestionId`, exercise the real bubbling `mousedown`/`click` sequence, and assert the requested question dialog remains mounted after the window handlers run. Cover `Next unanswered`-style navigation, repeated navigation to another question, and an intentional outside click that still closes the card. Keep the existing highlight-card and Phase 21a focus assertions unchanged.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts`

- [ ] T180h [P] [US8] Extend `tests/e2e/us8-note-rail.spec.ts`: click a specific question in the right-side rail, assert its card stays visible, then exercise `Next unanswered`; verify the existing highlight click and outside-dismissal paths still work on desktop and at a 375×667 viewport. Do not rewrite the existing rail ordering/assertion coverage.
  **Verify:** `npx playwright test tests/e2e/us8-note-rail.spec.ts`

### Implementation for Phase 21b

- [ ] T180i [US8] Mark note-rail navigation as an intentional internal interaction and preserve card opening.
  **Read:** `web/src/lib/editor/NoteReader.svelte` `isInternalInteractiveTarget`, `isInternalInteractiveEvent`, `dismissFromOutside`, and the `focusQuestionId` reactive block; `web/src/lib/components/NoteQuestionRail.svelte` root markup.
  **Do:** Add a stable, semantic navigation marker to the rail root (for example `data-note-reader-navigation`) and recognize that marker through the event target/composed path before the `readerElement`/`openQuestion` outside branches. A rail button or `Next unanswered` click must not clear the card opened by the resulting `focusQuestionId` update. Keep the marker narrow to the note-question rail; preserve normal pointer events and focus. Do not reset `focusQuestionId` as a workaround, move the rail into `NoteReader`, or alter card anchoring.
  **Do not:** bypass outside dismissal for arbitrary siblings or call `preventDefault` on rail controls.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/NoteQuestionRail.test.ts`

- [ ] T180j [US8] Verify the note-rail dismissal fix and Phase 21a regressions.
  **Read:** `tests/e2e/us8-note-rail.spec.ts`, `tests/e2e/us12-capture-focus.spec.ts`, `tests/e2e/us12-floating-selection-ux.spec.ts`, `tests/e2e/us12-keyboard-capture.spec.ts`, and `tests/e2e/accessibility-notes.spec.ts`.
  **Do:** Run the rail regression on desktop and mobile plus the capture-focus, floating-toolbar, keyboard, highlight, and accessibility specs. Record only tests actually executed.
  **Verify:** `npx playwright test tests/e2e/us8-note-rail.spec.ts tests/e2e/us12-capture-focus.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts`

**Landmines**:
- The rail is outside `readerElement`, so checking only `readerElement.contains(event.target)` will continue to classify its buttons as outside.
- Both `svelte:window` `mousedown` and `click` handlers can observe the same rail gesture; guard the event path before the open-card dismissal branch.
- Do not suppress the `focusQuestionId` update or repeatedly restore the browser selection while a rail card is opening.
- Preserve `role="dialog"`, Escape dismissal, card controls, and mobile bottom-sheet behavior.

**Checkpoint**: Clicking a question or Next unanswered in the right-side rail opens the requested card and keeps it open, while highlight clicks and genuine outside dismissal continue to work.

---

## Phase 22: Wire Delete note (US6) — gap remediation

**Purpose**: The note page already shows a Delete note button, but clicking it does nothing. Make that existing control delete the current note (cancel-safe), then leave the note page.

**For a small model.** Do not scan the repo. Do not reread US1–US12 or Phases 11–21a internals. Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement **one task**, run its verify command, stop.

**Independent Test**: Create a workspace and a note. Open the note. Click Delete note, then cancel — the note is unchanged and still listed under Notes. Click Delete note again and confirm — the app leaves the editor, `/notes` no longer lists that note, and opening its old URL does not show the editor. A note that has questions still deletes; the user is not stuck on a dead button.

**Do not**: add a second Delete note button; change capture, highlights, `/notes/new` labels, Next queue, or question lifecycle rules; rebuild the Phase 8 deletion backend unless the client has no delete operation to call.

### Tests for Phase 22 (write first; they must fail)

The Delete note **button is on** `web/src/routes/notes/[noteId]/+page.svelte`, not inside `NoteEditor.svelte`. `notesApi.previewDeletion` and `notesApi.deleteReviewed` already exist. `DeleteNoteReview.svelte` already exists. Fill gaps; do not add a second button.

- [ ] T181 [P] [US6] Extend `web/tests/component/DeleteNoteReview.test.ts` (and a note-page test if one exists; do not put the button into `NoteEditor.test.ts`): an existing saved note's Delete note control opens `DeleteNoteReview`; confirming calls delete / `onConfirm` once; cancel does not call delete. `/notes/new` must not expose a working delete for an unsaved note.
- [ ] T182 [P] [US6] Add Playwright `tests/e2e/us6-delete-note-button.spec.ts`: create a workspace and a uniquely titled note; open it; Delete note → cancel → still on the note page and the title still appears on `/notes`; Delete note → confirm → land on `/notes` without that title. Keep `tests/e2e/accessibility-notes.spec.ts` passing.

### Implementation for Phase 22

- [ ] T183 [P] [US6] Confirm the deletion client in `web/src/lib/api/notes.ts`.
  **Read:** that file (already has `previewDeletion` and `deleteReviewed`) and `specs/001-track-note-questions/contracts/openapi.yaml` Phase 8 operations.
  **Do:** If those methods already match the contract, stop. Only add a method when an OpenAPI operation is actually missing. Do not invent a new path.
  **Do not:** add a backend route; do not log note bodies.
  **Verify:** `npm --prefix web test -- --run tests/component/DeleteNoteReview.test.ts`

- [ ] T184 [US6] Keep the existing Delete note button on the note page.
  **Read:** `web/src/routes/notes/[noteId]/+page.svelte` only (`<button class="delete" onclick={reviewDelete}>Delete note</button>` and `{#if preview}<DeleteNoteReview ... />`). Do **not** edit `NoteEditor.svelte`.
  **Do:** If `onclick={reviewDelete}` is missing, attach it to that existing button. Confirm uses preview token + `deleteReviewed`. Cancel sets `preview = null` with no API call. Do not add another Delete note button. `/notes/new` does not have this button — leave it that way.
  **Do not:** delete on the first click with no confirm.
  **Verify:** `npm --prefix web test -- --run tests/component/DeleteNoteReview.test.ts`

- [ ] T185 [US6] After a successful delete, leave the note page.
  **Read:** `web/src/routes/notes/[noteId]/+page.svelte` `confirmDelete`, `web/src/lib/stores/query.ts` `invalidate`, `web/src/lib/stores/toast.ts`.
  **Do:** The page already toasts `Note deleted` and `goto('/notes')`. Add `invalidate` for note/question list keys if that call is missing. On 409: show the conflict and do not navigate. On other errors: stay on the page. Do not invent a new deletion policy.
  **Verify:** `npx playwright test tests/e2e/us6-delete-note-button.spec.ts tests/e2e/accessibility-notes.spec.ts`

**Landmines:** Do not add a second Delete note button. `/notes/new` Title/Note labels unchanged. `questionLinks` are irrelevant after the note is gone — do not save the note as part of delete. Optimistic `version` must be sent. Never log note/question/answer bodies. Cancel must mutate nothing.

**Checkpoint**: The existing Delete note button deletes the current note after confirm and does nothing on cancel.

---

## Phase 23: Edit canonical question text (dog-food gap remediation)

**Purpose**: Close the dog-food gap where question records can be updated by the API but the user-facing question surfaces do not expose an obvious way to edit the question text. The edit must update the one canonical question, not create a replacement or alter note directives.

**For a small model.** Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement one task at a time and run its verify command.

**Independent Test**: Create a question on a note, open it from Active Questions, edit and save its text, and reload the active list, question page, and source-note highlight card. The revised text appears everywhere while the question ID, answer, status, priority, due date, tags, and links remain unchanged. Repeat the edit from the linked note's card. Cancel and blank-text attempts make no request, and a version conflict keeps the draft and does not overwrite newer data.

**Do not**: create a second question, change `questionLinks` or directive order, change lifecycle rules, silently drop answer/status/schedule/tag fields, or broaden this phase to a new consequence field. The existing question update endpoint is the source of truth; do not add a route unless the current contract is actually missing `questionText`.

### Tests for Phase 23 (write first; they must fail)

`QuestionLifecycle.svelte` currently edits answer + status only and already sends `questionText: question.questionText` on update. `AnnotationCard` shows `<p class="text">{question.questionText}</p>` and renders `<QuestionLifecycle {question} />` **without** `onSave`. `NoteReader` holds `openQuestion` separately. `QuestionContext` already passes `{onSave}`. The question page `saved(value)` already replaces `question` (the `<h1>` uses `question.questionText`). `QuestionLifecycle` already calls `invalidate` for `questions` and `question:{id}`.

- [ ] T186 [P] [US2] Extend `web/tests/component/QuestionLifecycle.test.ts` only: accessible Question text editor; save sends a trimmed non-empty `questionText` with the other fields; blank/failed saves keep the draft; Cancel performs no mutation.
- [ ] T186a [P] [US2] Extend `web/tests/component/AnnotationCard.test.ts` only: after a successful question-text save, the card shows the returned canonical text (via `onSave`) with no second editor on the card.
- [ ] T186b [P] [US2] Extend `web/tests/component/QuestionContext.test.ts` only: the context surface shows the updated text after `onSave`.
- [ ] T187 [P] [US2] Add `tests/e2e/dog-food-edit-question.spec.ts`: edit from the central/detail flow and from a linked-note highlight card, reload each surface, same question ID, no duplicate. Include cancel/blank validation. Keep `tests/e2e/accessibility-notes.spec.ts` passing.

### Implementation for Phase 23

- [ ] T188 [US2] Add the question-text editor to `web/src/lib/components/QuestionLifecycle.svelte` only.
  **Read:** that file, `web/src/lib/api/questions.ts`, `web/src/lib/types/question.ts`.
  **Do:** Local `questionText` draft from `question.questionText`. Labelled input/textarea with Save and Cancel. Reject whitespace-only before the API call. Send `questionText` together with the current answer, status, priority, dueDate, tagIds, and `version` (same payload shape as today's `save()`). Update local state only after success. Keep failed drafts and show typed errors. Match `export let`.
  **Do not:** mutate the note body or `questionLinks`; do not convert an annotation into a question; do not add a second editor on `AnnotationCard`.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts`

- [ ] T189 [US2] Pass `onSave` through the highlight card.
  **Read:** `web/src/lib/components/AnnotationCard.svelte` only.
  **Do:** `export let onSave: ((question: Question) => void) | undefined = undefined`. Pass it into `QuestionLifecycle`. On save, replace the card's `question` with the returned record so `<p class="text">` updates. Do not add another question-text textarea on the card.
  **Verify:** `npm --prefix web test -- --run tests/component/AnnotationCard.test.ts tests/component/QuestionLifecycle.test.ts`

- [ ] T189a [US2] Keep `NoteReader`'s open card in sync.
  **Read:** `web/src/lib/editor/NoteReader.svelte` (`openQuestion` / `<AnnotationCard ... />`).
  **Do:** Pass `onSave` that sets `openQuestion = value`. Do not change capture.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/AnnotationCard.test.ts`

- [ ] T189b [US2] Confirm the central page already refreshes.
  **Read:** `web/src/lib/components/QuestionContext.svelte` and `web/src/routes/questions/[questionId]/+page.svelte` (`saved(value) { question = value }` and `<h1>{question.questionText}</h1>`).
  **Do:** If heading/context already update via existing `onSave`, do not add a second fetch. Only edit if the heading stays stale after T188. Do not change `QuestionList.svelte` (it renders `question.questionText` from props; invalidate already exists).
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts tests/component/AnnotationCard.test.ts tests/component/QuestionContext.test.ts` then `npx playwright test tests/e2e/dog-food-edit-question.spec.ts`

**Landmines:** `questionText` is canonical and must be sent with the current optimistic `version`; preserve all other editable fields when composing the update; do not log question or answer content; question text edits do not rewrite Markdown directives.

**Checkpoint**: A user can edit one canonical question from the central question flow or a linked-note card, and the new text is visible everywhere after reload.

---

## Phase 24: Optional “Why it matters if unanswered” context (dog-food gap remediation)

**Purpose**: Add the approved minimal optional context behind the dog-food example: a short explanation of why a question matters if it remains unanswered. This is informational context, not a second lifecycle or task-management system.

**Status**: The minimal scope is approved. T190 records the decision; T191–T195 remain implementation work. Delete-note behavior is already covered by Phase 22, and question-text editing is covered by Phase 23.

**Approved minimal design**: Add one optional, user-authored **plain-text** value named `consequenceText` (stored as `consequence_text`) to `kind=question`, labelled **Why it matters if unanswered**. It describes the consequence or risk of leaving the question unanswered. Show it in the question highlight card and question detail/answer-in-context controls, hide the field when empty, and allow it to be edited or cleared. Keep it informational only: it must not affect status, answer validation, priority, due date, reminders, or the Next queue. It is searchable and included in export/import. It does not apply to annotations, the note-question rail, or central list rows; deletion previews need no separate consequence presentation because the value follows the question if it is retained and is deleted with the question if the user chooses deletion.

**Independent Test**: Create a question without `consequenceText` and verify the normal workflow is unchanged. Add the value from the question detail or highlight card, view it in both approved surfaces, edit and clear it, and verify the canonical value survives reload, search, and export/import without becoming required for answering. Verify an annotation cannot receive the field and no status, priority, scheduling, reminder, or Next behavior changes.

### Tests and implementation for Phase 24

**Decision already recorded:** T190 is in [`completed-task.md`](./completed-task.md). T191–T195 remain implementation work. Put the field on `QuestionLifecycle` so the highlight card and the context page both get it (both already render that component). Annotations use the `AnnotationCard` textarea, not `QuestionLifecycle`.

JSON decoder rejects unknown fields — OpenAPI + Go struct before any request body includes `consequenceText`. `questions` INSERT/SELECT/UPDATE column lists live in `server/internal/store/sqlite/capture.go`. `UpdateLifecycle` copies a `QuestionWrite` literal and will drop a new column if omitted. FTS is `questions_fts(question_id, workspace_id, question_text, answer_markdown)` with triggers in `001_initial.sql` — add a **new** migration; do not edit `001_initial.sql`. Import allow-list is `server/internal/importexport/apply.go` `questions` map.

#### Tests (write first; they must fail)

- [ ] T191 [P] [US3] Domain tests only (`server/internal/domain/question_test.go` or a new `consequence_test.go`): null/empty `ConsequenceText` is valid; whitespace-only normalizes to nil; annotations cannot receive a non-empty value; answering still does not require it.
- [ ] T191a [P] [US3] SQLite round-trip tests in `server/internal/store/sqlite/capture_test.go` (or sibling): insert/update/get nullable `consequence_text`; optimistic `version` still increments; missing column on old rows reads as nil.
- [ ] T191b [P] [US3] API contract tests (`server/internal/api/handlers/capture_contract_test.go` or sibling): request/response round-trip `consequenceText`; absent field remains valid; unknown extra fields still 400.
- [ ] T191c [P] [US3] Search tests (`server/internal/store/sqlite/search_test.go`): a question is found by consequence text without duplicating results; annotations are unchanged.
- [ ] T191d [P] [US3] Import/export tests (`server/internal/importexport` / sqlite importexport tests): `consequenceText` round-trips; omitted field stays valid.
- [ ] T191e [P] [US3] Component tests in `web/tests/component/QuestionLifecycle.test.ts` (and `AnnotationCard.test.ts` only if the card needs extra chrome): display, edit, clear, retain draft; hide when empty; annotations show no such field.

#### Contracts then code

- [ ] T192 [US3] OpenAPI only: add nullable `consequenceText` on question request/response in `specs/001-track-note-questions/contracts/openapi.yaml` and copy to `tests/fixtures/contracts/openapi.yaml`. Existing records/requests stay valid when the value is absent.
  **Do not:** edit SQLite or Go yet.
- [ ] T192a [US3] Add `consequence_text` to `specs/001-track-note-questions/data-model.md` (nullable, question-only, informational, no lifecycle side effects).
- [ ] T192b [US3] Document the field, label **Why it matters if unanswered**, and no-lifecycle-side-effects rule in `spec.md` and `plan.md`.

- [ ] T193 [US3] Domain field in `server/internal/domain/question.go`.
  **Do:** Nullable `ConsequenceText *string` `json:"consequenceText,omitempty"` on `Question` and `QuestionWrite`. `Validate()`: trim; whitespace-only → nil; if `Kind == KindAnnotation` and non-nil, error; do not require it for `StatusAnswered`.
  **Verify:** `go -C server test ./internal/domain -count=1`

- [ ] T193a [US3] Migration `server/internal/store/migrations/003_question_consequence.sql` (name may increment if 003 exists): `ALTER TABLE questions ADD COLUMN consequence_text TEXT;` plus any CHECK that annotations stay null if cheap. Register in `embed.go` the same way `002_question_kind.sql` is registered. Do not edit `001_initial.sql`.
  **Verify:** `go -C server test ./internal/store/sqlite -count=1`

- [ ] T193b [US3] Persist the column in `server/internal/store/sqlite/capture.go`.
  **Read:** `CreateQuestion` INSERT, `getQuestionBase` / `GetQuestion` SELECT, `UpdateQuestion` UPDATE, note-link summary SELECT, and `UpdateLifecycle` in `lifecycle.go`.
  **Do:** Add `consequence_text` to every questions column list. Copy `ConsequenceText` through `UpdateLifecycle`'s `QuestionWrite` literal so lifecycle saves do not wipe it. Preserve workspace checks and `version`.
  **Verify:** `go -C server test ./internal/store/sqlite -count=1`

- [ ] T193c [US3] API mapping.
  **Read:** `server/internal/api/handlers/capture.go` `updateQuestion` / create (they decode `domain.QuestionWrite`).
  **Do:** If the struct json tag is enough, do not add a route. Confirm create/update/get responses include the field when set. Keep unknown JSON fields 400.
  **Verify:** `go -C server test ./internal/api/handlers -count=1`

- [ ] T193d [US3] FTS index in a new migration (not `001_initial.sql`).
  **Read:** `questions_fts` and its insert/update/delete triggers in `001_initial.sql` as a template.
  **Do:** Rebuild or extend `questions_fts` so `consequence_text` is searchable. Update `server/internal/store/sqlite/search.go` snippet column indexes if the FTS column list changes (`snippet(questions_fts,2,...)` today).
  **Verify:** `go -C server test ./internal/store/sqlite -count=1`

- [ ] T193e [US3] Import allow-list in `server/internal/importexport/apply.go`.
  **Do:** Add `"consequenceText": true` to the `questions` allowed map. Export is `SELECT *` — do not special-case it. Deletion keeps/deletes the question as today; no extra preview chrome.
  **Verify:** `go -C server test ./internal/importexport ./internal/store/sqlite -count=1`

- [ ] T194 [US3] UI on `web/src/lib/components/QuestionLifecycle.svelte` (covers card + context).
  **Read:** that file and `web/src/lib/types/question.ts`.
  **Do:** Add optional `consequenceText` to the frontend `Question` type and include it on `questionsApi.update` payloads. Label **Why it matters if unanswered**. Edit and clear. Hide the display when empty but still allow adding. Whitespace-only → omit (null). Failed-save keeps the draft. Do not change status/answer/priority/dueDate/reminder/Next logic.
  **Do not:** add the field to `NoteQuestionRail`, `QuestionList`, or the annotation branch of `AnnotationCard`.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts tests/component/AnnotationCard.test.ts tests/component/QuestionContext.test.ts`

- [ ] T195 [US3] Add `tests/e2e/dog-food-question-consequence.spec.ts` covering absent, add, edit, clear, reload, search, and export/import. Annotations and questions without the value remain fully supported.
- [ ] T195a [US3] Update `traceability.md` / `validation-results.md` only for tests actually run in this phase.

**Landmines:** Do not make answering depend on `consequenceText`; do not infer it from an answer or generate it automatically; do not apply it to annotations; do not add the OpenAPI or SQLite field before the contract/model task; treat whitespace-only input as absent; never log question or answer content.

**Checkpoint**: An optional plain-text consequence is available on question cards and detail pages, survives the approved persistence/search/portability paths, and has no effect on required question lifecycle behavior.

---

## Dependencies & Execution Order

**Completed phases 1–12** (and T139–T142, T190) are archived in [`completed-task.md`](./completed-task.md). The remaining active IDs start at T143.


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
21. **Phase 20a — Floating capture toolbar & mobile highlight positioning** addresses UX feedback on long notes and touch devices: floats the capture toolbar adjacent to highlighted text, anchors the composer/card in-context, and preserves responsive mobile touch access. Depends on Phase 20.
22. **Phase 21 — Map/trace** depends on the phases you actually shipped.
23. **Phase 21a — Capture composer focus hotfix** depends on the Phase 20a floating interaction code and should be completed before Phase 21b and Phase 22. It is client-only and has no backend, schema, or contract dependency.
24. **Phase 21b — Note-rail card-opening dismissal hotfix** depends on the Phase 14 note-question rail / `focusQuestionId` behavior and the Phase 21a window dismissal guards. It is client-only and must preserve both rail navigation and genuine outside dismissal.
25. **Phase 22 — Wire Delete note** depends on US1 note get/update and the existing note page. Use US6 preview/execute if those handlers already exist (Phase 8). Do not rebuild deletion backend unless the client has no delete operation to call. Independent of Phases 13–21b.
26. **Phase 23 — Edit canonical question text** depends on the existing US2 question update contract and the Phase 11 highlight card / Phase 13 context surfaces. It is independent of the delete-note wiring and does not require a schema or route change.
27. **Phase 24 — Optional consequence context** has its minimal scope approved in T190. Implement T191–T195 after the question lifecycle surfaces are available, updating contracts before the schema/API work and preserving the optional, question-only, informational behavior.

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
                                                         → Phase 20a floating capture & mobile highlight positioning
                                                         → Phase 21a capture composer focus hotfix
                                                         → Phase 22 wire Delete note button (US6)
                                                         → Phase 23 edit canonical question text
                                                         → Phase 24 optional consequence context (minimal scope approved)
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

### Phases 13–24 (small-model slices)

```text
One agent = one task. Never start T14x while T13x tests are unwritten.
Letter siblings are sequential (T143 then T143a).
13: T139–T142 done → T143 load+QuestionContext → T143a note switch → T144 list links
14: T145 rail test → T146 NoteReader focus test → T146a e2e → T147 rail → T148 focus prop → T149 mount rail → T149a focusId walk
15: T150 helper test → T151 NextQueue test → T151a e2e → T152 helper → T153 NextQueue.svelte → T153a /next page → T153b paginate → T154 nav
16: T155 domain test → T156 lifecycle test → T156a e2e → T157 domain+2 callers → T158 422 field → T159 reject → T159a resume-date input
17: T160–T162 tests (T161 markdown / T161a card) → T163 insert helper → T164 renderNoteHtml → T164a CSS → T164b pass questions → T165 card button → T165a NoteEditor insert → T165b NoteReader wire → T165c context page
18: T166/T166a/T166b unit → T167 composer error → T167a e2e → T168 exact → T168a whitespace → T168b markers → T168c wrapSelection → T169 match-before-create
19: T170 reader test → T170b unlinked-create test → T171 e2e → T172 picker UI → T172a wrap-without-create → T173 New question → T174 unlinked still listed
20: T175 keydown test → T175a e2e → T176 onWindowKeydown → T177 toolbar hint
20a: T177a–T177c tests (reader / card / e2e) → T177d selection coordinates & touch listener → T177e floating toolbar/composer → T177f anchored card/sheet → T177g e2e regression
21: T178 map → T179 FR-041–047 → T179a FR-048–055 → T179b SC-011–014 → T180 only tests you ran
21a: T180a–T180c focus-regression tests → T180d focus-safe event guards → T180e stable layer transition → T180f e2e/accessibility verification
21b: T180g–T180h rail-dismissal regression tests → T180i narrow rail-navigation event marker/guard → T180j e2e and Phase 21a regression verification
22: T181–T182 tests → T183 confirm notesApi → T184 existing page button (not NoteEditor) → T185 invalidate+goto
23: T186 lifecycle test → T186a card → T186b context → T187 e2e → T188 editor → T189 card onSave → T189a NoteReader → T189b page heading
24: T190 done → T191–T191e tests by layer → T192 OpenAPI → T192a data-model → T192b spec/plan → T193 domain → T193a migration → T193b store → T193c API → T193d FTS → T193e import → T194 QuestionLifecycle UI → T195 e2e → T195a evidence
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
10. **Phase 21a** → restore focus in the floating question/annotation capture composer before destructive or later dog-food work.
11. **Phase 21b** → keep right-side note-rail question navigation from being mistaken for an outside click while preserving explicit dismissal.
12. **Phase 22** → make the existing Delete note button actually delete (cancel-safe).
13. **Phase 23** → make canonical question text explicitly editable from central and linked-note surfaces.
14. **Phase 24** → implement the approved optional plain-text consequence/why-it-matters context without changing question lifecycle behavior.

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
- Phases 13–21b are further gap remediation (US8–US12). A new session reads `PROJECT_MAP.md`, **only the current phase in this file**, and only named files. One task per session for small models. Letter-suffixed IDs (`T143a`) are the next slice of the parent task — do not skip ahead. Do not start a later phase to “also add” extra product ideas (AI, flashcards, saved filters, Feynman modes).
- Phase 20a is UX gap remediation for floating/in-context capture and mobile highlight positioning. Only touches `NoteReader.svelte`, `AnnotationCard.svelte`, and their tests.
- Phase 21a is emergency gap remediation for the Phase 20a capture-composer focus regression. A new session reads `PROJECT_MAP.md`, **only Phase 21a**, and only named files. Keep native textbox focus, floating anchors, keyboard shortcuts, and explicit dismissal intact; do not fix it by disabling outside-click handling globally.
- Phase 21b is emergency gap remediation for the Phase 21a note-rail card-opening dismissal regression. A new session reads `PROJECT_MAP.md`, **only Phase 21b**, and only named files. Mark only intentional rail navigation as internal; preserve genuine outside-click dismissal and do not move the rail or add a second card.
- Phase 22 is gap remediation for Delete note (US6). The control lives on `web/src/routes/notes/[noteId]/+page.svelte`, not `NoteEditor`. A new session reads `PROJECT_MAP.md`, **only Phase 22**, and only named files. Do not add a second delete button.
- Phase 23 is dog-food gap remediation for canonical question-text editing. A new session reads `PROJECT_MAP.md`, **only Phase 23**, and only named files; do not redo the existing question update API. Put the editor in `QuestionLifecycle`; pass `onSave` through `AnnotationCard` / `NoteReader`.
- Phase 24 has an approved minimal scope: optional plain-text `consequenceText` for questions only. Follow T191–T191e tests, T192–T192b contracts, T193–T193e backend slices, then T194 UI on `QuestionLifecycle`. Do not broaden the field or give it lifecycle side effects.
