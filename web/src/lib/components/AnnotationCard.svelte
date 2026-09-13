<script lang="ts">
	import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
	import type { Question } from '$lib/types/question';
	export let question: Question;
	export let passage: string;
	export let onClose: (() => void) | undefined = undefined;
	export let onInsertAnswer: (() => Promise<void> | void) | undefined = undefined;
	type AnchorRect = {
		top: number;
		left: number;
		right: number;
		bottom: number;
		width?: number;
		height?: number;
	};
	type CardPosition = { top: number; left: number };
	export let anchorRect: AnchorRect | null = null;
	export let cardPosition: CardPosition | null = null;
	$: isAnnotation = question.kind === 'annotation';
	$: canInsertAnswer = !isAnnotation && question.status === 'answered' && Boolean(question.answerMarkdown);
	let draft = question.questionText;
	$: draft = question.questionText;

	function finiteNumber(value: unknown): value is number {
		return typeof value === 'number' && Number.isFinite(value);
	}

	function cardStyle(): string {
		const viewportWidth = Math.max(
			1,
			typeof window === 'undefined' ? 1 : window.innerWidth || document.documentElement.clientWidth || 1
		);
		const viewportHeight = Math.max(
			1,
			typeof window === 'undefined' ? 1 : window.innerHeight || document.documentElement.clientHeight || 1
		);
		const margin = 16;
		const gap = 12;
		const cardWidth = Math.min(512, Math.max(1, viewportWidth - margin * 2));
		const cardHeight = Math.min(640, Math.max(1, viewportHeight * 0.8));
		const minLeft = Math.min(margin, Math.max(0, viewportWidth - cardWidth));
		const maxLeft = Math.max(minLeft, viewportWidth - cardWidth - minLeft);
		const minTop = Math.min(margin, Math.max(0, viewportHeight - cardHeight));
		const maxTop = Math.max(minTop, viewportHeight - cardHeight - minTop);

		let top: number;
		let left: number;
		if (cardPosition && finiteNumber(cardPosition.top) && finiteNumber(cardPosition.left)) {
			// A caller-provided cardPosition is already an intended placement. Keep
			// its coordinates stable while preventing a malformed value from putting
			// the overlay outside the viewport entirely.
			return [
				'position: fixed',
				`top: ${Math.round(Math.min(Math.max(cardPosition.top, 0), viewportHeight))}px`,
				`left: ${Math.round(Math.min(Math.max(cardPosition.left, 0), viewportWidth))}px`,
				'right: auto',
				'bottom: auto'
			].join('; ');
		} else if (
			anchorRect &&
			finiteNumber(anchorRect.top) &&
			finiteNumber(anchorRect.left) &&
			finiteNumber(anchorRect.right)
		) {
			const anchorRight = anchorRect.right;
			const anchorLeft = anchorRect.left;
			top = anchorRect.top;
			left = anchorRight + gap;
			if (left + cardWidth > viewportWidth - margin) {
				left = anchorLeft - cardWidth - gap;
			}
			if (top + cardHeight > viewportHeight - margin && top - cardHeight - gap >= minTop) {
				top -= cardHeight + gap;
			}
		} else {
			return '';
		}

		return [
			'position: fixed',
			`top: ${Math.round(Math.min(Math.max(top, minTop), maxTop))}px`,
			`left: ${Math.round(Math.min(Math.max(left, minLeft), maxLeft))}px`,
			'right: auto',
			'bottom: auto'
		].join('; ');
	}
</script>

<div
	class="card"
	role="dialog"
	aria-modal="true"
	aria-labelledby="annotation-card-title"
	style={cardStyle()}
>
	<div class="header">
		<h2 id="annotation-card-title">{isAnnotation ? 'Annotation details' : 'Question'}</h2>
		<button type="button" class="close" onclick={() => onClose?.()}>Close</button>
	</div>
	<blockquote>{passage}</blockquote>
	<p class="text">{question.questionText}</p>
	{#if canInsertAnswer}
		<button type="button" class="insert-answer" onclick={() => onInsertAnswer?.()}>Insert answer into note</button>
	{/if}
	{#if isAnnotation}
		<label for="annotation-body">Annotation</label>
		<textarea id="annotation-body" rows="4" bind:value={draft}></textarea>
	{:else}
		<QuestionLifecycle {question} />
	{/if}
</div>

<style>
	.card {
		position: fixed;
		top: auto;
		right: 1rem;
		bottom: 1rem;
		left: auto;
		box-sizing: border-box;
		z-index: 20;
		width: min(32rem, calc(100vw - 2rem));
		max-height: min(80vh, 40rem);
		overflow: auto;
		background: #fff;
		border: 1px solid #cbd5e1;
		border-radius: 0.7rem;
		padding: 1rem;
		box-shadow: 0 12px 40px rgb(15 23 42 / 0.18);
	}
	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.75rem;
	}
	h2 {
		margin: 0;
		color: #172554;
		font-size: 1.1rem;
	}
	.close {
		border: 1px solid #cbd5e1;
		background: #f8fafc;
		border-radius: 0.35rem;
		padding: 0.35rem 0.65rem;
	}
	blockquote {
		margin: 0.75rem 0;
		padding: 0.65rem 0.75rem;
		background: #fef3c7;
		border-left: 3px solid #d97706;
		color: #78350f;
	}
	.text {
		margin: 0 0 0.75rem;
		color: #0f172a;
	}
	.insert-answer {
		margin: 0 0 0.75rem;
		border: 1px solid #2563eb;
		background: #eff6ff;
		color: #1d4ed8;
		border-radius: 0.35rem;
		padding: 0.45rem 0.7rem;
	}
	label {
		font-weight: 650;
	}
	textarea {
		width: 100%;
		margin-top: 0.35rem;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	@media (max-width: 640px) {
		.card {
			top: auto !important;
			right: 0 !important;
			bottom: 0 !important;
			left: 0 !important;
			width: 100vw;
			max-width: 100vw;
			max-height: min(85vh, 40rem);
			border-radius: 0.9rem 0.9rem 0 0;
		}
	}
</style>
