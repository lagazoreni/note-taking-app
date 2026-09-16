import { fireEvent, render, screen, within } from '@testing-library/svelte';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import NoteReader from '$lib/editor/NoteReader.svelte';
import type { Question } from '$lib/types/question';

const id = '550e8400-e29b-41d4-a716-446655440000';

const question: Question = {
	id,
	workspaceId: '11111111-1111-4111-8111-111111111111',
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
	kind: 'question'
};

const annotation: Question = {
	...question,
	id: '650e8400-e29b-41d4-a716-446655440000',
	questionText: 'Existing annotation',
	kind: 'annotation'
};

type SelectionBounds = {
	top: number;
	left: number;
	right: number;
	bottom: number;
	width: number;
	height: number;
};

function mockSelection(
	text: string,
	bounds: SelectionBounds = { top: 120, left: 80, right: 260, bottom: 145, width: 180, height: 25 }
) {
	vi.spyOn(window, 'getSelection').mockReturnValue({
		toString: () => text,
		rangeCount: text ? 1 : 0,
		isCollapsed: !text,
		removeAllRanges: () => undefined,
		addRange: () => undefined,
		getRangeAt: () => ({
			collapse: () => undefined,
			getBoundingClientRect: () => bounds
		})
	} as unknown as Selection);
}

function expectFloating(element: HTMLElement) {
	expect(['absolute', 'fixed']).toContain(element.style.position);
	expect(element.style.top).not.toBe('');
	expect(element.style.left).not.toBe('');
}

function expectInsideViewport(element: HTMLElement) {
	expectFloating(element);
	const top = Number.parseFloat(element.style.top);
	const left = Number.parseFloat(element.style.left);
	expect(Number.isFinite(top)).toBe(true);
	expect(Number.isFinite(left)).toBe(true);
	expect(top).toBeGreaterThanOrEqual(0);
	expect(left).toBeGreaterThanOrEqual(0);
	expect(top).toBeLessThanOrEqual(window.innerHeight);
	expect(left).toBeLessThanOrEqual(window.innerWidth);
}

async function selectPassageWithPointer(article: HTMLElement, passage: string) {
	mockSelection(passage);
	await fireEvent.mouseDown(article, { button: 0 });
	await fireEvent.mouseUp(article, { button: 0 });
}

async function clickWithPointer(element: HTMLElement) {
	await fireEvent.mouseDown(element, { button: 0 });
	await fireEvent.mouseUp(element, { button: 0 });
	await fireEvent.click(element);
}

async function focusAndFillWithPointer(
	control: HTMLInputElement | HTMLTextAreaElement,
	value: string
) {
	await fireEvent.mouseDown(control, { button: 0 });
	// jsdom does not perform the browser's default focus action for fireEvent.
	// Keep the pointer sequence, then model that native action before mouseup/click.
	control.focus();
	await fireEvent.input(control, { target: { value } });
	await fireEvent.mouseUp(control, { button: 0 });
	await fireEvent.click(control);
}

describe('NoteReader', () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it('shows selection menu actions after text is selected', async () => {
		render(NoteReader, { markdown: 'The mitochondria is the powerhouse.', questions: [] });
		mockSelection('The mitochondria is the powerhouse.');
		await fireEvent.mouseUp(screen.getByRole('article', { name: 'Reading note' }));
		expect(screen.getByRole('button', { name: 'Ask a question' })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Add annotation' })).toBeInTheDocument();
		expect(screen.getByText('Q ask · A annotate · L link · Esc cancel')).toBeInTheDocument();
	});

	it('positions the selection toolbar next to the selected passage', async () => {
		render(NoteReader, { markdown: 'The mitochondria is the powerhouse.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });
		mockSelection('The mitochondria is the powerhouse.', {
			top: 180,
			left: 96,
			right: 300,
			bottom: 208,
			width: 204,
			height: 28
		});

		await fireEvent.mouseUp(article);

		const toolbar = screen.getByRole('toolbar', { name: 'Selection actions' });
		expectFloating(toolbar);
		const top = Number.parseFloat(toolbar.style.top);
		const left = Number.parseFloat(toolbar.style.left);
		if (toolbar.dataset.placement === 'above') {
			expect(top).toBeLessThanOrEqual(180 - 8);
		} else {
			expect(top).toBeGreaterThanOrEqual(208 + 8);
		}
		expect(Math.abs(left - 96)).toBeLessThanOrEqual(48);
	});

	it('keeps the composer anchored to the selected passage', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });
		mockSelection('Select this passage.', {
			top: 120,
			left: 80,
			right: 220,
			bottom: 148,
			width: 140,
			height: 28
		});

		await fireEvent.mouseUp(article);
		const toolbar = screen.getByRole('toolbar', { name: 'Selection actions' });
		const toolbarTop = Number.parseFloat(toolbar.style.top);
		const toolbarLeft = Number.parseFloat(toolbar.style.left);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));

		const composer = document.querySelector('.composer');
		expect(composer).toBeInstanceOf(HTMLElement);
		const composerElement = composer as HTMLElement;
		expectFloating(composerElement);
		expect(Number.parseFloat(composerElement.style.top)).toBeCloseTo(toolbarTop, 0);
		expect(Number.parseFloat(composerElement.style.left)).toBeCloseTo(toolbarLeft, 0);
	});

	it('clamps the toolbar and composer inside the viewport near its edges', async () => {
		const originalWidth = window.innerWidth;
		const originalHeight = window.innerHeight;
		Object.defineProperty(window, 'innerWidth', { configurable: true, value: 320 });
		Object.defineProperty(window, 'innerHeight', { configurable: true, value: 240 });

		try {
			render(NoteReader, { markdown: 'Select this passage.', questions: [] });
			const article = screen.getByRole('article', { name: 'Reading note' });
			mockSelection('Select this passage.', {
				top: 220,
				left: 300,
				right: 319,
				bottom: 239,
				width: 19,
				height: 19
			});

			await fireEvent.mouseUp(article);
			const toolbar = screen.getByRole('toolbar', { name: 'Selection actions' });
			expectInsideViewport(toolbar);

			await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));
			const composer = document.querySelector('.composer') as HTMLElement;
			expectInsideViewport(composer);
		} finally {
			Object.defineProperty(window, 'innerWidth', { configurable: true, value: originalWidth });
			Object.defineProperty(window, 'innerHeight', { configurable: true, value: originalHeight });
		}
	});

	it('opens a card when a highlight is clicked', async () => {
		render(NoteReader, {
			markdown: `Lead {{question:${id}}}selected passage{{/question}} tail`,
			questions: [question]
		});
		expect(screen.queryByText('{{question:')).not.toBeInTheDocument();
		await fireEvent.click(screen.getByText('selected passage'));
		expect(screen.getByRole('dialog', { name: /question/i })).toBeInTheDocument();
		expect(screen.getByText('What causes this?')).toBeInTheDocument();
	});

	it('keeps a positioned question card open while editing its answer', async () => {
		render(NoteReader, {
			markdown: `Lead {{question:${id}}}selected passage{{/question}} tail`,
			questions: [question]
		});
		const mark = document.querySelector('mark[data-annotation-id]');
		expect(mark).toBeInstanceOf(HTMLElement);
		Object.defineProperty(mark, 'getBoundingClientRect', {
			configurable: true,
			value: () => ({
				top: 120,
				left: 80,
				right: 220,
				bottom: 145,
				width: 140,
				height: 25
			})
		});

		await clickWithPointer(mark as HTMLElement);
		const card = screen.getByRole('dialog', { name: /question/i });
		expectFloating(card);

		const answer = within(card).getByLabelText('Answer') as HTMLTextAreaElement;
		const enteredAnswer = 'Because the surrounding conditions changed.';
		await focusAndFillWithPointer(answer, enteredAnswer);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));

		expect(document.activeElement).toBe(answer);
		expect(answer).toHaveValue(enteredAnswer);
		expect(screen.getByRole('dialog', { name: /question/i })).toBeInTheDocument();

		await clickWithPointer(document.body);
		expect(screen.queryByRole('dialog', { name: /question/i })).not.toBeInTheDocument();
	});

	it('keeps a positioned annotation card open while editing its annotation', async () => {
		render(NoteReader, {
			markdown: `Lead {{question:${annotation.id}}}annotated passage{{/question}} tail`,
			questions: [annotation]
		});
		const mark = document.querySelector('mark[data-annotation-id]');
		expect(mark).toBeInstanceOf(HTMLElement);
		Object.defineProperty(mark, 'getBoundingClientRect', {
			configurable: true,
			value: () => ({
				top: 160,
				left: 100,
				right: 250,
				bottom: 185,
				width: 150,
				height: 25
			})
		});

		await clickWithPointer(mark as HTMLElement);
		const card = screen.getByRole('dialog', { name: /annotation/i });
		expectFloating(card);

		const annotationField = within(card).getByLabelText('Annotation') as HTMLTextAreaElement;
		const enteredAnnotation = 'Keep this example for the next review.';
		await focusAndFillWithPointer(annotationField, enteredAnnotation);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));

		expect(document.activeElement).toBe(annotationField);
		expect(annotationField).toHaveValue(enteredAnnotation);
		expect(screen.getByRole('dialog', { name: /annotation/i })).toBeInTheDocument();

		await clickWithPointer(document.body);
		expect(screen.queryByRole('dialog', { name: /annotation/i })).not.toBeInTheDocument();
	});

	it('opens the existing-question picker and links the selected passage', async () => {
		const onLinkExisting = vi.fn().mockResolvedValue(question);
		render(NoteReader, {
			markdown: 'Select this passage.',
			questions: [],
			linkableQuestions: [question],
			onLinkExisting
		});
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Link existing question' }));

		expect(screen.getByRole('region', { name: 'Find an existing question' })).toBeInTheDocument();
		await fireEvent.click(screen.getByRole('button', { name: /What causes this\?/ }));

		expect(onLinkExisting).toHaveBeenCalledOnce();
		expect(onLinkExisting).toHaveBeenCalledWith('Select this passage.', question);
	});

	it('opens the question composer with Q for a selected passage', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.keyDown(window, { key: 'q' });

		expect(screen.getByLabelText('Question text')).toBeInTheDocument();
	});

	it('opens the annotation composer with A for a selected passage', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.keyDown(window, { key: 'a' });

		expect(screen.getByLabelText('Annotation')).toBeInTheDocument();
	});

	it('opens the existing-question picker with L when linking is available', async () => {
		const onLinkExisting = vi.fn().mockResolvedValue(question);
		render(NoteReader, {
			markdown: 'Select this passage.',
			questions: [],
			linkableQuestions: [question],
			onLinkExisting
		});
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.keyDown(window, { key: 'l' });

		expect(screen.getByRole('region', { name: 'Find an existing question' })).toBeInTheDocument();
	});

	it('dismisses the selection toolbar with Escape', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		expect(screen.getByRole('toolbar', { name: 'Selection actions' })).toBeInTheDocument();

		await fireEvent.keyDown(window, { key: 'Escape' });

		expect(screen.queryByRole('toolbar', { name: 'Selection actions' })).not.toBeInTheDocument();
	});

	it('ignores Q while the composer is open', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Add annotation' }));
		expect(screen.getByLabelText('Annotation')).toBeInTheDocument();

		await fireEvent.keyDown(window, { key: 'q' });

		expect(screen.getByLabelText('Annotation')).toBeInTheDocument();
		expect(screen.queryByLabelText('Question text')).not.toBeInTheDocument();
	});

	it('ignores Q when focus is in an input or textarea', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);

		const input = document.createElement('input');
		const textarea = document.createElement('textarea');
		document.body.append(input, textarea);
		try {
			await fireEvent.keyDown(input, { key: 'q' });
			await fireEvent.keyDown(textarea, { key: 'q' });
		} finally {
			input.remove();
			textarea.remove();
		}

		expect(screen.getByRole('toolbar', { name: 'Selection actions' })).toBeInTheDocument();
		expect(screen.queryByLabelText('Question text')).not.toBeInTheDocument();
		expect(screen.queryByLabelText('Annotation')).not.toBeInTheDocument();
	});

	it('opens the matching highlight card when focusQuestionId is set', () => {
		render(NoteReader, {
			markdown: `Lead {{question:${id}}}selected passage{{/question}} tail`,
			questions: [question],
			focusQuestionId: id
		});

		const dialog = screen.getByRole('dialog', { name: /question/i });
		expect(dialog).toBeInTheDocument();
		expect(within(dialog).getByText('selected passage')).toBeInTheDocument();
		expect(within(dialog).getByText('What causes this?')).toBeInTheDocument();
	});

	it('keeps the question composer focused through internal pointer and selection events', async () => {
		const onCapture = vi.fn().mockResolvedValue(undefined);
		render(NoteReader, { markdown: 'Select this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });
		const questionText = 'Why does this happen?';

		await selectPassageWithPointer(article, 'Select this passage.');
		await clickWithPointer(screen.getByRole('button', { name: 'Ask a question' }));

		const input = screen.getByLabelText('Question text') as HTMLInputElement;
		await focusAndFillWithPointer(input, questionText);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));

		expect(document.activeElement).toBe(input);
		expect(input).toHaveValue(questionText);
		expect(document.querySelector('.composer')).toBeInTheDocument();

		await clickWithPointer(screen.getByRole('button', { name: 'Save question' }));
		expect(onCapture).toHaveBeenCalledOnce();
		expect(onCapture).toHaveBeenCalledWith('question', 'Select this passage.', questionText);
	});

	it('keeps the annotation composer focused through internal pointer and selection events', async () => {
		const onCapture = vi.fn().mockResolvedValue(undefined);
		render(NoteReader, { markdown: 'Annotate this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });
		const annotationText = 'Remember this context.';

		await selectPassageWithPointer(article, 'Annotate this passage.');
		await clickWithPointer(screen.getByRole('button', { name: 'Add annotation' }));

		const textarea = screen.getByLabelText('Annotation') as HTMLTextAreaElement;
		await focusAndFillWithPointer(textarea, annotationText);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));

		expect(document.activeElement).toBe(textarea);
		expect(textarea).toHaveValue(annotationText);
		expect(document.querySelector('.composer')).toBeInTheDocument();

		await clickWithPointer(screen.getByRole('button', { name: 'Save annotation' }));
		expect(onCapture).toHaveBeenCalledOnce();
		expect(onCapture).toHaveBeenCalledWith('annotation', 'Annotate this passage.', annotationText);
	});

	it('shows capture errors and keeps the composer open', async () => {
		const errorMessage = 'Could not find that passage in the note. Try selecting plain text.';
		const onCapture = vi.fn().mockRejectedValue(new Error(errorMessage));
		render(NoteReader, { markdown: 'Select this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));
		const input = screen.getByLabelText('Question text');
		await fireEvent.input(input, { target: { value: 'Why does this happen?' } });
		await fireEvent.click(screen.getByRole('button', { name: 'Save question' }));

		expect(await screen.findByRole('alert')).toHaveTextContent(errorMessage);
		expect(screen.getByLabelText('Question text')).toBeInTheDocument();
		expect(screen.queryByRole('dialog', { name: /question/i })).not.toBeInTheDocument();
		expect(onCapture).toHaveBeenCalledOnce();
	});

	it('dismisses selection actions and the composer without capturing', async () => {
		const onCapture = vi.fn();
		render(NoteReader, { markdown: 'Select this passage.', questions: [], onCapture });
		const article = screen.getByRole('article', { name: 'Reading note' });

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));
		expect(screen.queryByRole('toolbar', { name: 'Selection actions' })).not.toBeInTheDocument();

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));
		expect(screen.getByLabelText('Question text')).toBeInTheDocument();
		await fireEvent.keyDown(window, { key: 'Escape' });
		expect(screen.queryByLabelText('Question text')).not.toBeInTheDocument();
		expect(onCapture).not.toHaveBeenCalled();

		mockSelection('Select this passage.');
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Add annotation' }));
		await fireEvent.click(document.body);
		expect(screen.queryByLabelText('Annotation')).not.toBeInTheDocument();
		expect(onCapture).not.toHaveBeenCalled();
	});

	it('keeps the toolbar outside the selected line with clearance', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });
		const bounds = {
			top: 430,
			left: 80,
			right: 280,
			bottom: 455,
			width: 200,
			height: 25
		};
		const originalHeight = window.innerHeight;
		Object.defineProperty(window, 'innerHeight', { configurable: true, value: 480 });

		try {
			mockSelection('Select this passage.', bounds);
			await fireEvent.mouseUp(article);

			const toolbar = screen.getByRole('toolbar', { name: 'Selection actions' });
			const toolbarTop = Number.parseFloat(toolbar.style.top);
			const estimatedToolbarHeight = 96;
			if (toolbar.dataset.placement === 'above') {
				expect(toolbarTop + estimatedToolbarHeight).toBeLessThanOrEqual(bounds.top - 8);
			} else {
				expect(toolbarTop).toBeGreaterThanOrEqual(bounds.bottom + 8);
			}
		} finally {
			Object.defineProperty(window, 'innerHeight', {
				configurable: true,
				value: originalHeight
			});
		}
	});

	it('waits for mouseup before showing a toolbar during drag selection', async () => {
		render(NoteReader, { markdown: 'Select this passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });
		mockSelection('Select this passage.');

		await fireEvent.mouseDown(article);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));
		expect(screen.queryByRole('toolbar', { name: 'Selection actions' })).not.toBeInTheDocument();

		await fireEvent.mouseUp(article);
		expect(screen.getByRole('toolbar', { name: 'Selection actions' })).toBeInTheDocument();
	});

	it('keeps selection actions while clicking or moving focus inside the reader', async () => {
		render(NoteReader, { markdown: 'First passage. Second passage.', questions: [] });
		const article = screen.getByRole('article', { name: 'Reading note' });
		mockSelection('First passage.', {
			top: 120,
			left: 80,
			right: 190,
			bottom: 145,
			width: 110,
			height: 25
		});
		await fireEvent.mouseUp(article);

		await fireEvent.focus(article);
		await fireEvent.click(article);
		expect(screen.getByRole('toolbar', { name: 'Selection actions' })).toBeInTheDocument();

		mockSelection('Second passage.', {
			top: 180,
			left: 80,
			right: 210,
			bottom: 205,
			width: 130,
			height: 25
		});
		await fireEvent.mouseDown(article);
		document.dispatchEvent(new Event('selectionchange', { bubbles: true }));
		await fireEvent.mouseUp(article);
		await fireEvent.click(screen.getByRole('button', { name: 'Ask a question' }));

		expect(screen.getByRole('blockquote')).toHaveTextContent('Second passage.');
	});

	it('places a top-edge selection below the passage without clipping it', async () => {
		const originalWidth = window.innerWidth;
		const originalHeight = window.innerHeight;
		Object.defineProperty(window, 'innerWidth', { configurable: true, value: 320 });
		Object.defineProperty(window, 'innerHeight', { configurable: true, value: 240 });

		try {
			render(NoteReader, { markdown: 'Select this passage.', questions: [] });
			const article = screen.getByRole('article', { name: 'Reading note' });
			const bounds = {
				top: 2,
				left: 24,
				right: 210,
				bottom: 24,
				width: 186,
				height: 22
			};
			mockSelection('Select this passage.', bounds);
			await fireEvent.mouseUp(article);

			const toolbar = screen.getByRole('toolbar', { name: 'Selection actions' });
			const toolbarTop = Number.parseFloat(toolbar.style.top);
			expect(toolbar.dataset.placement).toBe('below');
			expect(toolbarTop).toBeGreaterThanOrEqual(bounds.bottom + 8);
			expect(toolbarTop + 96).toBeLessThanOrEqual(window.innerHeight - 8);
		} finally {
			Object.defineProperty(window, 'innerWidth', { configurable: true, value: originalWidth });
			Object.defineProperty(window, 'innerHeight', { configurable: true, value: originalHeight });
		}
	});
});
