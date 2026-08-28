<script lang="ts">
	import { onMount } from 'svelte';
	import { currentWorkspaceId } from '$lib/stores/workspace';
	import { questionsApi } from '$lib/api/questions';
	import QuestionList from '$lib/components/QuestionList.svelte';
	import type { Question } from '$lib/types/question';
	let questions: Question[] = [];
	let loading = false;
	let error = '';
	$: workspaceId = $currentWorkspaceId;
	$: if (workspaceId) load(workspaceId);
	onMount(() => {
		if (workspaceId) load(workspaceId);
	});
	async function load(id: string) {
		loading = true;
		try {
			questions = (
				await questionsApi.list({ workspaceId: id, status: ['answered'], kind: 'question' })
			).items;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load answered questions';
		} finally {
			loading = false;
		}
	}
</script>

<section>
	<p class="eyebrow">History</p>
	<h1>Answered questions</h1>
	{#if loading}<div role="status">Loading answered questions…</div>{:else if error}<div
			class="error"
			role="alert"
		>
			{error}
		</div>{:else}<QuestionList
			{questions}
			answered
			emptyMessage="No answered questions yet."
		/>{/if}
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
		padding: 1rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.5rem;
	}
</style>
