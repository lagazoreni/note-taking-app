import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import QuestionList from '$lib/components/QuestionList.svelte';

describe('QuestionList', () => {
	it('shows an empty state when there are no questions', () => {
		render(QuestionList, { questions: [] });
		expect(screen.getByText(/no questions/i)).toBeInTheDocument();
	});
});
