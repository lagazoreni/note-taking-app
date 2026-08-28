import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import WorkspaceSwitcher from '$lib/components/WorkspaceSwitcher.svelte';

describe('Organization', () => {
	it('has an accessible workspace selector', () => {
		render(WorkspaceSwitcher);
		expect(screen.getByRole('combobox', { name: /current workspace/i })).toBeInTheDocument();
	});
});
