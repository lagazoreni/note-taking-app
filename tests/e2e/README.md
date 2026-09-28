# Browser tests

Current Playwright verification runs in Chromium only. Firefox and WebKit are intentionally
skipped for now; use `--project=chromium` when running E2E tests until those browsers are
re-enabled. Tests block external requests while allowing same-origin localhost traffic so
“offline” means no internet, not an unavailable local Go service. Install browsers with
`npx playwright install` when cross-browser coverage is resumed.

The local SvelteKit server is started by `playwright.config.ts`; API-backed journeys
start the Go service through the test harness or use the configured development proxy.
