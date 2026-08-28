import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import NoteEditor from '$lib/editor/NoteEditor.svelte';

describe('NoteEditor', () => {
	it('renders title, body, and save controls', () => {
		render(NoteEditor, { workspaceId: '11111111-1111-4111-8111-111111111111' });
		expect(screen.getByLabelText('Title')).toBeInTheDocument();
		expect(screen.getByLabelText('Note')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Save note' })).toBeInTheDocument();
	});
});
