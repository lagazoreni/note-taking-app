<script lang="ts">
	import { onMount } from 'svelte';
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
	type AnchorRect = {
		top: number;
		left: number;
		right: number;
		bottom: number;
		width: number;
		height: number;
	};
	let openAnchorRect: AnchorRect | null = null;
	let actionLayer: HTMLDivElement | null = null;
	let readerElement: HTMLDivElement | null = null;
	let toolbarPosition: { top: number; left: number } | null = null;
	let layerPlacement: 'above' | 'below' = 'below';
	let mouseSelectionActive = false;

	const POSITION_MARGIN = 8;
	const POSITION_GAP = 8;
	const ESTIMATED_LAYER_WIDTH = 360;
	const ESTIMATED_LAYER_HEIGHT = 96;

	$: html = renderNoteHtml(markdown, questions);

	function selectedText(): string {
		return window.getSelection()?.toString().trim() ?? '';
	}

	function clamp(value: number, min: number, max: number): number {
		return Math.min(Math.max(value, min), max);
	}

	function selectionRange(): Range | null {
		const selection = window.getSelection();
		if (!selection || selection.rangeCount === 0 || selection.isCollapsed) return null;
		return selection.getRangeAt(0);
	}

	function isRangeInsideReader(range: Range): boolean {
		if (!readerElement) return true;
		const ancestor = (range as Range & { commonAncestorContainer?: Node }).commonAncestorContainer;
		if (!ancestor) return true;
		const node = ancestor.nodeType === 3 ? ancestor.parentNode : ancestor;
		return Boolean(node && readerElement.contains(node));
	}

	function selectionRect(): DOMRect | null {
		const range = selectionRange();
		if (!range || !isRangeInsideReader(range)) return null;
		if (typeof range.getBoundingClientRect !== 'function') return null;
		const rect = range.getBoundingClientRect();
		if (!rect) return null;
		const values = [rect.top, rect.left, rect.right, rect.bottom, rect.width, rect.height];
		if (values.some((value) => !Number.isFinite(value))) return null;
		if (rect.width === 0 && rect.height === 0) return null;
		return rect;
	}

	function calculateToolbarPosition(rect: DOMRect): { top: number; left: number } {
		const viewportWidth = Math.max(window.innerWidth || document.documentElement.clientWidth || 1, 1);
		const viewportHeight = Math.max(window.innerHeight || document.documentElement.clientHeight || 1, 1);
		const layerWidth = Math.min(
			ESTIMATED_LAYER_WIDTH,
			Math.max(1, viewportWidth - POSITION_MARGIN * 2)
		);
		const layerHeight = Math.min(
			ESTIMATED_LAYER_HEIGHT,
			Math.max(1, viewportHeight - POSITION_MARGIN * 2)
		);
		const minLeft = Math.min(POSITION_MARGIN, viewportWidth);
		const maxLeft = Math.max(minLeft, viewportWidth - layerWidth - minLeft);
		const minTop = Math.min(POSITION_MARGIN, viewportHeight);
		const maxTop = Math.max(minTop, viewportHeight - layerHeight - minTop);

		let top = rect.bottom + POSITION_GAP;
		layerPlacement = 'below';
		if (top > maxTop) {
			const aboveSelection = rect.top - layerHeight - POSITION_GAP;
			if (aboveSelection >= minTop) {
				top = aboveSelection;
				layerPlacement = 'above';
			} else {
				top = maxTop;
			}
		}

		return {
			top: Math.round(clamp(top, minTop, maxTop)),
			left: Math.round(clamp(rect.left, minLeft, maxLeft))
		};
	}

	function layerStyle(): string {
		if (!toolbarPosition) return '';
		return `position: fixed; top: ${toolbarPosition.top}px; left: ${toolbarPosition.left}px;`;
	}

	function updateSelectionPosition(): boolean {
		if (composerOpen || pickerOpen || mouseSelectionActive) return false;
		const range = selectionRange();
		if (range && !isRangeInsideReader(range)) return false;
		const text = selectedText();
		if (!text) {
			// A click inside the reader can collapse the browser selection. Keep the
			// action layer anchored until the user explicitly dismisses it or makes
			// another non-empty selection.
			if (showToolbar && selectedPassage) return true;
			clearSelection();
			return false;
		}

		const rect = selectionRect();
		if (!rect) {
			if (showToolbar && selectedPassage) return true;
			clearSelection();
			return false;
		}

		selectedPassage = text;
		toolbarPosition = calculateToolbarPosition(rect);
		showToolbar = true;
		return true;
	}

	function clearSelection() {
		showToolbar = false;
		selectedPassage = '';
		toolbarPosition = null;
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

	function onReaderMouseDown(event: MouseEvent) {
		if (event.button === 0) mouseSelectionActive = true;
	}

	function onMouseUp() {
		mouseSelectionActive = false;
		updateSelectionPosition();
	}

	function onWindowMouseUp() {
		if (!mouseSelectionActive) return;
		mouseSelectionActive = false;
		updateSelectionPosition();
	}

	function onContextMenu(event: MouseEvent) {
		if (updateSelectionPosition()) event.preventDefault();
	}

	function onSelectionChange() {
		// Browsers can emit several selectionchange events while the pointer is
		// still dragging. Wait for mouseup so the toolbar cannot flash over the
		// passage or replace a stable action layer with an intermediate range.
		if (mouseSelectionActive) return;
		updateSelectionPosition();
	}

	function onTouchEnd() {
		updateSelectionPosition();
	}

	function getAnchorRect(element: Element): AnchorRect | null {
		if (typeof element.getBoundingClientRect !== 'function') return null;
		const rect = element.getBoundingClientRect();
		const values = [rect.top, rect.left, rect.right, rect.bottom];
		if (values.some((value) => !Number.isFinite(value))) return null;
		const width = Number.isFinite(rect.width) ? rect.width : Math.max(0, rect.right - rect.left);
		const height = Number.isFinite(rect.height) ? rect.height : Math.max(0, rect.bottom - rect.top);
		return {
			top: rect.top,
			left: rect.left,
			right: rect.right,
			bottom: rect.bottom,
			width,
			height
		};
	}

	onMount(() => {
		if (!readerElement) return;

		readerElement.addEventListener('selectionchange', onSelectionChange);
		readerElement.addEventListener('touchend', onTouchEnd, { passive: true });
		// Browsers dispatch selectionchange on document rather than the selected
		// element. Keep the reader listener for direct/test dispatches and use the
		// document listener as the browser fallback.
		document.addEventListener('selectionchange', onSelectionChange);

		return () => {
			readerElement?.removeEventListener('selectionchange', onSelectionChange);
			readerElement?.removeEventListener('touchend', onTouchEnd);
			document.removeEventListener('selectionchange', onSelectionChange);
		};
	});

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
		openAnchorRect = getAnchorRect(mark);
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
			openAnchorRect = null;
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
		// Do not treat pointer movement outside the reader during an active
		// selection as an outside-click dismissal. The window mouseup handler
		// finalizes the selection once the drag ends.
		if (mouseSelectionActive) return;
		const target = event.target as Node | null;
		if (actionLayer && target && actionLayer.contains(target)) return;
		if (readerElement && target && readerElement.contains(target)) return;

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
	onmouseup={onWindowMouseUp}
	onkeydown={onWindowKeydown}
/>

<div class="reader" bind:this={readerElement}>
	<!-- eslint-disable-next-line svelte/no-unused-svelte-ignore -->
	<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -- Text selection and highlight clicks are intentionally handled on the reading surface. -->
	<article
		aria-label="Reading note"
		onmousedown={onReaderMouseDown}
		onmouseup={onMouseUp}
		oncontextmenu={onContextMenu}
		onclick={onClick}
	>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -- renderNoteHtml sanitizes the generated HTML with DOMPurify. -->
		{@html html}
	</article>
	{#if showToolbar}
		<div
			class="toolbar"
			role="toolbar"
			aria-label="Selection actions"
			data-placement={layerPlacement}
			style={layerStyle()}
			bind:this={actionLayer}
		>
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
		<div
			class="picker-layer"
			data-placement={layerPlacement}
			style={layerStyle()}
			bind:this={actionLayer}
		>
			<QuestionPicker questions={linkableQuestions} onSelect={selectExisting} />
			{#if pickerError}<div class="error" role="alert">{pickerError}</div>{/if}
			<button type="button" class="secondary" onclick={cancelPicker} disabled={pickerSaving}>
				Cancel
			</button>
		</div>
	{/if}
	{#if composerOpen}
		<div
			class="composer"
			data-placement={layerPlacement}
			style={layerStyle()}
			bind:this={actionLayer}
		>
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
			anchorRect={openAnchorRect}
			onClose={() => {
				openQuestion = null;
				openAnchorRect = null;
			}}
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
		position: fixed;
		z-index: 30;
		box-sizing: border-box;
		max-width: calc(100vw - 1rem);
		background: #f8fafc;
		border: 1px solid #cbd5e1;
		border-radius: 0.7rem;
		padding: 0.75rem;
		box-shadow: 0 12px 30px rgb(15 23 42 / 18%), 0 2px 8px rgb(15 23 42 / 10%);
		isolation: isolate;
		will-change: top, left;
	}
	.toolbar {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-items: center;
		width: max-content;
	}
	.composer,
	.picker-layer {
		display: grid;
		gap: 0.65rem;
		width: min(22.5rem, calc(100vw - 1rem));
	}
	.toolbar::before,
	.composer::before,
	.picker-layer::before {
		content: '';
		position: absolute;
		z-index: 0;
		width: 0.75rem;
		height: 0.75rem;
		background: #f8fafc;
		border: 1px solid #cbd5e1;
		transform: rotate(45deg);
		pointer-events: none;
	}
	.toolbar[data-placement='below']::before,
	.composer[data-placement='below']::before,
	.picker-layer[data-placement='below']::before {
		top: -0.4rem;
		left: 1.35rem;
		border-right: 0;
		border-bottom: 0;
	}
	.toolbar[data-placement='above']::before,
	.composer[data-placement='above']::before,
	.picker-layer[data-placement='above']::before {
		bottom: -0.4rem;
		left: 1.35rem;
		border-top: 0;
		border-left: 0;
	}
	.composer blockquote {
		max-height: 7rem;
		margin: 0;
		padding: 0.5rem 0.65rem;
		background: #fef3c7;
		color: #78350f;
		overflow: auto;
	}
	.composer-actions {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	.keyboard-hint {
		color: #475569;
		font-size: 0.75rem;
		white-space: nowrap;
	}
	label {
		font-weight: 650;
	}
	input,
	textarea {
		box-sizing: border-box;
		width: 100%;
		min-height: 44px;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	button {
		width: max-content;
		min-width: 44px;
		min-height: 44px;
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
	.picker-layer :global(button) {
		min-width: 44px;
		min-height: 44px;
	}
	.picker-layer :global(input) {
		box-sizing: border-box;
		min-height: 44px;
	}
	.error {
		padding: 0.6rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.4rem;
	}
	@media (max-width: 640px) {
		.toolbar {
			left: 0.5rem !important;
			width: calc(100vw - 1rem);
			max-width: calc(100vw - 1rem);
			padding: 0.5rem;
			align-items: stretch;
		}
		.toolbar button {
			flex: 1 1 calc(50% - 0.5rem);
		}
		.keyboard-hint {
			flex-basis: 100%;
			order: -1;
			white-space: normal;
		}
		.toolbar[data-placement='above'],
		.composer[data-placement='above'],
		.picker-layer[data-placement='above'] {
			top: auto !important;
			bottom: 0.5rem;
		}
		.composer,
		.picker-layer {
			left: 0.5rem !important;
			width: calc(100vw - 1rem);
			max-width: calc(100vw - 1rem);
			max-height: calc(100vh - 1rem);
			max-height: calc(100dvh - 1rem);
			overflow-y: auto;
			overscroll-behavior: contain;
		}
		.toolbar::before,
		.composer::before,
		.picker-layer::before {
			display: none;
		}
	}
</style>
