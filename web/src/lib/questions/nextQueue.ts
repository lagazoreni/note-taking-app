import type { Question } from '$lib/types/question';

export interface NextQueueSection {
	id: string;
	title: string;
	items: Question[];
}

/**
 * Group active questions into the sections shown by the daily Next queue.
 *
 * A question is assigned to at most one section. The comparisons intentionally
 * use the YYYY-MM-DD strings supplied by the API and caller, so they remain
 * date-only comparisons rather than becoming timezone-sensitive timestamps.
 */
export function buildNextQueue(questions: Question[], today: string): NextQueueSection[] {
	const active = questions.filter(
		(question) =>
			(question.kind === undefined || question.kind === 'question') &&
			(question.status === 'unanswered' ||
				question.status === 'in_progress' ||
				question.status === 'deferred')
	);

	const listed = new Set<string>();
	const sections: NextQueueSection[] = [];

	function addSection(id: string, title: string, items: Question[]): void {
		const unlisted = items.filter((question) => !listed.has(question.id));
		if (unlisted.length === 0) return;

		for (const question of unlisted) listed.add(question.id);
		sections.push({ id, title, items: unlisted });
	}

	addSection(
		'overdue',
		'Overdue',
		active.filter(
			(question) =>
				question.status !== 'deferred' &&
				question.dueDate !== null &&
				question.dueDate < today
		)
	);

	addSection(
		'due-today',
		'Due today',
		active.filter(
			(question) =>
				question.status !== 'deferred' && question.dueDate === today
		)
	);

	addSection(
		'in-progress',
		'In progress',
		active.filter((question) => question.status === 'in_progress')
	);

	addSection(
		'deferred-ready',
		'Deferred ready',
		active.filter(
			(question) =>
				question.status === 'deferred' &&
				question.dueDate !== null &&
				question.dueDate <= today
		)
	);

	addSection(
		'high-priority',
		'High priority',
		active.filter(
			(question) =>
			(question.priority === 'high' || question.priority === 'urgent') &&
			(question.status === 'unanswered' || question.status === 'in_progress')
		)
	);

	return sections;
}
