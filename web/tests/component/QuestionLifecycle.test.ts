import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { questionsApi } from '$lib/api/questions';
import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
import type { Question } from '$lib/types/question';

vi.mock('$lib/api/questions', () => ({
	questionsApi: {
		update: vi.fn()
	}
}));

function makeQuestion(overrides: Partial<Question> = {}): Question {
	return {
		id: 'q',
		workspaceId: 'w',
		questionText: 'Why?',
		answerMarkdown: null,
		status: 'unanswered',
		priority: 'none',
		dueDate: null,
		reminder: null,
		tagIds: [],
		linkedNotes: [],
		createdAt: '',
		updatedAt: '',
		version: 1,
		kind: 'question',
		...overrides
	};
}

beforeEach(() => {
	const update = vi.mocked(questionsApi.update);
	update.mockReset();
	update.mockResolvedValue(makeQuestion());
});

describe('QuestionLifecycle', () => {
	it('requires answer text before answered status', () => {
		render(QuestionLifecycle, { question: makeQuestion() });
		expect(screen.getByLabelText('Answer')).toBeInTheDocument();
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
	});

	it('rejects deferred status without a resume date', async () => {
		const update = vi.mocked(questionsApi.update);
		render(QuestionLifecycle, { question: makeQuestion() });

		await fireEvent.change(screen.getByLabelText('Status'), { target: { value: 'deferred' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Save question' }));

		expect(
			screen.getByText('A resume date is required to defer a question.')
		).toBeInTheDocument();
		expect(update).not.toHaveBeenCalled();
	});

	it('saves deferred status with its resume date', async () => {
		const update = vi.mocked(questionsApi.update);
		update.mockResolvedValue(makeQuestion({ status: 'deferred', dueDate: '2026-04-15' }));
		render(QuestionLifecycle, {
			question: makeQuestion({ dueDate: '2026-04-15' })
		});

		await fireEvent.change(screen.getByLabelText('Status'), { target: { value: 'deferred' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Save question' }));

		await waitFor(() => expect(update).toHaveBeenCalledTimes(1));
		expect(update).toHaveBeenCalledWith(
			'q',
			expect.objectContaining({ status: 'deferred', dueDate: '2026-04-15' })
		);
	});
});
