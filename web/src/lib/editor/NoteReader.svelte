<script lang="ts">
	import AnnotationCard from '$lib/components/AnnotationCard.svelte';
	import { renderNoteHtml } from './markdown';
	import { tokenizeDirectives } from './directives';
	import type { Question, QuestionKind } from '$lib/types/question';

	export let markdown: string;
	export let questions: Question[] = [];
	export let onCapture:
		((kind: QuestionKind, passage: string, text: string) => Promise<Question | void>) | undefined =
		undefined;

	let showToolbar = false;
	let selectedPassage = '';
	let composerOpen = false;
	let composerKind: QuestionKind = 'question';
	let composerText = '';
	let composerSaving = false;
	let composerError = '';
	let openQuestion: Question | null = null;
	let openPassage = '';

	$: html = renderNoteHtml(markdown);

	function selectedText(): string {
		return window.getSelection()?.toString().trim() ?? '';
	}

	function onMouseUp() {
		const text = selectedText();
		if (!text) {
			showToolbar = false;
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
		showToolbar = false;
	}

	function openComposer(kind: QuestionKind) {
		composerKind = kind;
		composerText = '';
		composerError = '';
		composerOpen = true;
		showToolbar = false;
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
		} catch (cause) {
			composerError = cause instanceof Error ? cause.message : 'Could not save highlight';
		} finally {
			composerSaving = false;
		}
	}
</script>

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
		<div class="toolbar" role="toolbar" aria-label="Selection actions">
			<button type="button" onclick={() => openComposer('question')}>Ask a question</button>
			<button type="button" onclick={() => openComposer('annotation')}>Add annotation</button>
		</div>
	{/if}
	{#if composerOpen}
		<div class="composer">
			<blockquote>{selectedPassage}</blockquote>
			{#if composerKind === 'question'}
				<label for="reader-question-text">Question text</label>
				<input id="reader-question-text" bind:value={composerText} />
				<button
					type="button"
					onclick={submitComposer}
					disabled={composerSaving || !composerText.trim()}
					>{composerSaving ? 'Saving…' : 'Save question'}</button
				>
			{:else}
				<label for="reader-annotation-text">Annotation</label>
				<textarea id="reader-annotation-text" rows="4" bind:value={composerText}></textarea>
				<button
					type="button"
					onclick={submitComposer}
					disabled={composerSaving || !composerText.trim()}
					>{composerSaving ? 'Saving…' : 'Save annotation'}</button
				>
			{/if}
			{#if composerError}<div class="error" role="alert">{composerError}</div>{/if}
		</div>
	{/if}
	{#if openQuestion}
		<AnnotationCard
			question={openQuestion}
			passage={openPassage}
			onClose={() => (openQuestion = null)}
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
	article :global(mark.annotation-chip) {
		display: inline-block;
		width: 0.65rem;
		height: 0.65rem;
		border-radius: 999px;
		padding: 0;
		vertical-align: middle;
	}
	.toolbar,
	.composer {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-items: end;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
		padding: 0.75rem;
	}
	.composer {
		display: grid;
	}
	.composer blockquote {
		margin: 0;
		padding: 0.5rem 0.65rem;
		background: #fef3c7;
		color: #78350f;
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
