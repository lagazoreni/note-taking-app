import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import TagSelector from '$lib/components/TagSelector.svelte';
import type { Tag } from '$lib/types/workspace';

const workspaceId = '11111111-1111-4111-8111-111111111111';
const tags: Tag[] = [
	{
		id: '22222222-2222-4222-8222-222222222222',
		ownerWorkspaceId: workspaceId,
		name: 'Work',
		availableWorkspaceIds: [workspaceId],
		createdAt: '2026-08-09T00:00:00.000Z',
		updatedAt: '2026-08-09T00:00:00.000Z',
		version: 1
	},
	{
		id: '33333333-3333-4333-8333-333333333333',
		ownerWorkspaceId: '44444444-4444-4444-8444-444444444444',
		name: 'Shared',
		availableWorkspaceIds: [workspaceId, '44444444-4444-4444-8444-444444444444'],
		createdAt: '2026-08-09T00:00:00.000Z',
		updatedAt: '2026-08-09T00:00:00.000Z',
		version: 1
	}
];

describe('TagSelector', () => {
	it('renders owned and shared tags and toggles selected values', async () => {
		const onChange = vi.fn();
		render(TagSelector, { tags, selected: [tags[0].id], onChange });
		expect(screen.getByRole('checkbox', { name: 'Work' })).toBeChecked();
		expect(screen.getByText('shared')).toBeInTheDocument();
		await fireEvent.click(screen.getByRole('checkbox', { name: 'Shared' }));
		expect(screen.getByRole('checkbox', { name: 'Shared' })).toBeChecked();
		expect(onChange).toHaveBeenCalledWith([tags[0].id, tags[1].id]);
	});

	it('shows an empty state when no tags are available', () => {
		render(TagSelector, { tags: [], selected: [] });
		expect(screen.getByText(/no workspace tags available/i)).toBeInTheDocument();
	});

	it('treats null list values from an older API response as empty arrays', () => {
		render(TagSelector, {
			tags: null as unknown as Tag[],
			selected: null as unknown as string[]
		});
		expect(screen.getByText(/no workspace tags available/i)).toBeInTheDocument();
	});
});
