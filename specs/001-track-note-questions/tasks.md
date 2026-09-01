# Tasks: Interactive Note Questions

**Input**: Design documents from `/specs/001-track-note-questions/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Completed work**: Phases 1–12 (T001–T138), Phase 13 T139–T142, and Phase 24 T190 are archived in [`completed-task.md`](./completed-task.md). Do not re-implement them.

**New-session start**: Read [`PROJECT_MAP.md`](../../PROJECT_MAP.md) at the repository root before opening other files. Then read **only the current phase in this file**, then only the files that phase names. For leftover Phase 13 wiring use **T143–T144**. For highlight/annotation follow-on work, Next queue, deferred dates, resolved highlights, sturdier wrapping, reading-view link/attach, and keyboard capture, use **Phases 14–20**. For the no-op Delete note button, use **Phase 22**; for canonical question-text editing, use **Phase 23**; for the approved optional consequence/why-it-matters behavior, use **Phase 24**. Do not scan the whole repo unless the map is stale.

**Tests**: Automated tests are included because the project constitution requires unit, integration, API contract, browser workflow, and container smoke coverage. Within each story, create the listed tests first and verify they fail for the expected missing behavior before implementation.

**Organization**: Remaining tasks are grouped by gap-remediation phase. Every phase ends with an independently executable browser scenario and checkpoint.

**Workspace prerequisite**: A fresh database can contain zero workspaces. Workspace bootstrap already shipped in Phases 1–3 (see `completed-task.md`). Remaining work must not regress that guard.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can be implemented in parallel with adjacent tasks because it changes different files and has no dependency on an unfinished adjacent task.
- **[Story]**: Maps the task to a user story in `spec.md`.
- Paths are relative to the repository root.
- `a`-suffixed IDs are gap-remediation tasks inserted without renumbering later tasks.

---

## Phase 13: Answer in context (US8) — gap remediation

**Purpose**: Opening a question shows the source note and highlighted passage on the same page as answer/status controls.

**For a small model.** Do not scan the repo. Do not reread US1–US7 or Phase 11 internals. Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement **one task**, run its verify command, stop.

**Independent Test**: Capture a question on a sentence. From Active Questions open it. The question page shows the note, the sentence is highlighted and in view, and the answer box is on the same page. An unlinked question still opens and says it is unlinked.

**Do not**: change capture, kind, lifecycle rules, NoteReader toolbar, `/notes/new` labels, or add APIs.

**Already done (do not rewrite):** T139–T142 — `NoteExcerpt.test.ts`, `QuestionContext.test.ts`, `us8-answer-in-context.spec.ts`, `NoteExcerpt.svelte`, `QuestionContext.svelte`. Details in [`completed-task.md`](./completed-task.md).

### Implementation remaining for Phase 13

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

## Phase 22: Wire Delete note (US6) — gap remediation

**Purpose**: The note page already shows a Delete note button, but clicking it does nothing. Make that existing control delete the current note (cancel-safe), then leave the note page.

**For a small model.** Do not scan the repo. Do not reread US1–US12 or Phases 11–21 internals. Read `PROJECT_MAP.md` once, this phase only, then only the files named in the current task. Implement **one task**, run its verify command, stop.

**Independent Test**: Create a workspace and a note. Open the note. Click Delete note, then cancel — the note is unchanged and still listed under Notes. Click Delete note again and confirm — the app leaves the editor, `/notes` no longer lists that note, and opening its old URL does not show the editor. A note that has questions still deletes; the user is not stuck on a dead button.

**Do not**: add a second Delete note button; change capture, highlights, `/notes/new` labels, Next queue, or question lifecycle rules; rebuild the Phase 8 deletion backend unless the client has no delete operation to call.

### Tests for Phase 22 (write first; they must fail)

- [ ] T181 [P] [US6] Extend `web/tests/component/NoteEditor.test.ts` (and `web/tests/component/DeleteNoteReview.test.ts` if the review dialog is used): an existing saved note renders a control named Delete note; clicking it opens confirmation or `DeleteNoteReview`; confirming calls delete / `onDelete` once; cancel does not call delete and leaves the note loaded. `/notes/new` must not expose a working delete for an unsaved note.
- [ ] T182 [P] [US6] Add Playwright `tests/e2e/us6-delete-note-button.spec.ts`: create a workspace and a uniquely titled note; open it; Delete note → cancel → still on the note page and the title still appears on `/notes`; Delete note → confirm → land on `/notes` without that title. Keep `tests/e2e/accessibility-notes.spec.ts` passing.

### Implementation for Phase 22

- [ ] T183 [P] [US6] Ensure the frontend can call note deletion in `web/src/lib/api/notes.ts`.
  **Read:** `web/src/lib/api/notes.ts`, `web/src/lib/api/client.ts`, and the existing note-deletion operations in `specs/001-track-note-questions/contracts/openapi.yaml` (Phase 8 / T096).
  **Do:** Add or finish `notesApi` methods that match those operations (preview + execute if that is the contract; otherwise `DELETE` with current `version`). Do not invent a new path. Typed errors must surface 404, 409 version conflict, and 422 the same way other note writes do.
  **Do not:** add a new backend route; do not log note bodies.
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts` (client compile is enough if the new tests still fail on missing UI).

- [ ] T184 [US6] Wire the existing Delete note button.
  **Read:** `web/src/lib/editor/NoteEditor.svelte` and `web/src/routes/notes/[noteId]/+page.svelte` only far enough to find the already-rendered Delete note control; `web/src/lib/components/DeleteNoteReview.svelte` if it exists.
  **Do:** Attach the handler to that existing control — do not add another button. Show `DeleteNoteReview` when that component exists; otherwise a dialog with Cancel and Delete note. Confirm runs the T183 client with the loaded note `id` and `version`. Cancel closes the dialog with no API call. Hide or disable delete on `/notes/new` until a note has been saved.
  **Do not:** delete on the first click with no confirm; do not navigate yet (T185).
  **Verify:** `npm --prefix web test -- --run tests/component/NoteEditor.test.ts tests/component/DeleteNoteReview.test.ts`

- [ ] T185 [US6] After a successful delete, leave the note page.
  **Read:** `web/src/routes/notes/[noteId]/+page.svelte`, `web/src/lib/stores/query.ts`, `web/src/lib/stores/toast.ts`.
  **Do:** On success: toast that the note was deleted, invalidate note/question list queries, `goto('/notes')`. On 409: show the conflict and do not navigate. On other errors: show the API envelope and stay on the page. If preview/execute is required, call preview then execute with the review decisions; do not invent a new deletion policy.
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

- [ ] T186 [P] [US2] Extend `web/tests/component/QuestionLifecycle.test.ts`, `web/tests/component/AnnotationCard.test.ts`, and `web/tests/component/QuestionContext.test.ts`: a question exposes an accessible Question text editor, saves a trimmed non-empty value through the supplied save path, preserves the other question fields, reports blank/failed saves without losing the draft, and Cancel performs no mutation. The question highlight card must receive the updated canonical question through its save callback.
- [ ] T187 [P] [US2] Add `tests/e2e/dog-food-edit-question.spec.ts`: edit a question from its central/detail flow and from a linked-note highlight card, reload each surface, and verify the same question ID has the new text with no duplicate. Include cancel/blank validation and preserve `tests/e2e/accessibility-notes.spec.ts`.

### Implementation for Phase 23

- [ ] T188 [US2] Add the question-text editing control to `web/src/lib/components/QuestionLifecycle.svelte`.
  **Read:** `web/src/lib/components/QuestionLifecycle.svelte`, `web/src/lib/api/questions.ts`, and `web/src/lib/types/question.ts`.
  **Do:** Keep a local `questionText` draft initialized from `question.questionText`; render a labelled, keyboard-accessible input or textarea with Save and Cancel behavior; reject whitespace-only text before the API call; send the edited `questionText` together with the current answer, status, priority, due date, tags, and optimistic `version`; update local state only after success; retain failed drafts and surface typed conflict/validation errors.
  **Do not:** make text editing mutate the note body or link rows, or allow an annotation to be converted into a question.
  **Verify:** `npm --prefix web test -- --run tests/component/QuestionLifecycle.test.ts`

- [ ] T189 [US2] Propagate canonical question-text saves through all existing question surfaces in `web/src/lib/components/AnnotationCard.svelte`, `web/src/lib/editor/NoteReader.svelte`, `web/src/lib/components/QuestionContext.svelte`, `web/src/routes/questions/[questionId]/+page.svelte`, and the relevant question-list/link components.
  **Do:** Pass `onSave` into the highlight card, replace the card/context/list's local question with the returned record, invalidate the existing question queries, and ensure the central page heading, linked-note card, and any inline question presentation refresh without a duplicate fetch or stale text. Keep the existing answer/status save behavior and `/notes/new` Title/Note labels unchanged.
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

**Decision already recorded:** T190 is in [`completed-task.md`](./completed-task.md). T191–T195 remain implementation work.

- [ ] T191 [P] [US3] Add failing domain, repository, API-contract, search, import/export, and component tests for optional `consequenceText`: null/empty values preserve the existing workflow; non-empty text round-trips; whitespace-only input is treated as absent; annotation records cannot receive it; question cards and detail controls display, edit, clear, and retain it; search can find it without duplicating results.
- [ ] T192 [US3] Update the approved product contracts and model in `spec.md`, `plan.md`, `data-model.md`, `contracts/openapi.yaml`, and `tests/fixtures/contracts/openapi.yaml`: define nullable `consequenceText`/`consequence_text`, the exact label and question-only scope, and the informational/no-lifecycle-side-effects rule. Keep existing records and requests valid when the value is absent.
- [ ] T193 [US3] Implement nullable `consequence_text` persistence, normalization, question-only validation, API request/response mapping, FTS/search indexing, deterministic export/import, and any required allow-list or deletion-retention behavior across the authoritative Go domain/store/services/handlers. Preserve workspace boundaries and optimistic versions.
- [ ] T194 [US3] Implement the optional field in the question highlight card and question detail/answer-in-context controls with accessible label **Why it matters if unanswered**, edit/clear behavior, empty-state hiding, failed-save draft preservation, and no automatic changes to lifecycle, scheduling, reminders, priority, or Next. Do not add it to annotation-only UI, the note-question rail, or central list rows.
- [ ] T195 [US3] Add `tests/e2e/dog-food-question-consequence.spec.ts` covering absent, add, edit, clear, reload, search, and export/import behavior; update traceability/validation evidence only for tests actually run and verify annotations and questions without the value remain fully supported.

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
21. **Phase 21 — Map/trace** depends on the phases you actually shipped.
22. **Phase 22 — Wire Delete note** depends on US1 note get/update and the existing note page. Use US6 preview/execute if those handlers already exist (Phase 8). Do not rebuild deletion backend unless the client has no delete operation to call. Independent of Phases 13–20.
23. **Phase 23 — Edit canonical question text** depends on the existing US2 question update contract and the Phase 11 highlight card / Phase 13 context surfaces. It is independent of the delete-note wiring and does not require a schema or route change.
24. **Phase 24 — Optional consequence context** has its minimal scope approved in T190. Implement T191–T195 after the question lifecycle surfaces are available, updating contracts before the schema/API work and preserving the optional, question-only, informational behavior.

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
13: T139–T140 tests → T141 NoteExcerpt → T142 QuestionContext → T143 page → T144 list links
14: T145–T146 tests → T147 rail → T148 NoteReader focus → T149 editor wire
15: T150–T151 tests → T152 helper → T153 /next page → T154 nav
16: T155–T156 tests → T157 domain → T158 API mapping → T159 lifecycle UI
17: T160–T162 tests → T163 insert helper → T164 render/CSS → T165 card + save
18: T166–T167 tests → T168 findSelection → T169 wrap-before-create
19: T170–T171 tests → T172 link existing → T173 unlinked create → T174 attach = link
20: T175 tests → T176 keydown → T177 toolbar hint
21: T178 map → T179 trace → T180 only tests you ran
22: T181–T182 tests → T183 notesApi.delete → T184 wire existing button → T185 navigate/invalidate
23: T186–T187 tests → T188 question-text editor → T189 surface/state propagation
24: T190 approved decision → T191 tests → T192 contracts/model → T193 persistence/API/search/import → T194 UI → T195 browser/evidence
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
10. **Phase 22** → make the existing Delete note button actually delete (cancel-safe).
11. **Phase 23** → make canonical question text explicitly editable from central and linked-note surfaces.
12. **Phase 24** → implement the approved optional plain-text consequence/why-it-matters context without changing question lifecycle behavior.

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
- Phase 22 is gap remediation for the existing no-op Delete note button (US6). A new session reads `PROJECT_MAP.md`, **only Phase 22**, and only named files. Do not add a second delete button.
- Phase 23 is dog-food gap remediation for canonical question-text editing. A new session reads `PROJECT_MAP.md`, **only Phase 23**, and only named files; do not redo the existing question update API.
- Phase 24 has an approved minimal scope: optional plain-text `consequenceText` for questions only. T191–T195 must still update contracts and add tests before implementation; do not broaden the field or give it lifecycle side effects.
