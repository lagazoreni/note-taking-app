<!--
Sync Impact Report
- Version change: template (unversioned) -> 1.0.0
- Modified principles:
	- Template Principle 1 -> I. Lightweight by Default
	- Template Principle 2 -> II. Notes Stay Fast and Focused
	- Template Principle 3 -> III. Explicit SvelteKit-Go Boundary
	- Template Principle 4 -> IV. Container-Ready Delivery
	- Template Principle 5 -> V. Tested, Observable, and Secure
- Added sections: Technology and Architecture Constraints; Development and Release Workflow
- Removed sections: none
- Follow-up TODOs: none
-->
# Note Taking App Constitution

## Core Principles

### I. Lightweight by Default
Every dependency, service, abstraction, and runtime process MUST have a demonstrated need for the
current product. Implementations MUST prefer standard-library capabilities and existing framework
features over custom infrastructure. A change that adds operational or architectural complexity
MUST document the user benefit and why a simpler design is insufficient. This keeps the application
small enough to understand, operate, and ship without a dedicated platform team.

### II. Notes Stay Fast and Focused
Creating, opening, editing, saving, searching, and deleting notes are the core workflows. Each MUST
remain usable on narrow and wide viewports, provide clear loading and failure states, and avoid
unnecessary navigation. User input MUST be protected from silent loss through explicit persistence
behavior and recoverable error handling. New capabilities MUST NOT degrade these workflows without
measured evidence and an approved tradeoff.

### III. Explicit SvelteKit-Go Boundary
SvelteKit MUST own browser rendering, client interaction, and frontend routing. Go MUST own business
rules, persistence, and backend APIs. Communication MUST use a documented, versioned HTTP contract
with stable request, response, validation, and error shapes. Business rules MUST NOT be duplicated
in the frontend; frontend validation MAY improve feedback but the backend remains authoritative.
Contract changes MUST update and test both consumers and providers in the same change.

### IV. Container-Ready Delivery
The application MUST build and run as reproducible containers from committed definitions. Runtime
configuration and secrets MUST enter through the environment or mounted files and MUST NOT be baked
into images. Containers MUST run as non-root users where supported, expose health checks, write logs
to standard output and error, and persist user data only through declared volumes or external
services. Production images MUST contain only runtime requirements and use pinned, reviewable base
image versions.

### V. Tested, Observable, and Secure
Changes MUST include automated tests at the lowest useful level. Core note workflows and API
contracts MUST have integration coverage; release candidates MUST pass frontend, backend, and
container smoke checks. The backend MUST emit structured logs with request context while excluding
note content, credentials, and secrets. All external input MUST be validated, errors MUST avoid
leaking internals, and dependencies MUST be reviewed for known vulnerabilities before release.

## Technology and Architecture Constraints

- The web frontend MUST use SvelteKit and MUST follow its established routing, rendering, and data
	loading conventions.
- The backend MUST use Go and MUST keep transport, business logic, and persistence concerns
	independently testable.
- The frontend and backend MUST be independently buildable. Shared behavior MUST be expressed
	through the API contract rather than cross-runtime source coupling.
- Persistence MUST use the simplest durable store that meets documented requirements. Introducing
	caches, queues, or additional services requires measured need and an operational plan.
- Local and deployed environments MUST use the same configuration keys and container entrypoints.
	Environment-specific values MUST remain outside source control.

## Development and Release Workflow

1. Each change MUST state its user-visible outcome and identify affected frontend, API, persistence,
	 and deployment contracts.
2. Implementations MUST remain scoped to the stated outcome; unrelated abstractions and services
	 require separate justification and review.
3. Before merge, formatting, static analysis, and automated tests for affected SvelteKit and Go code
	 MUST pass. API changes MUST also pass contract tests.
4. Before release, container images MUST build from a clean checkout and pass startup, health, and
	 core note workflow smoke checks using production-like configuration.
5. Reviews MUST verify compliance with every applicable principle. Any approved exception MUST be
	 recorded with its owner, rationale, risk, and removal or review date.

## Governance

This constitution governs all specifications, plans, implementation tasks, and reviews for the
project. Where another project document conflicts with it, this constitution takes precedence.

Amendments MUST be proposed as a documented change that states the motivation, affected principles,
migration impact, and requested semantic version bump. Approval requires review by the project
maintainers. A change that removes or incompatibly redefines governance requires a MAJOR bump; a new
principle or material expansion requires a MINOR bump; a clarification with no governance impact
requires a PATCH bump.

Every feature plan and pull request MUST include a constitution compliance check. Maintainers MUST
review the constitution before each release and at least once every twelve months. Violations MUST
be corrected before release or documented as time-bounded exceptions under the workflow above.

**Version**: 1.0.0 | **Ratified**: 2026-08-09 | **Last Amended**: 2026-08-09
