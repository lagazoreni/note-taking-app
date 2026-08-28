import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import QuestionFilters from '$lib/components/QuestionFilters.svelte';

describe('QuestionDiscovery', () => {
	it('renders filter and sort controls', () => {
		render(QuestionFilters, { status: [] });
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
		expect(screen.getByLabelText('Sort')).toBeInTheDocument();
	});

	it('uses a compact accessible status checkbox filter and preserves selections', async () => {
		const onChange = vi.fn();
		render(QuestionFilters, { status: ['unanswered'], onChange });
		await fireEvent.click(screen.getByRole('button', { name: 'Status' }));
		expect(screen.queryByRole('listbox')).not.toBeInTheDocument();
		expect(screen.getByRole('checkbox', { name: 'Unanswered' })).toBeChecked();
		await fireEvent.click(screen.getByRole('checkbox', { name: 'In progress' }));
		expect(onChange).toHaveBeenLastCalledWith({
			status: ['unanswered', 'in_progress'],
			sort: 'updatedAt',
			direction: 'desc'
		});
	});
});
