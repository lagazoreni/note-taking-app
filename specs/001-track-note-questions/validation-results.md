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
| Playwright matrix | PENDING | Browser binaries/service-backed harness require release environment |
| Container smoke | PENDING | Requires OCI image build/runtime |
| Usability success criteria | PENDING | Protocol and empty result record committed |

The Playwright attempt was blocked by the execution environment's permission to bind localhost ports, and the container attempt was blocked because the Docker daemon was unavailable. These pending items are release-environment evidence, not silently treated as passes.
