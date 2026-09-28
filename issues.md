# Issues

## E2E: `us1-notes-highlight.spec.ts` has a strict locator duplicate

- **Status:** Open
- **Test:** `tests/e2e/us1-notes-highlight.spec.ts`
- **Browser:** Chromium
- **Observed failure:** `getByText('The mitochondria is the powerhouse of the cell.').toBeVisible()` fails in strict mode because the locator resolves to two elements:
  1. The note paragraph containing the sentence.
  2. A highlight-card blockquote containing the exact sentence.
- **Suggested follow-up:** Scope the assertion to the intended reader/highlight element or use a more specific locator so the test does not require the sentence to be unique on the page.

## E2E browser executables unavailable

- **Status:** Environment/setup
- **Tests:** `tests/e2e/us12-passage-match.spec.ts`, `tests/e2e/us1-notes-highlight.spec.ts`
- **Browsers:** Firefox and WebKit
- **Observed failure:** Playwright could not find the installed Firefox and WebKit browser executables.
- **Suggested follow-up:** Run `npx playwright install` before running the full multi-browser E2E suite.
