<script lang="ts">
	import type { Question } from '$lib/types/question';

	export let questions: Question[] = [];
	export let orderedIds: string[] = [];
	export let onSelect: ((id: string) => void) | undefined = undefined;
	export let onNext: (() => void) | undefined = undefined;

	$: openQuestions = (questions ?? [])
		.filter((question) => question.kind !== 'annotation' && question.status !== 'answered')
		.map((question, index) => ({ question, index }))
		.sort((a, b) => {
			const aOrder = orderedIds.indexOf(a.question.id);
			const bOrder = orderedIds.indexOf(b.question.id);

			if (aOrder === -1 && bOrder === -1) return a.index - b.index;
			if (aOrder === -1) return 1;
			if (bOrder === -1) return -1;
			return aOrder - bOrder;
		})
		.map(({ question }) => question);
</script>

<aside
	class="rail"
	aria-label="Open questions in this note"
	data-note-reader-navigation
>
	{#if openQuestions.length === 0}
		<p class="empty">No open questions in this note</p>
	{:else}
		<ul>
			{#each openQuestions as question (question.id)}
				<li>
					<button type="button" class="question" onclick={() => onSelect?.(question.id)}>
						<span class="question-text">{question.questionText}</span>
						<span class="status">{question.status}</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}

	<button
		type="button"
		class="next"
		disabled={openQuestions.length === 0 || !onNext}
		onclick={() => onNext?.()}>Next unanswered</button
	>
</aside>

<style>
	.rail {
		position: relative;
		z-index: 21;
		padding: 0.8rem;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	ul {
		list-style: none;
		padding: 0;
		margin: 0;
		display: grid;
		gap: 0.4rem;
	}
	.question,
	.next {
		width: 100%;
		padding: 0.6rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
		background: white;
		text-align: left;
	}
	.question {
		display: grid;
		gap: 0.2rem;
	}
	.question:hover,
	.question:focus-visible,
	.next:hover:not(:disabled),
	.next:focus-visible:not(:disabled) {
		border-color: #2563eb;
	}
	.question-text {
		color: #172554;
	}
	.status {
		color: #64748b;
		font-size: 0.8rem;
	}
	.next {
		margin-top: 0.7rem;
		font-weight: 650;
		text-align: center;
	}
	.next:disabled {
		cursor: not-allowed;
		opacity: 0.55;
	}
	.empty {
		margin: 0;
		color: #64748b;
	}
</style>
