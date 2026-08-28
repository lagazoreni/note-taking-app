ALTER TABLE questions ADD COLUMN kind TEXT NOT NULL DEFAULT 'question' CHECK (kind IN ('question', 'annotation'));
