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

| SC-001, SC-002, SC-010 | `tests/usability/protocol.md`; results are pending participant evidence |
| SC-003–SC-009 | optimistic versions, indexed query/FTS, atomic deletion/import tests, validation results |

Pending browser/container/usability checks are explicitly recorded in
`validation-results.md`; no exception waives a requirement.
