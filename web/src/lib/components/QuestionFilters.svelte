<script lang="ts">
	import type { QuestionStatus } from '$lib/types/question';
	export let status: QuestionStatus[] = [];
	export let sort = 'updatedAt';
	export let direction: 'asc' | 'desc' = 'desc';
	export let onChange:
		| ((value: { status: QuestionStatus[]; sort: string; direction: 'asc' | 'desc' }) => void)
		| undefined = undefined;
	function update() {
		onChange?.({ status, sort, direction });
	}
</script>

<div class="filters" aria-label="Question filters">
	<label for="status">Status</label><select
		id="status"
		multiple
		size="4"
		bind:value={status}
		onchange={update}
		><option value="unanswered">Unanswered</option><option value="in_progress">In progress</option
		><option value="deferred">Deferred</option><option value="answered">Answered</option></select
	><label for="sort">Sort</label><select id="sort" bind:value={sort} onchange={update}
		><option value="updatedAt">Updated</option><option value="createdAt">Created</option><option
			value="dueDate">Due date</option
		><option value="priority">Priority</option><option value="questionText">Question text</option
		><option value="status">Status</option></select
	><label for="direction">Direction</label><select
		id="direction"
		bind:value={direction}
		onchange={update}
		><option value="desc">Newest / descending</option><option value="asc">Oldest / ascending</option
		></select
	>
</div>

<style>
	.filters {
		display: flex;
		align-items: end;
		gap: 0.6rem;
		flex-wrap: wrap;
		margin: 1rem 0;
		padding: 0.75rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.filters label {
		font-size: 0.8rem;
		font-weight: 650;
		color: #475569;
	}
	.filters select {
		padding: 0.45rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
		min-width: 8rem;
	}
	.filters select[multiple] {
		min-width: 10rem;
	}
</style>
