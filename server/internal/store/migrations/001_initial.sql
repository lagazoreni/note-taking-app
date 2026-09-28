PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    checksum TEXT NOT NULL CHECK (length(checksum) = 64),
    applied_at TEXT NOT NULL
);

CREATE TABLE workspaces (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 100),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    UNIQUE (name COLLATE NOCASE)
);

CREATE TABLE topics (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE RESTRICT,
    name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 100),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    UNIQUE (workspace_id, name COLLATE NOCASE)
);

CREATE TABLE notes (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE RESTRICT,
    topic_id TEXT REFERENCES topics(id) ON DELETE SET NULL,
    parent_note_id TEXT REFERENCES notes(id) ON DELETE RESTRICT,
    title TEXT NOT NULL CHECK (length(trim(title)) BETWEEN 1 AND 300),
    body_markdown TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    CHECK (parent_note_id IS NULL OR parent_note_id <> id)
);

CREATE TABLE questions (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE RESTRICT,
    question_text TEXT NOT NULL CHECK (length(trim(question_text)) BETWEEN 1 AND 10000),
    answer_markdown TEXT,
    status TEXT NOT NULL CHECK (status IN ('unanswered', 'in_progress', 'deferred', 'answered')),
    priority TEXT NOT NULL DEFAULT 'none' CHECK (priority IN ('none', 'low', 'medium', 'high', 'urgent')),
    due_date TEXT CHECK (due_date IS NULL OR due_date GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]'),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    CHECK (status <> 'answered' OR (answer_markdown IS NOT NULL AND length(trim(answer_markdown)) > 0))
);

CREATE TABLE note_questions (
    note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE RESTRICT,
    question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    display_mode TEXT NOT NULL DEFAULT 'collapsed' CHECK (display_mode IN ('expanded', 'collapsed', 'link')),
    position INTEGER NOT NULL CHECK (position >= 0),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    PRIMARY KEY (note_id, question_id),
    UNIQUE (note_id, position)
);

CREATE TABLE tags (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    owner_workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE RESTRICT,
    name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 100),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    UNIQUE (owner_workspace_id, name COLLATE NOCASE)
);

CREATE TABLE tag_workspace_access (
    tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    workspace_id TEXT NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL,
    PRIMARY KEY (tag_id, workspace_id)
);

CREATE TABLE note_tags (
    note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE RESTRICT,
    tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE RESTRICT,
    created_at TEXT NOT NULL,
    PRIMARY KEY (note_id, tag_id)
);

CREATE TABLE question_tags (
    question_id TEXT NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE RESTRICT,
    created_at TEXT NOT NULL,
    PRIMARY KEY (question_id, tag_id)
);

CREATE TABLE reminders (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    question_id TEXT NOT NULL UNIQUE REFERENCES questions(id) ON DELETE CASCADE,
    scheduled_at TEXT NOT NULL,
    state TEXT NOT NULL CHECK (state IN ('pending', 'delivered', 'missed', 'dismissed')),
    last_evaluated_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0)
);

CREATE TABLE import_sessions (
    id TEXT PRIMARY KEY CHECK (lower(id) = id),
    archive_path TEXT NOT NULL,
    archive_sha256 TEXT NOT NULL CHECK (archive_sha256 GLOB '[a-f0-9]*' AND length(archive_sha256) = 64),
    format_version INTEGER NOT NULL,
    state TEXT NOT NULL CHECK (state IN ('validated', 'applied', 'cancelled', 'expired')),
    conflict_count INTEGER NOT NULL CHECK (conflict_count >= 0),
    created_at TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    applied_at TEXT
);

CREATE TABLE import_conflicts (
    import_session_id TEXT NOT NULL REFERENCES import_sessions(id) ON DELETE CASCADE,
    conflict_id TEXT NOT NULL CHECK (lower(conflict_id) = conflict_id),
    entity_type TEXT NOT NULL CHECK (entity_type IN ('workspace', 'topic', 'note', 'question', 'noteQuestion', 'tag', 'tagWorkspaceAccess', 'noteTag', 'questionTag', 'reminder')),
    entity_id TEXT NOT NULL,
    reason TEXT NOT NULL CHECK (reason IN ('different_content', 'missing_dependency', 'name_collision', 'newer_existing')),
    existing_version INTEGER CHECK (existing_version IS NULL OR existing_version > 0),
    imported_version INTEGER CHECK (imported_version IS NULL OR imported_version > 0),
    PRIMARY KEY (import_session_id, conflict_id)
);

CREATE INDEX notes_workspace_parent_title_id ON notes(workspace_id, parent_note_id, title, id);
CREATE INDEX notes_workspace_topic_updated_id ON notes(workspace_id, topic_id, updated_at, id);
CREATE INDEX questions_workspace_status_updated_id ON questions(workspace_id, status, updated_at, id);
CREATE INDEX questions_workspace_due_priority_id ON questions(workspace_id, due_date, priority, id);
CREATE INDEX questions_workspace_created_id ON questions(workspace_id, created_at, id);
CREATE INDEX note_questions_question_note ON note_questions(question_id, note_id);
CREATE INDEX note_questions_note_position ON note_questions(note_id, position);
CREATE INDEX note_tags_tag_note ON note_tags(tag_id, note_id);
CREATE INDEX question_tags_tag_question ON question_tags(tag_id, question_id);
CREATE INDEX tag_workspace_access_workspace_tag ON tag_workspace_access(workspace_id, tag_id);
CREATE INDEX reminders_state_scheduled_id ON reminders(state, scheduled_at, id);
CREATE INDEX import_sessions_state_expires ON import_sessions(state, expires_at);

CREATE VIRTUAL TABLE notes_fts USING fts5(
    note_id UNINDEXED,
    workspace_id UNINDEXED,
    title,
    body_markdown
);

CREATE VIRTUAL TABLE questions_fts USING fts5(
    question_id UNINDEXED,
    workspace_id UNINDEXED,
    question_text,
    answer_markdown
);

CREATE TRIGGER notes_fts_insert AFTER INSERT ON notes BEGIN
    INSERT INTO notes_fts(note_id, workspace_id, title, body_markdown)
    VALUES (new.id, new.workspace_id, new.title, new.body_markdown);
END;

CREATE TRIGGER notes_fts_update AFTER UPDATE OF workspace_id, title, body_markdown ON notes BEGIN
    DELETE FROM notes_fts WHERE note_id = old.id;
    INSERT INTO notes_fts(note_id, workspace_id, title, body_markdown)
    VALUES (new.id, new.workspace_id, new.title, new.body_markdown);
END;

CREATE TRIGGER notes_fts_delete AFTER DELETE ON notes BEGIN
    DELETE FROM notes_fts WHERE note_id = old.id;
END;

CREATE TRIGGER questions_fts_insert AFTER INSERT ON questions BEGIN
    INSERT INTO questions_fts(question_id, workspace_id, question_text, answer_markdown)
    VALUES (new.id, new.workspace_id, new.question_text, COALESCE(new.answer_markdown, ''));
END;

CREATE TRIGGER questions_fts_update AFTER UPDATE OF workspace_id, question_text, answer_markdown ON questions BEGIN
    DELETE FROM questions_fts WHERE question_id = old.id;
    INSERT INTO questions_fts(question_id, workspace_id, question_text, answer_markdown)
    VALUES (new.id, new.workspace_id, new.question_text, COALESCE(new.answer_markdown, ''));
END;

CREATE TRIGGER questions_fts_delete AFTER DELETE ON questions BEGIN
    DELETE FROM questions_fts WHERE question_id = old.id;
END;
