# Browser tests

Playwright runs every journey in Chromium, Firefox, and WebKit. Tests block external
requests while allowing same-origin localhost traffic so “offline” means no internet,
not an unavailable local Go service. Install browsers with `npx playwright install`.

The local SvelteKit server is started by `playwright.config.ts`; API-backed journeys
start the Go service through the test harness or use the configured development proxy.
