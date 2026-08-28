import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import QuestionFilters from '$lib/components/QuestionFilters.svelte';

describe('QuestionDiscovery', () => {
	it('renders filter and sort controls', () => {
		render(QuestionFilters, { status: [] });
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
		expect(screen.getByLabelText('Sort')).toBeInTheDocument();
	});
});
