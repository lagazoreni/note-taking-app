import { fireEvent, render, screen, waitFor } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import NoteEditor from '$lib/editor/NoteEditor.svelte';
import { notesApi } from '$lib/api/notes';
import { questionsApi } from '$lib/api/questions';
import { tagsApi } from '$lib/api/tags';
import type { Note } from '$lib/types/note';
import type { Tag } from '$lib/types/workspace';

const workspaceId = '11111111-1111-4111-8111-111111111111';
const noteId = '22222222-2222-4222-8222-222222222222';
const tag: Tag = {
	id: '33333333-3333-4333-8333-333333333333',
	ownerWorkspaceId: workspaceId,
	name: 'Important',
	availableWorkspaceIds: [workspaceId],
	createdAt: '2026-08-09T00:00:00.000Z',
	updatedAt: '2026-08-09T00:00:00.000Z',
	version: 1
};
const note: Note = {
	id: noteId,
	workspaceId,
	topicId: null,
	parentNoteId: null,
	title: 'Original title',
	bodyMarkdown: 'Original body',
	questionLinks: [],
	tagIds: [],
	createdAt: '2026-08-09T00:00:00.000Z',
	updatedAt: '2026-08-09T00:00:00.000Z',
	version: 4
};

beforeEach(() => {
	vi.restoreAllMocks();
	vi.spyOn(tagsApi, 'list').mockResolvedValue({ items: [tag] });
	vi.spyOn(questionsApi, 'list').mockResolvedValue({ items: [], nextCursor: null });
});

describe('NoteEditor', () => {
	it('renders title, body, and save controls', () => {
		render(NoteEditor, { workspaceId });
		expect(screen.getByLabelText('Title')).toBeInTheDocument();
		expect(screen.getByLabelText('Note')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Save note' })).toBeInTheDocument();
	});

	it('edits an existing note with bound fields and advances the save version', async () => {
		const firstSave = {
			...note,
			title: 'Edited title',
			bodyMarkdown: 'Original body\nMore text',
			version: 5
		};
		const secondSave = { ...firstSave, title: 'Edited again', version: 6 };
		const update = vi
			.spyOn(notesApi, 'update')
			.mockResolvedValueOnce(firstSave)
			.mockResolvedValueOnce(secondSave);
		render(NoteEditor, { workspaceId, existing: note });

		await fireEvent.click(screen.getByRole('button', { name: 'Edit note' }));
		await fireEvent.input(screen.getByLabelText('Title'), { target: { value: 'Edited title' } });
		await fireEvent.input(screen.getByLabelText('Note'), {
			target: { value: 'Original body\nMore text' }
		});
		await fireEvent.click(screen.getByRole('button', { name: 'Save note' }));
		await waitFor(() =>
			expect(update).toHaveBeenCalledWith(
				noteId,
				expect.objectContaining({
					version: 4,
					title: 'Edited title',
					bodyMarkdown: 'Original body\nMore text'
				})
			)
		);

		await fireEvent.input(screen.getByLabelText('Title'), { target: { value: 'Edited again' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Save note' }));
		await waitFor(() =>
			expect(update).toHaveBeenLastCalledWith(
				noteId,
				expect.objectContaining({ version: 5, title: 'Edited again' })
			)
		);
	});

	it('keeps existing-note editing available when the empty tag response is null', async () => {
		vi.mocked(tagsApi.list).mockResolvedValueOnce({
			items: null as unknown as Tag[]
		});
		render(NoteEditor, { workspaceId, existing: note });
		await fireEvent.click(screen.getByRole('button', { name: 'Edit note' }));
		await waitFor(() => expect(screen.getByLabelText('Title')).toBeInTheDocument());
	});

	it('persists selected workspace tags with an existing note', async () => {
		const update = vi
			.spyOn(notesApi, 'update')
			.mockResolvedValue({ ...note, tagIds: [tag.id], version: 5 });
		render(NoteEditor, { workspaceId, existing: note });
		await fireEvent.click(screen.getByRole('button', { name: 'Edit note' }));
		await waitFor(() => expect(screen.getByLabelText(tag.name)).toBeInTheDocument());
		await fireEvent.click(screen.getByLabelText(tag.name));
		await fireEvent.click(screen.getByRole('button', { name: 'Save note' }));
		await waitFor(() =>
			expect(update).toHaveBeenCalledWith(noteId, expect.objectContaining({ tagIds: [tag.id] }))
		);
	});
});
