import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import ExportData from '$lib/components/ExportData.svelte';

describe('DataPortability', () => {
	it('renders export action', () => {
		render(ExportData);
		expect(screen.getByRole('button', { name: /export/i })).toBeInTheDocument();
	});
});
