import type { DisplayMode, QuestionSummary } from './question';

export interface NoteQuestionLinkWrite {
	questionId: string;
	displayMode: DisplayMode;
	position: number;
}

export interface NoteQuestionLink extends NoteQuestionLinkWrite {
	question: QuestionSummary;
}

export interface NoteWrite {
	workspaceId: string;
	topicId: string | null;
	parentNoteId: string | null;
	title: string;
	bodyMarkdown: string;
	questionLinks: NoteQuestionLinkWrite[];
	tagIds: string[];
}

export interface Note extends NoteWrite {
	id: string;
	createdAt: string;
	updatedAt: string;
	version: number;
	questionLinks: NoteQuestionLink[];
}

export interface NotePage {
	items: Note[];
	nextCursor?: string | null;
}
