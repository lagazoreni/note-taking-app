import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { get } from 'svelte/store';
import { api } from '$lib/api/client';
import WorkspacesPage from '../../src/routes/workspaces/+page.svelte';
import NewNotePage from '../../src/routes/notes/new/+page.svelte';
import {
	currentWorkspaceId,
	setWorkspaceList,
	workspaceHydrated,
	workspaces
} from '$lib/stores/workspace';
import type { Workspace } from '$lib/types/workspace';

const workspace: Workspace = {
	id: '11111111-1111-4111-8111-111111111111',
	name: 'Work',
	createdAt: '2026-08-09T00:00:00.000Z',
	updatedAt: '2026-08-09T00:00:00.000Z',
	version: 1
};

beforeEach(() => {
	localStorage.clear();
	workspaces.set([]);
	currentWorkspaceId.set(null);
	workspaceHydrated.set(true);
	vi.restoreAllMocks();
});

describe('workspace setup', () => {
	it('creates a workspace, selects it, and persists the returned selection', async () => {
		vi.spyOn(api, 'post').mockResolvedValue(workspace);

		render(WorkspacesPage);
		await fireEvent.input(screen.getByLabelText('New workspace'), {
			target: { value: 'Work' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Create' }));

		await waitFor(() => expect(screen.getByText('Work')).toBeInTheDocument());
		expect(get(currentWorkspaceId)).toBe(workspace.id);
		expect(localStorage.getItem('noted.currentWorkspace')).toBe(workspace.id);
		expect(screen.getByRole('button', { name: /Work Version 1/ })).toBeInTheDocument();
	});

	it('replaces a stale persisted selection with the first available workspace', () => {
		currentWorkspaceId.set('99999999-9999-4999-8999-999999999999');
		setWorkspaceList([workspace]);

		expect(get(currentWorkspaceId)).toBe(workspace.id);
		expect(localStorage.getItem('noted.currentWorkspace')).toBe(workspace.id);
	});

	it('exposes setup guidance and does not render note controls without a workspace', () => {
		setWorkspaceList([]);
		render(NewNotePage);

		expect(screen.getByText(/select or create a workspace first/i)).toBeInTheDocument();
		expect(screen.getByRole('link', { name: /go to workspaces/i })).toHaveAttribute(
			'href',
			'/workspaces'
		);
		expect(screen.queryByRole('button', { name: 'Save note' })).not.toBeInTheDocument();
	});
});
