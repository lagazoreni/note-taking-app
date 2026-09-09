<script lang="ts">
	import AnnotationCard from '$lib/components/AnnotationCard.svelte';
	import QuestionPicker from '$lib/components/QuestionPicker.svelte';
	import { renderNoteHtml } from './markdown';
	import { tokenizeDirectives } from './directives';
	import type { Question, QuestionKind } from '$lib/types/question';

	export let markdown: string;
	export let questions: Question[] = [];
	export let focusQuestionId: string | null = null;
	export let onCapture:
		((kind: QuestionKind, passage: string, text: string) => Promise<Question | void>) | undefined =
		undefined;
	export let onInsertAnswer: ((question: Question) => Promise<void> | void) | undefined = undefined;
	export let onLinkExisting:
		((passage: string, question: Question) => Promise<Question | void>) | undefined = undefined;
	export let linkableQuestions: Question[] = [];

	let showToolbar = false;
	let selectedPassage = '';
	let composerOpen = false;
	let composerKind: QuestionKind = 'question';
	let composerText = '';
	let composerSaving = false;
	let composerError = '';
	let pickerOpen = false;
	let pickerSaving = false;
	let pickerError = '';
	let openQuestion: Question | null = null;
	let openPassage = '';
	let actionLayer: HTMLDivElement | null = null;

	$: html = renderNoteHtml(markdown, questions);

	function selectedText(): string {
		return window.getSelection()?.toString().trim() ?? '';
	}

	function clearSelection() {
		showToolbar = false;
		selectedPassage = '';
	}

	function cancelComposer() {
		if (composerSaving) return;
		composerOpen = false;
		composerText = '';
		composerError = '';
		clearSelection();
	}

	function cancelPicker() {
		if (pickerSaving) return;
		pickerOpen = false;
		pickerError = '';
		clearSelection();
	}

	function onMouseUp() {
		const text = selectedText();
		if (!text) {
			clearSelection();
			return;
		}
		selectedPassage = text;
		showToolbar = true;
	}

	function onContextMenu(event: MouseEvent) {
		const text = selectedText();
		if (!text) return;
		event.preventDefault();
		selectedPassage = text;
		showToolbar = true;
	}

	function onClick(event: MouseEvent) {
		const target = event.target as HTMLElement | null;
		const mark = target?.closest?.('mark[data-annotation-id]');
		if (!mark) return;
		const id = mark.getAttribute('data-annotation-id');
		if (!id) return;
		const question = questions.find((item) => item.id === id) ?? null;
		if (!question) return;
		const token = tokenizeDirectives(markdown).find(
			(item) => item.kind === 'question' && item.directive?.id === id
		);
		openPassage = token?.directive?.snippet || mark.textContent || '';
		openQuestion = question;
		clearSelection();
	}

	$: if (focusQuestionId !== null) {
		const question = questions.find((item) => item.id === focusQuestionId) ?? null;
		if (question) {
			const token = tokenizeDirectives(markdown).find(
				(item) => item.kind === 'question' && item.directive?.id === focusQuestionId
			);
			openPassage = token?.directive?.snippet || '';
			openQuestion = question;
			clearSelection();
		}
	}

	function openComposer(kind: QuestionKind) {
		composerKind = kind;
		composerText = '';
		composerError = '';
		composerOpen = true;
		showToolbar = false;
	}

	function openLinkPicker() {
		if (!onLinkExisting) return;
		pickerError = '';
		pickerOpen = true;
		showToolbar = false;
	}

	async function selectExisting(question: Question) {
		if (!onLinkExisting || pickerSaving || !selectedPassage) return;
		const passage = selectedPassage;
		pickerSaving = true;
		pickerError = '';
		try {
			await onLinkExisting(passage, question);
			pickerOpen = false;
			clearSelection();
		} catch (cause) {
			pickerError = cause instanceof Error ? cause.message : 'Could not link question';
		} finally {
			pickerSaving = false;
		}
	}

	function dismissFromOutside(event: MouseEvent) {
		if (!showToolbar && !composerOpen && !pickerOpen) return;
		const target = event.target as Node | null;
		if (actionLayer && target && actionLayer.contains(target)) return;

		// The action that opens the composer removes the toolbar and adds the composer
		// during the same click. By the time the window click handler runs, the
		// bind:this value can briefly be null, even though the click originated in
		// the toolbar. Inspect the event path as well so that action clicks are not
		// mistaken for outside clicks during that transition.
		const cameFromActionLayer = event
			.composedPath()
			.some(
				(node) =>
					node instanceof HTMLElement &&
					(node.classList.contains('toolbar') ||
						node.classList.contains('composer') ||
						node.classList.contains('picker-layer'))
			);
		if (cameFromActionLayer) return;
		if (composerOpen) {
			if (!composerSaving) cancelComposer();
			return;
		}
		if (pickerOpen) {
			if (!pickerSaving) cancelPicker();
			return;
		}
		clearSelection();
	}

	function onWindowClick(event: MouseEvent) {
		dismissFromOutside(event);
	}

	function onWindowMouseDown(event: MouseEvent) {
		dismissFromOutside(event);
	}

	function onWindowKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			if (composerOpen) {
				if (!composerSaving) cancelComposer();
				return;
			}
			if (pickerOpen) {
				if (!pickerSaving) cancelPicker();
				return;
			}
			if (showToolbar) {
				event.preventDefault();
				clearSelection();
			}
			return;
		}

		if (event.defaultPrevented) return;
		const target = event.target as Element | null;
		if (target?.closest?.('input, textarea, select, [contenteditable]')) return;
		if (composerOpen) return;
		if (event.ctrlKey || event.metaKey || event.altKey) return;
		if (!selectedPassage && !showToolbar) return;

		const key = event.key.toLowerCase();
		if (key === 'q') {
			event.preventDefault();
			openComposer('question');
		} else if (key === 'a') {
			event.preventDefault();
			openComposer('annotation');
		} else if (key === 'l' && onLinkExisting) {
			event.preventDefault();
			openLinkPicker();
		}
	}

	async function submitComposer() {
		if (!selectedPassage || !composerText.trim()) return;
		composerSaving = true;
		composerError = '';
		try {
			const captured = await onCapture?.(composerKind, selectedPassage, composerText.trim());
			if (captured) {
				openPassage = selectedPassage;
				openQuestion = captured;
			}
			composerOpen = false;
			composerText = '';
			selectedPassage = '';
		} catch (cause) {
			composerError = cause instanceof Error ? cause.message : 'Could not save highlight';
		} finally {
			composerSaving = false;
		}
	}
</script>

<svelte:window
	onclick={onWindowClick}
	onmousedown={onWindowMouseDown}
	onkeydown={onWindowKeydown}
/>

<div class="reader">
	<!-- eslint-disable-next-line svelte/no-unused-svelte-ignore -->
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -- Text selection and highlight clicks are intentionally handled on the reading surface. -->
	<article
		aria-label="Reading note"
		onmouseup={onMouseUp}
		oncontextmenu={onContextMenu}
		onclick={onClick}
	>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -- renderNoteHtml sanitizes the generated HTML with DOMPurify. -->
		{@html html}
	</article>
	{#if showToolbar}
		<div class="toolbar" role="toolbar" aria-label="Selection actions" bind:this={actionLayer}>
			<button type="button" onclick={() => openComposer('question')}>Ask a question</button>
			<button type="button" onclick={() => openComposer('annotation')}>Add annotation</button>
			{#if onLinkExisting}
				<button type="button" onclick={openLinkPicker}>Link existing question</button>
			{/if}
			<span class="keyboard-hint">Q ask · A annotate · L link · Esc cancel</span>
			<button type="button" class="secondary" onclick={clearSelection}>Cancel</button>
		</div>
	{/if}
	{#if pickerOpen}
		<div class="picker-layer" bind:this={actionLayer}>
			<QuestionPicker questions={linkableQuestions} onSelect={selectExisting} />
			{#if pickerError}<div class="error" role="alert">{pickerError}</div>{/if}
			<button type="button" class="secondary" onclick={cancelPicker} disabled={pickerSaving}>
				Cancel
			</button>
		</div>
	{/if}
	{#if composerOpen}
		<div class="composer" bind:this={actionLayer}>
			<blockquote>{selectedPassage}</blockquote>
			{#if composerKind === 'question'}
				<label for="reader-question-text">Question text</label>
				<input id="reader-question-text" bind:value={composerText} />
				<div class="composer-actions">
					<button
						type="button"
						onclick={submitComposer}
						disabled={composerSaving || !composerText.trim()}
						>{composerSaving ? 'Saving…' : 'Save question'}</button
					>
					<button
						type="button"
						class="secondary"
						onclick={cancelComposer}
						disabled={composerSaving}
					>
						Cancel
					</button>
				</div>
			{:else}
				<label for="reader-annotation-text">Annotation</label>
				<textarea id="reader-annotation-text" rows="4" bind:value={composerText}></textarea>
				<div class="composer-actions">
					<button
						type="button"
						onclick={submitComposer}
						disabled={composerSaving || !composerText.trim()}
						>{composerSaving ? 'Saving…' : 'Save annotation'}</button
					>
					<button
						type="button"
						class="secondary"
						onclick={cancelComposer}
						disabled={composerSaving}
					>
						Cancel
					</button>
				</div>
			{/if}
			{#if composerError}<div class="error" role="alert">{composerError}</div>{/if}
		</div>
	{/if}
	{#if openQuestion}
		<AnnotationCard
			question={openQuestion}
			passage={openPassage}
			onClose={() => (openQuestion = null)}
			onInsertAnswer={() => onInsertAnswer?.(openQuestion!)}
		/>
	{/if}
</div>

<style>
	.reader {
		display: grid;
		gap: 0.75rem;
	}
	article {
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.65rem;
		padding: 1.1rem 1.2rem;
		line-height: 1.6;
		color: #0f172a;
	}
	article :global(mark[data-annotation-id]) {
		background: #fef3c7;
		border-bottom: 2px solid #d97706;
		cursor: pointer;
		padding: 0 0.1em;
	}
	article :global(mark[data-status="answered"]) {
		background: #d1fae5;
		border-bottom: 2px solid #047857;
	}
	article :global(mark.annotation-chip) {
		display: inline-block;
		width: 0.65rem;
		height: 0.65rem;
		border-radius: 999px;
		padding: 0;
		vertical-align: middle;
	}
	.toolbar,
	.composer,
	.picker-layer {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-items: end;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
		padding: 0.75rem;
	}
	.composer,
	.picker-layer {
		display: grid;
	}
	.composer blockquote {
		margin: 0;
		padding: 0.5rem 0.65rem;
		background: #fef3c7;
		color: #78350f;
	}
	.composer-actions {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	label {
		font-weight: 650;
	}
	input,
	textarea {
		width: 100%;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	button {
		width: max-content;
		border: 0;
		background: #2563eb;
		color: #fff;
		border-radius: 0.4rem;
		padding: 0.55rem 0.8rem;
	}
	button.secondary {
		background: #fff;
		border: 1px solid #94a3b8;
		color: #334155;
	}
	button:disabled {
		opacity: 0.55;
	}
	.error {
		padding: 0.6rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.4rem;
	}
</style>
