# Noted — Project Overview

## What this project is about

Noted is a local-first note-taking application for turning reading and note-taking into actionable follow-up. It lets one local user write and organize notes, select passages that need clarification or context, and attach either a question or an annotation directly to those passages.

Questions are shared, canonical items rather than copies of text embedded in individual notes. A question can appear in several notes and in the central question views, so an answer, status change, or other edit is reflected everywhere. Annotations are comments on passages and remain available from their highlights without adding noise to the Active Questions inbox.

The application is organized around workspaces such as Work, Personal, or Studies. Each workspace keeps its notes, hierarchy, topics, tags, and default question views separate, while explicitly shared tags may be used in more than one workspace.

## Product goal

**Make reading interactive by connecting the context of a note with a dependable workflow for returning to unresolved questions.**

The primary experience is:

1. Create or select a workspace.
2. Write or open a note from the workspace Notes index.
3. Read the note and select a passage.
4. Add a question or annotation from the selection.
5. See the passage as a sanitized highlight, without exposing internal directive tokens.
6. Open the highlight to view the question or annotation card.
7. Find unanswered questions in Active Questions, answer or defer them, and reopen them later when needed.

The goal is not merely to store text. It is to preserve the relationship between a thought and the exact note context that produced it, while making unresolved questions easy to find and act on.

## Core capabilities

- Create, read, edit, and delete Markdown notes, including paragraphs and ordered or unordered lists.
- Arrange notes in parent-child hierarchies without allowing circular relationships.
- Capture multiple questions or annotations from passages in a note.
- Render captured passages as clickable highlights or marker chips instead of raw `{{question:<uuid>}}` tokens.
- Keep each question as one shared record across all linked notes and central views.
- Manage question lifecycle states: Unanswered, In Progress, Deferred, and Answered. An answer is required before a question can be marked Answered.
- Prioritize questions with priority, due dates, reminders, filters, sorting, and full-text search across notes, questions, and answers.
- Keep annotations out of Active Questions and Answered Questions while retaining them on their source passages.
- Organize content with workspaces, one topic per note, multiple tags, and explicit cross-workspace tag sharing.
- Preview note deletion and require decisions before a sole linked question or child note can be lost.
- Export all supported data and relationships to a documented portable archive, then validate, review, and restore it atomically without silent overwrites.

## Local-first model

Noted is designed for a single local user and does not require authentication, an online account, or an internet connection for core use. “Offline” means that the local Go service is still running and reachable while external internet access is unavailable. The browser provides the application shell, but the Go API is the authoritative source for validation, business rules, persistence, and question/note relationships.

The application deliberately avoids silently queuing or replaying writes in the browser. If the local service is unavailable, unsaved edits remain visible and retryable rather than being reported as successfully saved.

## High-level architecture

```text
SvelteKit browser application
        │ same-origin /api/v1 requests
        ▼
Go HTTP API — validation, lifecycle, search, import/export
        ▼
SQLite database with FTS5 — canonical local data
```

The frontend is a static SvelteKit application. In production it is served by the Go process from the same origin. SQLite stores workspaces, notes, questions, links, topics, tags, reminders, and import state; full-text indexes support local search.

## MVP boundaries

The first release is web-only and focused on one local application state. It does not include authentication, collaboration, permissions, cross-device synchronization, moving notes or questions between workspaces, revision history, backup workflows, images or file attachments, advanced links/previews, saved filters, or native desktop and Android clients.

## Definition of success

Noted succeeds when a user can move from a passage to a useful follow-up with little friction: create a note, capture a question, find it centrally, and return to the original context. The product should keep edits consistent across linked views, preserve user data during deletion and import operations, remain useful without internet access, and make the question-capture workflow feel easy and reliable.
