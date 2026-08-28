<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { questionsApi } from '$lib/api/questions';
	import type { Question } from '$lib/types/question';
	import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
	import QuestionSchedule from '$lib/components/QuestionSchedule.svelte';
	let question: Question | null = null;
	let questionId = '';
	let loading = true;
	let error = '';
	$: questionId = page.params.questionId ?? '';
	onMount(async () => {
		try {
			question = await questionsApi.get(questionId);
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load question';
		} finally {
			loading = false;
		}
	});
	function saved(value: Question) {
		question = value;
	}
</script>

{#if loading}<div role="status">Loading question…</div>{:else if error}<div
		role="alert"
		class="error"
	>
		{error}
	</div>{:else if question}<section class="page">
		<p class="eyebrow">Question</p>
		<h1>{question.questionText}</h1>
		<QuestionLifecycle {question} onSave={saved} /><QuestionSchedule {question} onSave={saved} />
		<div class="linked">
			<h2>Source notes</h2>
			{#each question.linkedNotes as note}<a href={`/notes/${note.id}`}>{note.title}</a>{:else}<p>
					This question is currently unlinked.
				</p>{/each}
		</div>
	</section>{/if}

<style>
	.page {
		max-width: 58rem;
		margin: auto;
	}
	.eyebrow {
		color: #2563eb;
		text-transform: uppercase;
		font-size: 0.8rem;
		font-weight: 700;
		letter-spacing: 0.08em;
	}
	h1 {
		color: #172554;
	}
	.linked {
		margin-top: 1.5rem;
		padding: 1rem;
		background: #fff;
		border-radius: 0.5rem;
		border: 1px solid #e2e8f0;
	}
	.linked a {
		display: block;
		color: #1d4ed8;
		margin: 0.3rem 0;
	}
	.error {
		padding: 1rem;
		background: #fef2f2;
		color: #991b1b;
	}
</style>
