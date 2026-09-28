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
	let newQuestionText = '';
	let newQuestionOpen = false;
	let creatingQuestion = false;
	let createError = '';
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

	async function createQuestion() {
		const questionText = newQuestionText.trim();
		if (!workspaceId || !questionText || creatingQuestion) return;
		creatingQuestion = true;
		createError = '';
		try {
			await questionsApi.create({
				workspaceId,
				questionText,
				kind: 'question',
				status: 'unanswered',
				priority: 'none',
				tagIds: []
			});
			newQuestionText = '';
			newQuestionOpen = false;
			await loadQuestions(workspaceId);
		} catch (cause) {
			createError = cause instanceof Error ? cause.message : 'Could not create question';
		} finally {
			creatingQuestion = false;
		}
	}

	onMount(() => {
		if (workspaceId) loadQuestions(workspaceId);
	});
</script>

<section>
	<p class="eyebrow">Follow-up</p>
	<h1>Active questions</h1>
	<button
		type="button"
		class="new-question-trigger"
		aria-expanded={newQuestionOpen}
		onclick={() => {
			newQuestionOpen = true;
			createError = '';
		}}
	>
		New question
	</button>
	{#if newQuestionOpen}
		<form
			class="new-question"
			onsubmit={(event) => {
				event.preventDefault();
				void createQuestion();
			}}
		>
			<label for="new-question-text">Question text</label>
			<div class="new-question-controls">
				<input
					id="new-question-text"
					bind:value={newQuestionText}
					required
					placeholder="What do I need to find out?"
				/>
				<button type="submit" disabled={creatingQuestion || !newQuestionText.trim()}>
					{creatingQuestion ? 'Creating…' : 'Create question'}
				</button>
			</div>
		</form>
	{/if}
	{#if createError}<div class="create-error" role="alert">{createError}</div>{/if}
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
	.new-question-trigger {
		border: 0;
		background: #2563eb;
		color: white;
		border-radius: 0.4rem;
		padding: 0.65rem 0.85rem;
	}
	.new-question {
		display: grid;
		gap: 0.4rem;
		max-width: 42rem;
		margin: 1rem 0;
		padding: 0.85rem;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.new-question label {
		font-weight: 650;
	}
	.new-question-controls {
		display: flex;
		gap: 0.5rem;
	}
	.new-question input {
		flex: 1;
		min-width: 0;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	.new-question button {
		border: 0;
		background: #2563eb;
		color: white;
		border-radius: 0.4rem;
		padding: 0.65rem 0.85rem;
		white-space: nowrap;
	}
	.new-question button:disabled {
		opacity: 0.55;
	}
	.create-error {
		margin: 1rem 0;
		padding: 0.7rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.5rem;
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
