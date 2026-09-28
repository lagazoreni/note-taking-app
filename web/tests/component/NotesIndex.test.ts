import { render, screen, waitFor } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { notesApi } from '$lib/api/notes';
import NotesPage from '../../src/routes/notes/+page.svelte';
import {
	currentWorkspaceId,
	setWorkspaceList,
	workspaceHydrated,
	workspaces
} from '$lib/stores/workspace';
import type { Workspace } from '$lib/types/workspace';
import type { Note } from '$lib/types/note';

const workspace: Workspace = {
	id: '11111111-1111-4111-8111-111111111111',
	name: 'Work',
	createdAt: '2026-08-09T00:00:00.000Z',
	updatedAt: '2026-08-09T00:00:00.000Z',
	version: 1
};

const note: Note = {
	id: '22222222-2222-4222-8222-222222222222',
	workspaceId: workspace.id,
	topicId: null,
	parentNoteId: null,
	title: 'Research notes',
	bodyMarkdown: 'Hello',
	questionLinks: [],
	tagIds: [],
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

describe('Notes index', () => {
	it('loads workspace notes and links to a note', async () => {
		setWorkspaceList([workspace]);
		vi.spyOn(notesApi, 'list').mockResolvedValue({ items: [note] });
		render(NotesPage);
		await waitFor(() => expect(screen.getByText('Research notes')).toBeInTheDocument());
		expect(notesApi.list).toHaveBeenCalledWith(workspace.id, expect.anything());
		expect(screen.getByRole('link', { name: 'Research notes' })).toHaveAttribute(
			'href',
			`/notes/${note.id}`
		);
	});

	it('shows an empty state when the workspace has no notes', async () => {
		setWorkspaceList([workspace]);
		vi.spyOn(notesApi, 'list').mockResolvedValue({ items: [] });
		render(NotesPage);
		await waitFor(() =>
			expect(screen.getByText(/your notes will appear here/i)).toBeInTheDocument()
		);
	});

	it('shows an error state when notes fail to load', async () => {
		setWorkspaceList([workspace]);
		vi.spyOn(notesApi, 'list').mockRejectedValue(new Error('offline'));
		render(NotesPage);
		await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent(/offline/i));
	});

	it('shows a workspace CTA when none is selected', () => {
		setWorkspaceList([]);
		render(NotesPage);
		expect(screen.getByText(/create a workspace/i)).toBeInTheDocument();
		expect(screen.getByRole('link', { name: /go to workspaces/i })).toHaveAttribute(
			'href',
			'/workspaces'
		);
	});
});
