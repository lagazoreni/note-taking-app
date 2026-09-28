import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import QuestionContext from '$lib/components/QuestionContext.svelte';
import type { Note } from '$lib/types/note';
import type { Question } from '$lib/types/question';

const workspaceId = '11111111-1111-4111-8111-111111111111';
const firstNoteId = '22222222-2222-4222-8222-222222222222';
const secondNoteId = '33333333-3333-4333-8333-333333333333';

function makeQuestion(overrides: Partial<Question> = {}): Question {
	return {
		id: '550e8400-e29b-41d4-a716-446655440000',
		workspaceId,
		questionText: 'What causes this?',
		answerMarkdown: null,
		status: 'unanswered',
		priority: 'none',
		dueDate: null,
		reminder: null,
		tagIds: [],
		linkedNotes: [],
		createdAt: '2026-08-09T00:00:00.000Z',
		updatedAt: '2026-08-09T00:00:00.000Z',
		version: 1,
		kind: 'question',
		...overrides
	};
}

function makeNote(overrides: Partial<Note> = {}): Note {
	return {
		id: firstNoteId,
		workspaceId,
		topicId: null,
		parentNoteId: null,
		title: 'Research notes',
		bodyMarkdown: 'Hello',
		questionLinks: [],
		tagIds: [],
		createdAt: '2026-08-09T00:00:00.000Z',
		updatedAt: '2026-08-09T00:00:00.000Z',
		version: 1,
		...overrides
	};
}

describe('QuestionContext', () => {
	it('shows the linked note title and question lifecycle', () => {
		const note = makeNote();
		const question = makeQuestion({
			linkedNotes: [{ id: note.id, title: note.title, displayMode: 'collapsed' }]
		});
		render(QuestionContext, {
			question,
			note,
			selectedNoteId: note.id,
			onSelectNote: vi.fn(),
			onSave: vi.fn()
		});
		expect(screen.getByText('Research notes')).toBeInTheDocument();
		expect(screen.getByLabelText('Answer')).toBeInTheDocument();
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
	});

	it('shows currently unlinked copy and still shows question lifecycle', () => {
		render(QuestionContext, {
			question: makeQuestion(),
			note: null,
			selectedNoteId: '',
			onSelectNote: vi.fn(),
			onSave: vi.fn()
		});
		expect(screen.getByText(/currently unlinked/i)).toBeInTheDocument();
		expect(screen.getByLabelText('Answer')).toBeInTheDocument();
		expect(screen.getByLabelText('Status')).toBeInTheDocument();
	});

	it('renders a control to switch notes when two notes are linked', async () => {
		const note = makeNote({ title: 'First note' });
		const onSelectNote = vi.fn();
		const question = makeQuestion({
			linkedNotes: [
				{ id: firstNoteId, title: 'First note', displayMode: 'collapsed' },
				{ id: secondNoteId, title: 'Second note', displayMode: 'collapsed' }
			]
		});
		render(QuestionContext, {
			question,
			note,
			selectedNoteId: firstNoteId,
			onSelectNote,
			onSave: vi.fn()
		});
		const switcher = screen.getByLabelText('Source note');
		expect(switcher).toBeInTheDocument();
		expect(screen.getByRole('option', { name: 'First note' })).toBeInTheDocument();
		expect(screen.getByRole('option', { name: 'Second note' })).toBeInTheDocument();
		await fireEvent.change(switcher, { target: { value: secondNoteId } });
		expect(onSelectNote).toHaveBeenCalledWith(secondNoteId);
	});
});
