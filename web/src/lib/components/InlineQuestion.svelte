<script lang="ts">
	import type { Question, DisplayMode } from '$lib/types/question';
	export let question: Question;
	export let displayMode: DisplayMode = 'collapsed';
	export let onChange: ((question: Question) => void) | undefined = undefined;
	let expanded = displayMode === 'expanded';
</script>

<article
	class="question"
	class:expanded
	class:link={displayMode === 'link'}
	aria-label="Inline question"
>
	<div class="head">
		<button
			class="toggle"
			type="button"
			aria-expanded={expanded}
			onclick={() => {
				expanded = !expanded;
				onChange?.(question);
			}}>{expanded ? '▾' : '▸'}</button
		><strong>{question.questionText}</strong><span class="status">{question.status}</span>
	</div>
	{#if displayMode === 'link'}<a href={`/questions/${question.id}`}>Open question</a
		>{:else if expanded}<div class="details">
			{#if question.answerMarkdown}<p>{question.answerMarkdown}</p>{:else}<p class="muted">
					No answer yet.
				</p>{/if}<a href={`/questions/${question.id}`}>Open and edit</a>
		</div>{/if}
</article>

<style>
	.question {
		border: 1px solid #bfdbfe;
		background: #eff6ff;
		border-radius: 0.5rem;
		padding: 0.65rem;
	}
	.head {
		display: flex;
		align-items: center;
		gap: 0.45rem;
	}
	.toggle {
		border: 0;
		background: transparent;
		padding: 0.1rem;
	}
	.status {
		margin-left: auto;
		color: #1d4ed8;
		font-size: 0.8rem;
	}
	.details {
		padding: 0.5rem 1.4rem 0;
		color: #334155;
	}
	.question a {
		color: #1d4ed8;
	}
	.muted {
		color: #64748b;
	}
</style>
