import { describe, expect, it } from 'vitest';
import { buildNextQueue } from '$lib/questions/nextQueue';
import type { Question } from '$lib/types/question';

const today = '2026-04-01';

function makeQuestion(id: string, overrides: Partial<Question> = {}): Question {
	return {
		id,
		workspaceId: '11111111-1111-4111-8111-111111111111',
		questionText: `Question ${id}`,
		status: 'unanswered',
		priority: 'none',
		dueDate: null,
		answerMarkdown: null,
		reminder: null,
		tagIds: [],
		linkedNotes: [],
		createdAt: '2026-04-01T00:00:00.000Z',
		updatedAt: '2026-04-01T00:00:00.000Z',
		version: 1,
		kind: 'question',
		...overrides
	};
}

describe('buildNextQueue', () => {
	it('groups active questions in priority order without duplicating items', () => {
		const questions = [
			makeQuestion('overdue', { dueDate: '2026-03-31' }),
			makeQuestion('overdue-high', { dueDate: '2026-03-30', priority: 'urgent' }),
			makeQuestion('due-today', { dueDate: today }),
			makeQuestion('due-today-progress', {
				dueDate: today,
				status: 'in_progress',
				priority: 'high'
			}),
			makeQuestion('in-progress', { status: 'in_progress' }),
			makeQuestion('in-progress-future', {
				status: 'in_progress',
				dueDate: '2026-04-02',
				priority: 'urgent'
			}),
			makeQuestion('deferred-ready', { status: 'deferred', dueDate: '2026-03-31' }),
			makeQuestion('deferred-today', { status: 'deferred', dueDate: today }),
			makeQuestion('high-priority', { priority: 'high' }),
			makeQuestion('legacy-high', { kind: undefined, priority: 'urgent' }),
			makeQuestion('answered', {
				status: 'answered',
				answerMarkdown: 'Done',
				dueDate: '2026-03-30',
				priority: 'urgent'
			}),
			makeQuestion('annotation', {
				kind: 'annotation',
				dueDate: '2026-03-30',
				priority: 'urgent'
			}),
			makeQuestion('deferred-future', {
				status: 'deferred',
				dueDate: '2026-04-02',
				priority: 'urgent'
			}),
			makeQuestion('deferred-no-date', {
				status: 'deferred',
				priority: 'urgent'
			})
		];

		const sections = buildNextQueue(questions, today);

		expect(
			sections.map(({ title, items }) => ({
				title,
				ids: items.map((item) => item.id)
			}))
		).toEqual([
			{ title: 'Overdue', ids: ['overdue', 'overdue-high'] },
			{ title: 'Due today', ids: ['due-today', 'due-today-progress'] },
			{ title: 'In progress', ids: ['in-progress', 'in-progress-future'] },
			{ title: 'Deferred ready', ids: ['deferred-ready', 'deferred-today'] },
			{ title: 'High priority', ids: ['high-priority', 'legacy-high'] }
		]);

		const ids = sections.flatMap((section) => section.items.map((item) => item.id));
		expect(new Set(ids).size).toBe(ids.length);
		expect(ids).not.toContain('answered');
		expect(ids).not.toContain('annotation');
		expect(ids).not.toContain('deferred-future');
		expect(ids).not.toContain('deferred-no-date');
	});

	it('omits empty sections and excluded questions', () => {
		const sections = buildNextQueue(
			[
				makeQuestion('answered', { status: 'answered' }),
				makeQuestion('annotation', { kind: 'annotation', priority: 'high' }),
				makeQuestion('deferred-future', { status: 'deferred', dueDate: '2026-04-02' }),
				makeQuestion('deferred-no-date', { status: 'deferred' })
			],
			today
		);

		expect(sections).toEqual([]);
	});
});
