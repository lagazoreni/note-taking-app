# Validation results

Run date: 2026-08-09

| Check | Result | Command / note |
|---|---|---|
| Go formatting | PASS | `gofmt -w server/internal server/cmd` |
| Go unit/integration tests | PASS | `cd server && go test ./...` |
| Svelte type check | PASS | `cd web && npm run check` |
| Web formatting/lint | PASS | `cd web && npm run format:check && npm run lint` |
| Vitest component/unit tests | PASS | `cd web && npm test -- --run` |
| Web production build | PASS | `cd web && npm run build` |
| Dependency audit | PASS (high threshold) | `cd web && npm audit --audit-level=high`; three low findings remain in SvelteKit's transitive cookie range and require upstream-compatible remediation |
| OpenAPI contract load | PASS | `server/internal/api/contracttest` |
| Playwright matrix | PENDING | Firefox/WebKit coverage still requires browser binaries; targeted Chromium evidence is recorded below |
| T180f targeted Playwright run (Chromium, 2026-09-16) | PASS — 8 passed | `npx playwright test tests/e2e/us12-capture-focus.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts --project=chromium --workers=1`; verified question/annotation typing, Q/Escape behavior, floating toolbar/card positioning, highlight capture, and note accessibility on desktop and mobile. |
| T180f configured-browser Playwright attempt (2026-09-16) | PARTIAL — 7 passed, 1 Chromium failure, 16 browser-launch failures | `npx playwright test tests/e2e/us12-capture-focus.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts`; Firefox and WebKit executables are unavailable; the parallel Chromium mobile floating-selection test missed the newly created highlight, then the serial Chromium rerun passed all 8 tests. |
| T180j targeted Playwright run (Chromium, 2026-09-17) | PASS — 11 passed | `npx playwright test tests/e2e/us8-note-rail.spec.ts tests/e2e/us12-capture-focus.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/accessibility-notes.spec.ts --project=chromium --workers=1`; verified note-rail card navigation and outside dismissal on desktop/mobile plus capture focus, floating positioning, keyboard shortcuts, highlight capture, and note accessibility. Firefox/WebKit executables are unavailable in this environment. |
| T180j component regression check (Chromium-independent, 2026-09-17) | PASS — 26 passed | `npm --prefix web test -- --run tests/component/NoteReader.test.ts tests/component/NoteQuestionRail.test.ts` |
| T180 targeted Playwright run (Chromium, 2026-09-14) | FAIL — 13 passed, 3 failed (16 tests) | `npx playwright test tests/e2e/us8-answer-in-context.spec.ts tests/e2e/us8-note-rail.spec.ts tests/e2e/us9-next-queue.spec.ts tests/e2e/us10-defer-date.spec.ts tests/e2e/us11-resolved-highlight.spec.ts tests/e2e/us12-passage-match.spec.ts tests/e2e/us12-link-attach.spec.ts tests/e2e/us12-keyboard-capture.spec.ts tests/e2e/us12-floating-selection-ux.spec.ts tests/e2e/us1-notes-highlight.spec.ts tests/e2e/us3-question-lifecycle.spec.ts tests/e2e/us2-shared-question.spec.ts tests/e2e/accessibility-notes.spec.ts --project=chromium --workers=1`; failures: `us12-link-attach.spec.ts` unlinked-create flow timed out waiting for the Question text textbox after `New question`, while `us2-shared-question.spec.ts` and `us3-question-lifecycle.spec.ts` hit strict-mode heading locators because the page also renders the matching empty-state heading. |
| Container smoke | PENDING | Requires OCI image build/runtime |
| Usability success criteria | PENDING | Protocol and empty result record committed |

The targeted T180 Playwright run was executed with Chromium. Firefox/WebKit coverage remains pending because those browser executables are unavailable in this environment. The container attempt was blocked because the Docker daemon was unavailable. These pending items are release-environment evidence, not silently treated as passes.

