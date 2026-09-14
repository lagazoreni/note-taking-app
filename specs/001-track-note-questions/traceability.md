# Requirements traceability

| Range | Evidence |
|---|---|
| FR-001–FR-011 | Capture repository/service, note editor, question list, US1/US2 component and browser tests |
| FR-012–FR-019 | `domain/question_lifecycle.go`, lifecycle store/service, QuestionLifecycle UI, US3 tests |
| FR-020–FR-027 | question query/filter, FTS search, organization store/API/UI, US4/US5 tests |
| FR-028–FR-031 | deletion impact service, one-transaction delete, review dialog, US6 tests |
| FR-032–FR-040 | service worker shell cache, same-origin API, import/export validation/apply, no-internet and portability tests |

## FR-041–FR-047 (Phases 13–15)

The entries below identify the assertions present in the named tests; they do not claim that a test has been run. Gaps in an assertion are called out rather than inferred as passing evidence.

| Requirement | Test evidence |
|---|---|
| FR-041 | `web/tests/component/QuestionContext.test.ts` checks that a linked note title and the answer/status controls render together; `web/tests/component/NoteExcerpt.test.ts` checks the matching passage is rendered in a `mark`; `tests/e2e/us8-answer-in-context.spec.ts` opens a captured question and checks the source passage and Answer field are visible on the question page. These tests do not directly assert scroll position. |
| FR-042 | `web/tests/component/QuestionContext.test.ts` supplies two linked notes, checks both source-note options, and checks that selecting the second calls `onSelectNote` with its id. There is no listed end-to-end assertion for default-first-note loading or fetching the second note. |
| FR-043 | `web/tests/component/QuestionContext.test.ts` checks that an unlinked question shows the currently-unlinked message while retaining the answer/status controls. |
| FR-044 | `web/tests/component/NoteQuestionRail.test.ts` checks directive-order rendering and excludes annotations and answered questions; `tests/e2e/us8-note-rail.spec.ts` checks the same exclusions and order in the reading view. |
| FR-045 | `web/tests/component/NoteQuestionRail.test.ts` checks that Next unanswered invokes the callback; `web/tests/component/NoteReader.test.ts` checks that a focused question opens its matching highlight card; `tests/e2e/us8-note-rail.spec.ts` walks the first and second open questions and checks the empty state after both are answered. The end-to-end test does not directly assert scroll position. |
| FR-046 | `web/tests/component/NextQueue.test.ts` checks the empty guidance and section-heading presentation; `tests/e2e/us9-next-queue.spec.ts` exercises the workspace's `/next` route and opens a queue item. No listed test asserts cross-workspace isolation or the absence of saved-filter controls directly. |
| FR-047 | `web/tests/unit/nextQueue.test.ts` checks section order, overdue/today/in-progress/deferred-ready/high-priority grouping, deduplication, and omission of answered, annotation, future-deferred, and null-date deferred items; `web/tests/component/NextQueue.test.ts` checks non-empty section headings and empty-section guidance; `tests/e2e/us9-next-queue.spec.ts` checks an overdue item appears, an answered item is omitted, and a queue item opens the question page. |

## FR-048–FR-055 (Phases 16–20)

The entries below identify implementation pointers and assertions present in the named tests; they do not claim that a test has been run. Gaps in an assertion are called out rather than inferred as passing evidence.

| Requirement | Test evidence |
|---|---|
| FR-048 | `server/internal/domain/question_lifecycle.go` requires a trimmed due date for Deferred while retaining the existing answer rule; `server/internal/domain/question_lifecycle_test.go` covers nil/empty/date Deferred transitions and an Answered transition without a due date; `web/tests/component/QuestionLifecycle.test.ts` covers client rejection and the date-bearing update payload; `tests/e2e/us10-defer-date.spec.ts` covers the browser rejection and successful save path. `server/internal/api/handlers/capture.go` maps the resume-date error to a `dueDate` validation field. No listed test directly exercises loading a legacy Deferred row with a null date or the annotation exemption. |
| FR-049 | `web/src/lib/editor/markdown.ts` carries live question status into `data-status`, and `web/src/lib/editor/NoteReader.svelte` / `web/src/lib/components/NoteExcerpt.svelte` provide the answered-mark styling. `web/tests/unit/markdown-status.test.ts` checks answered status metadata without exposing directives; `tests/e2e/us11-resolved-highlight.spec.ts` reloads an answered note and checks `data-status="answered"`. No listed test compares computed active/resolved colors directly. |
| FR-050 | `web/src/lib/editor/directives.ts` implements blockquote insertion and idempotence; `web/tests/unit/insertAnswer.test.ts` checks trimming, multiline blockquotes, repeated insertion, unknown ids, and empty answers; `web/tests/component/AnnotationCard.test.ts` checks the answered-question button and callback; `tests/e2e/us11-resolved-highlight.spec.ts` checks that answering alone does not add a blockquote, that an explicit insert does, and that it remains after reload. The UI path is user-initiated; the repeated-click guarantee is covered at helper level rather than by a browser assertion. |
| FR-051 | `web/src/lib/editor/directives.ts` implements exact, whitespace-collapsed, and Markdown-marker-aware matching before wrapping. `web/tests/unit/directives.test.ts` checks first exact match, collapsed whitespace, emphasis markers, directive order, and the stable no-match error; `web/tests/component/NoteReader.test.ts` checks a rejected capture keeps the composer open; `tests/e2e/us12-passage-match.spec.ts` checks a rendered bold selection wraps the source Markdown. No listed browser test proves that an unmappable selection makes neither the create request nor the note mutation. |
| FR-052 | `web/src/lib/editor/NoteEditor.svelte` filters linkable workspace questions and rejects ids already present in the note; `web/tests/component/NoteReader.test.ts` checks the reading toolbar/picker callback, and `web/tests/component/QuestionPicker.test.ts` checks search, selection, and unlinked labeling. `tests/e2e/us12-link-attach.spec.ts` links one question from a first note onto a second note and checks both source notes on the question page. No listed negative test covers a cross-workspace candidate or a duplicate link on the same note. |
| FR-053 | `web/src/routes/questions/+page.svelte` creates an unlinked question without a note mutation, while `web/src/lib/editor/NoteEditor.svelte` wraps an existing question when it is later attached. `web/tests/component/ActiveQuestions.test.ts` checks the question-only create payload and `web/tests/component/QuestionPicker.test.ts` checks that unlinked questions remain selectable; the second scenario in `tests/e2e/us12-link-attach.spec.ts` creates an unlinked question, attaches it from reading, and opens the resulting highlight. The E2E flow does not explicitly assert the question id/count before and after attachment. |
| FR-054 | `web/src/lib/editor/NoteReader.svelte` guards shortcuts by selection, modifier keys, focused controls, and composer state. `web/tests/component/NoteReader.test.ts` checks Q question capture, A annotation capture, L picker opening, Escape dismissal, and ignored Q input while the composer or an input/textarea is active; `tests/e2e/us12-keyboard-capture.spec.ts` checks Q followed by Escape on a selected passage. The listed browser test does not independently exercise A or L. |
| FR-055 | `web/src/lib/components/QuestionList.svelte` and `web/src/lib/components/NextQueue.svelte` include the first linked note as `noteId` in their question-page links. `tests/e2e/us8-answer-in-context.spec.ts` checks an Active Questions item opens the question page with the source passage and Answer field together; `tests/e2e/us9-next-queue.spec.ts` checks a Next item opens `/questions/{id}` while allowing the query string. No listed assertion checks the exact `noteId` query value or the source passage after following a Next link. |

| SC-001, SC-002, SC-010 | `tests/usability/protocol.md`; results are pending participant evidence |
| SC-003–SC-009 | optimistic versions, indexed query/FTS, atomic deletion/import tests, validation results |

Pending browser/container/usability checks are explicitly recorded in
`validation-results.md`; no exception waives a requirement.
