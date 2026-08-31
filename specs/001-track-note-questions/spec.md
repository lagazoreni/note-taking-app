# Feature Specification: Interactive Note Questions

**Feature Branch**: `001-track-note-questions`

**Created**: 2026-08-09

**Status**: Draft

**Input**: User description: "Build a local-first note-taking application that lets users capture questions in notes, manage shared questions centrally, organize them across workspaces, and work offline."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Capture and Track Questions (Priority: P1)

While reading or taking notes, a user creates one or more notes and can attach a question or annotation to a selected passage. Notes appear in the workspace's Notes index, highlighted passages open cards, and questions remain available from the central active Questions view while annotations stay out of that inbox.

**Why this priority**: Capturing follow-up work or context in place and returning to it centrally is the product's primary value.

**Independent Test**: Create a workspace and two notes, verify both appear on the Notes page, select a sentence in one note, add a question and an annotation, verify each passage is highlighted and opens its card, and verify only the question appears in Active Questions.

**Acceptance Scenarios**:

1. **Given** notes in the current workspace, **When** the user opens the Notes page, **Then** every note appears as a navigable entry and an empty workspace has actionable guidance.
2. **Given** a selected passage in a note, **When** the user adds a question or annotation, **Then** the passage is highlighted without exposing directive tokens and the new record is saved with the note.
3. **Given** a highlighted passage, **When** the user clicks it, **Then** a card opens with the quoted passage and question lifecycle or annotation controls as appropriate.
4. **Given** an unanswered question and an annotation in the same workspace, **When** the user opens Active Questions, **Then** the question appears once and the annotation does not appear.
5. **Given** no internet connection and an available local Go service, **When** the user creates or reads notes and captures questions or annotations, **Then** all core capture and tracking actions remain available.

---

### User Story 2 - Share a Question Across Notes (Priority: P1)

A user links one question to multiple notes when a topic crosses contexts. The question remains one shared item, and edits made from any linked note or central view are visible in every other location.

**Why this priority**: A shared question prevents duplicate, inconsistent copies and supports cross-cutting topics.

**Independent Test**: Link one question to two notes, edit it from the second note, and verify the revised content in the first note and central Questions view.

**Acceptance Scenarios**:

1. **Given** an existing question, **When** the user links it to another note, **Then** both notes refer to the same question and list each other among its linked notes.
2. **Given** a question linked to multiple notes, **When** the user edits its question text from any linked note, **Then** the updated text appears in every linked note and central question view.
3. **Given** a question in a central view, **When** the user edits it, **Then** the edit is visible from all linked notes.
4. **Given** multiple possible existing questions, **When** the user creates or selects a link, **Then** the application provides enough identifying context to avoid linking the wrong question.

---

### User Story 3 - Answer and Reopen Questions (Priority: P1)

A user records an answer, changes the question's status, and later reopens it to revise the answer. Active and Answered pages are status-based views over the same questions rather than separate copies.

**Why this priority**: Managing the lifecycle of unresolved questions turns note-taking into actionable follow-up.

**Independent Test**: Answer an active question, mark it Answered, find it in the Answered Questions view, reopen it, and revise the preserved answer.

**Acceptance Scenarios**:

1. **Given** a question with no answer, **When** the user attempts to mark it Answered, **Then** the status change is rejected and the user is told that an answer is required.
2. **Given** a question with an answer, **When** the user marks it Answered, **Then** it disappears from the default active view and appears in the Answered Questions view without losing its links or history-relevant dates.
3. **Given** an answered question, **When** the user reopens it, **Then** its existing answer is preserved and the question returns to an active status.
4. **Given** an answered question, **When** its answer is removed, **Then** it cannot remain Answered and the user must choose an active status.
5. **Given** a question, **When** its status changes from a note or central view, **Then** all views reflect the new status.
6. **Given** a question linked to one or more notes, **When** the user edits its answer from any linked note or a central view, **Then** the updated answer appears in every linked note and central question view.

---

### User Story 4 - Prioritize and Find Work (Priority: P2)

A user assigns due dates and reminders, then searches, filters, and sorts questions to decide what to answer next.

**Why this priority**: Central tracking is only useful when a growing question collection can be prioritized and searched efficiently.

**Independent Test**: Create questions with different statuses, tags, dates, and answers; filter and sort the list; then search within a selected scope and open a result.

**Acceptance Scenarios**:

1. **Given** questions with varied metadata, **When** the user filters by status, topic, tag, workspace, creation date, updated date, due date, answer presence, or linked-note presence, **Then** only matching questions are shown; a question matches a topic when at least one of its linked notes has that topic.
2. **Given** a question list, **When** the user sorts by creation date, updated date, due date, priority, title or question text, or status, **Then** the list follows the chosen order.
3. **Given** a question with a past incomplete due date, **When** it is shown in a list, **Then** it is clearly identified as overdue.
4. **Given** an optional reminder, **When** its time is reached while reminder delivery is available, **Then** the user is alerted and can open the question.
5. **Given** search text, **When** the user scopes the search to the current workspace, all workspaces, notes, questions, answers, or everything, **Then** results are limited to that scope and identify result type, workspace, matching context, and destination.
6. **Given** no internet connection and an available local Go service, **When** the user searches, filters, sorts, or manages reminders, **Then** those capabilities remain available within the platform's local notification limits.

---

### User Story 5 - Organize Separate Areas of Life (Priority: P2)

A user creates Work, Personal, or Studies workspaces and organizes notes using parent-child pages, topics, and tags. Each workspace presents its own notes and question views by default.

**Why this priority**: Separate contexts prevent unrelated notes and actions from becoming mixed together.

**Independent Test**: Create two workspaces with notes and questions, switch between them, and verify that default views contain only the current workspace's content.

**Acceptance Scenarios**:

1. **Given** multiple workspaces, **When** the user switches workspace, **Then** notes, topics, tags, and default Active and Answered Questions views change to the selected workspace.
2. **Given** notes in a workspace, **When** the user establishes a parent-child relationship, **Then** the hierarchy is visible and can be navigated.
3. **Given** a workspace-scoped tag, **When** the user enables cross-workspace sharing for that tag, **Then** the tag becomes selectable in other workspaces while remaining the same shared tag.
4. **Given** a tag that has not been explicitly shared, **When** the user works in another workspace, **Then** that tag is not offered there by default.
5. **Given** a note, **When** the user assigns a topic and separate tags, **Then** the topic organizes the note while tags provide independent labels.

---

### User Story 6 - Delete Notes Safely (Priority: P2)

A user deletes a note and makes an informed decision about questions that would otherwise lose their only note link. Questions linked elsewhere are preserved.

**Why this priority**: Safe deletion prevents accidental loss of questions and answers.

**Independent Test**: Delete one note containing both singly linked and multiply linked questions, then verify the chosen result for each affected question.

**Acceptance Scenarios**:

1. **Given** a note is the only linked note for a question, **When** the user deletes the note, **Then** the user must choose to keep the question unlinked, delete the question, or cancel the note deletion.
2. **Given** a question linked to multiple notes, **When** one linked note is deleted, **Then** the application explains that the note link will be removed, requires confirmation, and preserves the question and its other links.
3. **Given** a note with multiple affected questions, **When** deletion is requested, **Then** one review presents every consequence and supports a decision for each singly linked question without repetitive prompts.
4. **Given** the user cancels the deletion review, **When** the review closes, **Then** the note and all links remain unchanged.

---

### User Story 7 - Export and Restore All Data (Priority: P3)

A user exports all application data to a portable package and later imports it to restore the complete application state without silent data loss.

**Why this priority**: Local users need control and portability even before automated backups or cross-device synchronization exist.

**Independent Test**: Build a representative data set, export it, import it into an empty application state, and compare all content, relationships, settings, and dates.

**Acceptance Scenarios**:

1. **Given** application data, **When** the user exports it, **Then** the export includes all workspaces, notes, hierarchy, topics, questions, answers, statuses, links, tags and sharing settings, due dates, reminders, and creation and update dates.
2. **Given** a valid complete export and an empty application state, **When** the user imports it, **Then** the prior application state is restored with all relationships intact.
3. **Given** invalid or unsupported import data, **When** validation runs, **Then** no existing data is changed and the user receives a clear report.
4. **Given** imported records that conflict with existing records, **When** validation completes, **Then** the user can review the conflict and choose to keep existing data, use imported data, or cancel before changes occur.

---

### User Story 8 - Answer in the Original Note Context (Priority: P1)

A user opens a question from Active Questions, Next, or a note and answers it beside the source passage. While reading a note, they can see that note's open questions and jump to the next unanswered gap.

**Why this priority**: The product's unique value is the passage, not a detached todo list. Answering without the sentence that caused the question breaks the primary loop.

**Independent Test**: Create a note, capture two questions on different sentences, open the first from Active Questions, verify the note and highlighted passage are visible next to answer controls, then from the note use Next unanswered to open the second question's card.

**Acceptance Scenarios**:

1. **Given** a question linked to a note, **When** the user opens it from Active Questions, **Then** the same page shows the source note, the passage is highlighted and scrolled into view, and answer/status controls are available without a second navigation.
2. **Given** a question linked to two notes, **When** the user opens it, **Then** the first linked note is shown by default and the user can switch to the other linked note.
3. **Given** an unlinked question, **When** the user opens it, **Then** answer/status controls remain available and the note pane explains that the question is currently unlinked.
4. **Given** a note with unanswered or in-progress questions, **When** the user reads the note, **Then** a note-local list shows those questions in directive order and excludes annotations and answered questions.
5. **Given** multiple open questions on a note, **When** the user chooses Next unanswered, **Then** the next open question in directive order is scrolled into view and its card opens; if none remain, the UI says so.

---

### User Story 9 - Decide What to Answer Next (Priority: P2)

A user opens a single opinionated Next queue for the current workspace instead of building a saved filter.

**Why this priority**: Due dates, priority, and status already exist; they are not usable as a daily workflow until they appear as one default queue.

**Independent Test**: Seed overdue, due-today, in-progress, deferred-ready, high-priority, future-deferred, annotation, and answered items; open Next; verify section order and that annotations, answered questions, and future-deferred questions are absent; open one item into answer-in-context.

**Acceptance Scenarios**:

1. **Given** active questions in the current workspace, **When** the user opens Next, **Then** only `kind=question` items in unanswered, in progress, or deferred appear, grouped and ordered as: Overdue; Due today; In progress; Deferred ready; High priority.
2. **Given** a deferred question whose due date is in the future, **When** Next is opened, **Then** that question is omitted from every section.
3. **Given** an annotation or an answered question, **When** Next is opened, **Then** it does not appear.
4. **Given** an empty section, **When** Next is rendered, **Then** that section is hidden; if the whole queue is empty, the page shows actionable guidance.
5. **Given** a Next item, **When** the user opens it, **Then** they land on the answer-in-context question page.
6. **Given** no internet access and an available local Go service, **When** the user opens Next, **Then** the queue still loads from the local API.

Next is not a saved-filter builder and MUST NOT add a query language or reusable views.

---

### User Story 10 - Defer with a Resume Date (Priority: P2)

A user who defers a question must set a resume date so the question can return in Next instead of disappearing into a junk drawer.

**Why this priority**: Deferred without a wake-up date silently kills the lifecycle the rest of the product depends on.

**Independent Test**: Attempt to mark a question Deferred with no due date and verify rejection; set a date, save as Deferred, verify it appears under Deferred ready or is omitted when the date is in the future; load a legacy deferred question that has no due date and verify it still opens, but saving it as Deferred again requires a date.

**Acceptance Scenarios**:

1. **Given** a question, **When** the user sets status to Deferred without a due date, **Then** the change is rejected and the user is told a resume date is required.
2. **Given** Deferred selected in the lifecycle UI, **When** no due date is present, **Then** a date control is required before save.
3. **Given** an existing deferred question with a null due date, **When** it is opened, **Then** it still loads; a later save that keeps Deferred requires a due date.
4. **Given** an annotation, **When** its comment is edited, **Then** this resume-date rule is not applied.

---

### User Story 11 - Resolved Highlights and Promote Answer (Priority: P2)

After a question is answered, its highlight looks resolved. The user may insert the answer into the note after the passage without turning the note into a live embed.

**Why this priority**: Reading should get clearer over time. Answering should not leave the note looking like an open problem, and the optional promote step folds follow-up back into the note the user owns.

**Independent Test**: Answer a question, reload the note, verify the highlight style is resolved and the card shows the answer; insert the answer into the note, reload, verify a blockquote follows the wrap; insert again and verify the note body is unchanged; reopen the question and verify the highlight uses the active style again.

**Acceptance Scenarios**:

1. **Given** an answered question wrap, **When** the note is read, **Then** that highlight uses a resolved style distinct from unanswered, in-progress, and deferred highlights.
2. **Given** an answered question card, **When** it is opened, **Then** the answer is visible.
3. **Given** an answered wrapped question, **When** the user chooses Insert answer into note, **Then** a Markdown blockquote containing the answer is inserted immediately after the wrap.
4. **Given** that blockquote already follows the wrap, **When** Insert is chosen again, **Then** the note body is not duplicated.
5. **Given** the user answers a question, **When** they do not choose Insert, **Then** the note body is not mutated.
6. **Given** an answered question that is reopened, **When** the note is read, **Then** the highlight uses the active style again. Inserted note text is left in place.

---

### User Story 12 - Capture and Relink from Reading (Priority: P1)

A user captures and links questions from the reading view, including existing questions and questions created without a passage. Selection wrapping tolerates light Markdown/whitespace differences. Keyboard shortcuts cover the capture actions.

**Why this priority**: Shared questions and unlinked follow-ups are already in the data model, but reading is the default surface. Exact substring wrap fails the “exact context” promise on ordinary formatted text.

**Independent Test**: Select rendered bold text whose Markdown is `**bold**` and capture a question; link an existing question onto a new passage from reading; create an unlinked question from Active Questions and attach it to a passage later; exercise Q / A / L / Escape with a selection and verify typing in the composer is not stolen.

**Acceptance Scenarios**:

1. **Given** a selection that matches note Markdown after collapsing whitespace or stripping emphasis markers, **When** the user captures a question or annotation, **Then** the wrap is applied to the source Markdown and no raw token is shown in reading view.
2. **Given** a selection that cannot be mapped to note Markdown, **When** capture is attempted, **Then** no question, annotation, or note mutation occurs and the user sees a clear error.
3. **Given** a text selection in reading view, **When** the user chooses Link existing question and picks a question not already linked to this note, **Then** the passage is wrapped with that question id and the canonical question is unchanged except for the new link.
4. **Given** Active Questions, **When** the user creates a question without a passage, **Then** it appears in Active Questions as unlinked.
5. **Given** an unlinked question and a reading selection, **When** the user attaches it to the passage, **Then** the wrap and note link are saved and the question remains the same record.
6. **Given** a reading selection and no open composer or focused text field, **When** the user presses Q, A, L, or Escape, **Then** those keys ask, annotate, link existing, or dismiss. The same keys do nothing while typing in an input, textarea, or composer.

### Edge Cases

- A question may remain intentionally unlinked after its last linked note is deleted; it remains accessible through central question views and search.
- Duplicate question text does not imply duplicate identity; users can keep distinct questions with identical wording.
- Empty or whitespace-only question text and answers are treated as missing content.
- Circular parent-child note relationships are rejected, and moving a parent beneath one of its descendants leaves the hierarchy unchanged.
- Deleting a parent note does not silently delete its child notes; the user must review whether children are re-parented or deleted.
- A reminder whose due time passed while the application could not deliver notifications is shown as missed or overdue when the application is next available.
- Changes to a shared tag are visible in every workspace where it was explicitly shared; disabling sharing does not remove the tag from content that already uses it without confirmation.
- Search handles no matches, punctuation-only input, and identical matches in multiple workspaces without presenting one as the other.
- An interrupted or failed import leaves the pre-import application state intact.
- Exported content remains complete when questions are unlinked or when notes have nested descendants.

## Requirements *(mandatory)*

### Functional Requirements

#### Notes and question links

- **FR-001**: Users MUST be able to create, view, edit, and delete notes without internet access while the local Go service is running.
- **FR-002**: Users MUST be able to arrange notes in navigable parent-child relationships.
- **FR-003**: The application MUST prevent circular parent-child relationships.
- **FR-004**: Notes MUST support rapid entry of paragraphs, ordered lists, and unordered lists; the exact editing and storage format is deferred to planning and evaluation.
- **FR-005**: A note MUST be able to contain multiple questions.
- **FR-006**: A question MUST be linkable to multiple notes within its owning workspace.
- **FR-007**: Each question MUST remain one shared item regardless of how many notes and views display it.
- **FR-008**: Users MUST be able to create, open, and edit a question from a linked note or a central question view.
- **FR-009**: Changes to a question's text, answer, status, priority, due date, reminder, and tags MUST be reflected in every location where it appears, without requiring duplicate edits.
- **FR-010**: Users MUST be able to select a passage in a note and attach a question or annotation; the passage MUST render as a sanitized clickable highlight that opens a card, and directive tokens MUST never be visible in reading view. Legacy expanded, collapsed, and link-only display modes may remain stored for compatibility.
- **FR-011**: Users MUST be able to see and navigate to every note linked to a question.

#### Question lifecycle and prioritization

- **FR-012**: Questions MUST support Unanswered, In Progress, Deferred, and Answered statuses.
- **FR-013**: A question MUST have a non-empty answer before it can enter Answered status.
- **FR-014**: Users MUST be able to reopen an answered question while preserving its answer.
- **FR-015**: Removing the answer from an answered question MUST require selection of a non-answered status before the change is completed.
- **FR-016**: The Active Questions and Answered Questions pages MUST be filtered views over the same question collection, not separate copies.
- **FR-017**: Questions MUST support optional due dates, reminders, and priority.
- **FR-018**: The application MUST visibly distinguish overdue, unanswered questions.
- **FR-019**: Reminder behavior MUST remain useful without internet access and clearly communicate when browser or operating-system conditions prevent timely delivery.

#### Find and organize

- **FR-020**: Users MUST be able to filter questions by status, topic, tag, workspace, creation date, updated date, due date, answer presence, and linked-note presence. A question MUST match a topic filter when at least one linked note has that topic.
- **FR-021**: Users MUST be able to sort questions by creation date, updated date, due date, priority, alphabetical title or question text, and status.
- **FR-022**: Users MUST be able to search the current workspace, all workspaces, notes, questions, answers, or all content.
- **FR-023**: Every search result MUST identify its content type, workspace, matching context, and navigable destination.
- **FR-024**: Users MUST be able to create and switch among workspaces, each with its own notes, hierarchy, topics, default tags, and Active and Answered Questions views.
- **FR-025**: Each note MUST support one directly assigned topic and multiple tags.
- **FR-026**: Tags MUST be workspace-scoped by default and MUST become available in other workspaces only through an explicit sharing control.
- **FR-027**: A question MUST belong to one workspace and MUST link only to notes in that workspace. Moving notes or questions between workspaces is not supported in the MVP.

#### Safe deletion

- **FR-028**: Deleting the sole linked note of a question MUST require the user to keep the question unlinked, delete the question, or cancel.
- **FR-029**: Deleting one of several notes linked to a question MUST preserve the question, explain that only the selected note link will be removed, and require confirmation.
- **FR-030**: A deletion review involving multiple questions MUST summarize all consequences and allow an explicit decision for each question at risk of becoming unlinked.
- **FR-031**: Cancelling a deletion review MUST leave the note, questions, and links unchanged.

#### Offline operation and portability

- **FR-032**: Core MVP capabilities MUST work without internet access while the local Go service is running, including note and question management, status changes, workspace organization, search, sorting, filtering, reminders subject to platform limits, import, and export.
- **FR-033**: The MVP MUST NOT require authentication or an online account.
- **FR-034**: Users MUST be able to export all application data in a documented, portable form.
- **FR-035**: Exports MUST preserve workspaces, notes, note hierarchy, topics, questions, answers, statuses, question-note links, tags and sharing settings, priorities, due dates, reminders, and creation and update dates.
- **FR-036**: Users MUST be able to validate an export and import it to restore a complete application state.
- **FR-037**: Import validation MUST report invalid, unsupported, duplicate, or conflicting content before changing existing data.
- **FR-038**: For conflicts, users MUST be able to keep existing content, use imported content, or cancel the import; no content may be overwritten silently.
- **FR-039**: Failed or cancelled imports MUST leave the pre-import application state unchanged.
- **FR-040**: The MVP MUST support the web platform with all core capabilities available without internet access while the local Go service is running. Desktop and Android support are deferred to future releases.

#### Reading context, next queue, and capture

- **FR-041**: Opening a question MUST show answer/status controls on the same page as the source note excerpt when a linked note exists, with the matching passage highlighted and scrolled into view.
- **FR-042**: When a question is linked to multiple notes, the context page MUST default to the first linked note and MUST allow switching among linked notes.
- **FR-043**: Unlinked questions MUST still open for answering and MUST explain that no source note is linked.
- **FR-044**: A note reading view MUST list that note's open questions (`kind=question`, status not `answered`) in directive order, excluding annotations.
- **FR-045**: A Next unanswered action on a note MUST open the next open question in directive order, scroll to its highlight, and open its card, or announce that none remain.
- **FR-046**: The application MUST provide one workspace-scoped Next queue. It MUST NOT be a saved-filter builder.
- **FR-047**: Next MUST include only current-workspace `kind=question` items with status unanswered, in_progress, or deferred, grouped in this order, hiding empty groups: Overdue (`due_date < today` and status is not deferred); Due today (`due_date = today` and status is not deferred); In progress (status `in_progress` and not already in Overdue or Due today); Deferred ready (status `deferred` and `due_date <= today`); High priority (priority `high` or `urgent`, status unanswered or in_progress, and not already listed). Future-deferred, answered, and annotation items MUST be omitted.
- **FR-048**: Transitioning a question to Deferred MUST require a non-empty due date. Existing deferred rows with a null due date MUST still load. Annotations are exempt.
- **FR-049**: Reading view MUST style answered-question highlights as resolved and distinct from active question highlights. Style is driven by the live question status.
- **FR-050**: Users MUST be able to insert an answered question's answer into the note as a Markdown blockquote immediately after the wrap. The action MUST be user-initiated, idempotent if that blockquote already follows the wrap, and MUST NOT run automatically on answer.
- **FR-051**: Passage wrapping MUST match note Markdown by exact substring first, then collapsed whitespace, then emphasis-stripped text (`*`, `_`, `**`, `` ` ``). If no match, the application MUST mutate nothing and MUST show an error.
- **FR-052**: Reading-view selection MUST be able to link an existing same-workspace question that is not already linked to the note.
- **FR-053**: Users MUST be able to create a question with no passage and later attach it to a selected passage from reading view. The question remains one canonical record.
- **FR-054**: With a reading selection and no focused text field or open composer, Q MUST start a question, A an annotation, L link-existing, and Escape MUST dismiss selection actions. Those shortcuts MUST NOT fire while typing in an input, textarea, or composer.
- **FR-055**: Opening a Next or Active Questions item MUST land on the answer-in-context question page, including `noteId` when a linked note exists.

### Scope Boundaries

The MVP includes structured text notes, question tracking, workspaces, topics, tags, parent-child note organization, due dates, reminders, search, sorting, filtering, offline operation, and complete import/export.

Post-MVP deepening in this specification (US8–US12) includes answer-in-context, a note-local open-question rail, one opinionated Next queue, deferred resume dates, resolved highlights with optional answer insertion, sturdier passage matching, reading-view link/attach, and keyboard capture.

The following are outside the MVP:

- Saved filters and reusable views
- Revision history
- Automatic or manual backup workflows
- Local or cloud backup destinations and backup encryption
- Cross-device synchronization and conflict resolution
- Moving notes or questions between workspaces
- Images and file attachments
- Advanced links or previews and code blocks
- Separate user profiles and user switching
- Authentication
- Collaboration, permissions, and sharing
- Mobile platforms not selected for the first release

### Key Entities *(include if feature involves data)*

- **Workspace**: A separate area such as Work, Personal, or Studies; owns notes, topics, questions, and workspace-scoped tags.
- **Note**: User-authored content with a title, structured text, optional parent, child notes, one optional topic, tags, linked questions, workspace, and creation and update dates.
- **Question**: A shared actionable inquiry or passage annotation with kind (`question` or `annotation`), question/comment text, optional answer for questions, status, priority, optional due date and reminder, tags, linked notes, owning workspace, and creation and update dates. Only questions enter Active and Answered views.
- **Question Link**: The relationship between one question and one note, including the question's presentation preference in that note.
- **Topic**: A workspace-owned organizational category directly assigned to notes.
- **Tag**: A flexible label that can classify notes and questions; belongs to a workspace by default and can be explicitly shared with additional workspaces.
- **Reminder**: An optional alert associated with a question and a scheduled date and time.
- **Export Package**: A portable representation of all user data and relationships needed for validation, import, and complete restoration.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In usability testing, at least 90% of users can create a note, add a question, and find it in the Active Questions view on their first attempt without assistance.
- **SC-002**: Users can create a note containing a question in under 60 seconds and link an existing question to a second note in under 30 seconds.
- **SC-003**: A question edit is visible from all linked notes and central views within one second of the user completing the edit under normal local use.
- **SC-004**: At least 95% of searches over collections containing up to 10,000 notes and 50,000 questions present complete matching results within two seconds under normal local use.
- **SC-005**: Users can filter and sort a 50,000-question collection and see the resulting view within two seconds under normal local use.
- **SC-006**: All core MVP journeys can be completed without internet access while the local Go service remains available, excluding only notification behavior that the active platform explicitly does not permit.
- **SC-007**: In deletion tests, 100% of questions linked to multiple notes survive deletion of any single linked note, and no singly linked question is deleted without an explicit user choice.
- **SC-008**: A complete export-import round trip preserves 100% of supported content, relationships, statuses, settings, and dates in the acceptance data set.
- **SC-009**: Invalid, cancelled, interrupted, or conflicting imports cause no unapproved loss or overwrite of existing data in all acceptance tests.
- **SC-010**: At least 85% of usability-test participants rate the question capture and follow-up workflow as easy or very easy.
- **SC-011**: In usability testing, at least 90% of users can open a question from Active Questions and see the source passage on the same page on the first attempt without assistance.
- **SC-012**: Given a seeded Next fixture, 100% of automated tests place overdue, due-today, in-progress, deferred-ready, and high-priority items in the specified sections and omit answered, annotation, and future-deferred items.
- **SC-013**: 100% of attempts to mark a question Deferred without a due date are rejected in domain, API, and browser tests.
- **SC-014**: Passage wrap succeeds for selections that differ from source Markdown only by collapsed whitespace or emphasis markers, and fails with no mutation when the passage cannot be mapped.

## Assumptions

- The MVP is web-only. Desktop and Android remain target platforms for future releases.
- The MVP serves one local user and one active local application state; profiles and authentication are future features.
- A question belongs to exactly one workspace and links only to notes in that workspace. This preserves workspace-specific master views and avoids unclear cross-workspace ownership. Cross-workspace note and question moves are deferred beyond the MVP.
- “Offline” and “without internet access” mean the local Go service remains running and reachable; the MVP does not support editing while that local service is stopped.
- A note has zero or one directly assigned topic. Tags cover overlapping or multiple classifications.
- Topics are workspace-specific and are not shared across workspaces.
- Newly created questions begin with Unanswered status unless the user explicitly selects another valid status.
- Reopening an answered question changes it to In Progress by default while preserving the answer.
- Question title sorting uses the question text when no separate title exists.
- Import conflicts are reviewed before import, with per-conflict choices to retain existing content, use imported content, or cancel.
- Reminder delivery uses the notification capabilities available to the selected platform and reports reminders missed while delivery was unavailable.
- The user-facing editor must support rapid structured note-taking, but selection among plain text, Markdown, and rich text is a planning decision informed by usability evaluation.
- Links, images, files, and code blocks are deferred even if the selected editor could support them cheaply; this keeps MVP acceptance criteria stable.
- Saved views, revision history, backups, synchronization, profiles, authentication, and collaboration are future features and do not block the MVP.
