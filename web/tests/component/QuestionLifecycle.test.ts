import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';

describe('QuestionLifecycle', () => {
	it('requires answer text before answered status', () => {
		render(QuestionLifecycle, {
			question: {
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
				version: 1
			}
		});
		expect(screen.getByLabelText('Answer')).toBeInTheDocument();
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
	});
});
