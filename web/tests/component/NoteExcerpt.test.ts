import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import NoteExcerpt from '$lib/components/NoteExcerpt.svelte';

const questionId = '550e8400-e29b-41d4-a716-446655440000';
const markdown = `Lead {{question:${questionId}}}selected passage{{/question}} tail`;

describe('NoteExcerpt', () => {
	it('renders a highlight mark for the given question without capture actions', () => {
		render(NoteExcerpt, { markdown, questionId });
		const mark = document.querySelector(`mark[data-annotation-id="${questionId}"]`);
		expect(mark).toBeInTheDocument();
		expect(mark).toHaveTextContent('selected passage');
		expect(screen.queryByRole('button', { name: 'Ask a question' })).not.toBeInTheDocument();
		expect(screen.queryByRole('button', { name: 'Add annotation' })).not.toBeInTheDocument();
	});
});
