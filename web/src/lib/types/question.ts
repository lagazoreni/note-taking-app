export type QuestionStatus = 'unanswered' | 'in_progress' | 'deferred' | 'answered';
export type Priority = 'none' | 'low' | 'medium' | 'high' | 'urgent';
export type DisplayMode = 'expanded' | 'collapsed' | 'link';
export type QuestionKind = 'question' | 'annotation';

export interface Reminder {
	id: string;
	scheduledAt: string;
	state: 'pending' | 'delivered' | 'missed' | 'dismissed';
	lastEvaluatedAt: string | null;
	createdAt: string;
	updatedAt: string;
	version: number;
}

export interface QuestionSummary {
	id: string;
	workspaceId: string;
	questionText: string;
	status: QuestionStatus;
	priority: Priority;
	dueDate: string | null;
	version: number;
	kind?: QuestionKind;
}

export interface LinkedNoteSummary {
	id: string;
	title: string;
	displayMode: DisplayMode;
}

export interface Question extends QuestionSummary {
	answerMarkdown: string | null;
	reminder: Reminder | null;
	tagIds: string[];
	linkedNotes: LinkedNoteSummary[];
	createdAt: string;
	updatedAt: string;
}

export interface QuestionPage {
	items: Question[];
	nextCursor?: string | null;
}
