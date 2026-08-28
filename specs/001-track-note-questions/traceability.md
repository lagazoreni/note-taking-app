# Requirements traceability

| Range | Evidence |
|---|---|
| FR-001–FR-011 | Capture repository/service, note editor, question list, US1/US2 component and browser tests |
| FR-012–FR-019 | `domain/question_lifecycle.go`, lifecycle store/service, QuestionLifecycle UI, US3 tests |
| FR-020–FR-027 | question query/filter, FTS search, organization store/API/UI, US4/US5 tests |
| FR-028–FR-031 | deletion impact service, one-transaction delete, review dialog, US6 tests |
| FR-032–FR-040 | service worker shell cache, same-origin API, import/export validation/apply, no-internet and portability tests |
| SC-001, SC-002, SC-010 | `tests/usability/protocol.md`; results are pending participant evidence |
| SC-003–SC-009 | optimistic versions, indexed query/FTS, atomic deletion/import tests, validation results |

Pending browser/container/usability checks are explicitly recorded in
`validation-results.md`; no exception waives a requirement.
