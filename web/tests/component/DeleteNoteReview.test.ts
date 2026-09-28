import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import DeleteNoteReview from '$lib/components/DeleteNoteReview.svelte';

describe('DeleteNoteReview', () => {
	it('shows explicit consequences and cancel', () => {
		render(DeleteNoteReview, {
			preview: {
				previewToken: 't',
				expiresAt: '',
				noteId: 'n',
				noteVersion: 1,
				singlyLinkedQuestions: [
					{
						id: 'q',
						workspaceId: 'w',
						questionText: 'Keep?',
						status: 'unanswered',
						priority: 'none',
						dueDate: null,
						version: 1
					}
				],
				multiplyLinkedQuestions: [],
				childNotes: []
			}
		});
		expect(screen.getByRole('heading', { name: /review note deletion/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument();
	});
});
