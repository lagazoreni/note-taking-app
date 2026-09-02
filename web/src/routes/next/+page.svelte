<script lang="ts">
	import { onMount } from 'svelte';
	import { questionsApi, type QuestionQuery } from '$lib/api/questions';
	import NextQueue from '$lib/components/NextQueue.svelte';
	import { buildNextQueue } from '$lib/questions/nextQueue';
	import { currentWorkspaceId } from '$lib/stores/workspace';
	import type { Question } from '$lib/types/question';

	let questions: Question[] = [];
	let loading = false;
	let error = '';
	let hasMoreQuestions = false;
	$: workspaceId = $currentWorkspaceId;
	$: if (workspaceId) loadQuestions(workspaceId);

	const today = new Date().toISOString().slice(0, 10);

	async function loadQuestions(id: string) {
		loading = true;
		error = '';
		hasMoreQuestions = false;
		try {
			const query: QuestionQuery = {
				workspaceId: id,
				status: ['unanswered', 'in_progress', 'deferred'],
				kind: 'question',
				pageSize: 200
			};
			const items: Question[] = [];
			let nextCursor: string | null = null;

			for (let page = 0; page < 5; page += 1) {
				const result = await questionsApi.list(
					nextCursor ? { ...query, cursor: nextCursor } : query
				);
				items.push(...result.items);
				nextCursor = result.nextCursor ?? null;
				if (!nextCursor) break;
			}

			questions = items;
			hasMoreQuestions = nextCursor !== null;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load questions';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		if (workspaceId) loadQuestions(workspaceId);
	});
</script>

<section>
	<p class="eyebrow">Follow-up</p>
	<h1>Next</h1>
	{#if loading}
		<div role="status">Loading questions…</div>
	{:else if error}
		<div class="error" role="alert">
			{error}<button onclick={() => workspaceId && loadQuestions(workspaceId)}>Retry</button>
		</div>
	{:else}
		{#if hasMoreQuestions}
			<p role="status">Showing the first loaded questions.</p>
		{/if}
		<NextQueue sections={buildNextQueue(questions, today)} />
	{/if}
</section>

<style>
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

	.error {
		margin: 1rem 0;
		padding: 1rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.5rem;
	}

	.error button {
		margin-left: 1rem;
	}
</style>
