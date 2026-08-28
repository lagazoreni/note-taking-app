<script lang="ts">
	import { onMount } from 'svelte';
	import { currentWorkspaceId } from '$lib/stores/workspace';
	import { questionsApi } from '$lib/api/questions';
	import QuestionList from '$lib/components/QuestionList.svelte';
	import QuestionFilters from '$lib/components/QuestionFilters.svelte';
	import type { Question, QuestionStatus } from '$lib/types/question';
	let questions: Question[] = [];
	let selectedStatuses: QuestionStatus[] = ['unanswered', 'in_progress', 'deferred'];
	let sort = 'updatedAt';
	let direction: 'asc' | 'desc' = 'desc';
	let loading = false;
	let error = '';
	$: workspaceId = $currentWorkspaceId;
	$: if (workspaceId) loadQuestions(workspaceId);
	async function loadQuestions(id: string) {
		loading = true;
		error = '';
		try {
			const result = await questionsApi.list({
				workspaceId: id,
				status: selectedStatuses,
				kind: 'question',
				sort,
				direction
			});
			questions = result.items;
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
	<h1>Active questions</h1>
	<QuestionFilters
		status={selectedStatuses}
		{sort}
		{direction}
		onChange={(value) => {
			selectedStatuses = value.status;
			sort = value.sort;
			direction = value.direction;
			if (workspaceId) loadQuestions(workspaceId);
		}}
	/>{#if loading}<div role="status">Loading questions…</div>{:else if error}<div
			class="error"
			role="alert"
		>
			{error}<button onclick={() => workspaceId && loadQuestions(workspaceId)}>Retry</button>
		</div>{:else}<QuestionList {questions} emptyMessage="No active questions yet." />{/if}
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
