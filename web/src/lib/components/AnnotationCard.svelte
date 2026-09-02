<script lang="ts">
	import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
	import type { Question } from '$lib/types/question';
	export let question: Question;
	export let passage: string;
	export let onClose: (() => void) | undefined = undefined;
	export let onInsertAnswer: (() => Promise<void> | void) | undefined = undefined;
	$: isAnnotation = question.kind === 'annotation';
	$: canInsertAnswer = !isAnnotation && question.status === 'answered' && Boolean(question.answerMarkdown);
	let draft = question.questionText;
	$: draft = question.questionText;
</script>

<div class="card" role="dialog" aria-modal="true" aria-labelledby="annotation-card-title">
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
		inset: auto 1rem 1rem auto;
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
</style>
