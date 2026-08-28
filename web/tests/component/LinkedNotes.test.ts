import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import LinkedNotes from '$lib/components/LinkedNotes.svelte';

describe('LinkedNotes', () => {
	it('renders source-note links', () => {
		render(LinkedNotes, {
			notes: [
				{ id: 'n1', title: 'First', displayMode: 'collapsed' },
				{ id: 'n2', title: 'Second', displayMode: 'link' }
			]
		});
		expect(screen.getByRole('link', { name: 'First' })).toHaveAttribute('href', '/notes/n1');
		expect(screen.getByRole('link', { name: 'Second' })).toBeInTheDocument();
	});
});
