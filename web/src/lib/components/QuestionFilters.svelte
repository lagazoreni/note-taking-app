<script lang="ts">
	import type { QuestionStatus } from '$lib/types/question';

	export let status: QuestionStatus[] = [];
	export let sort = 'updatedAt';
	export let direction: 'asc' | 'desc' = 'desc';
	export let onChange:
		| ((value: { status: QuestionStatus[]; sort: string; direction: 'asc' | 'desc' }) => void)
		| undefined = undefined;

	let statusOpen = false;
	let statusControl: HTMLDivElement | null = null;
	const statusOptions: { value: QuestionStatus; label: string }[] = [
		{ value: 'unanswered', label: 'Unanswered' },
		{ value: 'in_progress', label: 'In progress' },
		{ value: 'deferred', label: 'Deferred' },
		{ value: 'answered', label: 'Answered' }
	];

	$: statusSummary = status.length
		? `${status.length} status${status.length === 1 ? '' : 'es'} selected`
		: 'All statuses';

	function update() {
		onChange?.({ status: [...status], sort, direction });
	}

	function toggleStatus(value: QuestionStatus) {
		status = status.includes(value) ? status.filter((item) => item !== value) : [...status, value];
		update();
	}

	function onWindowClick(event: MouseEvent) {
		if (!statusOpen) return;
		const target = event.target as Node | null;
		if (!statusControl || !target || !statusControl.contains(target)) statusOpen = false;
	}

	function onWindowKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && statusOpen) {
			statusOpen = false;
			event.preventDefault();
		}
	}
</script>

<svelte:window onclick={onWindowClick} onkeydown={onWindowKeydown} />

<div class="filters" aria-label="Question filters">
	<div class="status-control" bind:this={statusControl}>
		<label for="status-filter-button">Status</label>
		<button
			id="status-filter-button"
			type="button"
			class="status-button"
			aria-label="Status"
			aria-haspopup="true"
			aria-expanded={statusOpen}
			aria-controls="status-filter-options"
			onclick={() => (statusOpen = !statusOpen)}
		>
			{statusSummary}
		</button>
		{#if statusOpen}
			<fieldset id="status-filter-options" class="status-options" aria-label="Status options">
				<legend>Filter by status</legend>
				{#each statusOptions as option}
					<label for={`status-${option.value}`}>
						<input
							id={`status-${option.value}`}
							type="checkbox"
							checked={status.includes(option.value)}
							onchange={() => toggleStatus(option.value)}
						/>
						{option.label}
					</label>
				{/each}
			</fieldset>
		{/if}
	</div>

	<div class="select-control">
		<label for="question-sort">Sort</label>
		<select id="question-sort" bind:value={sort} onchange={update}>
			<option value="updatedAt">Updated</option>
			<option value="createdAt">Created</option>
			<option value="dueDate">Due date</option>
			<option value="priority">Priority</option>
			<option value="questionText">Question text</option>
			<option value="status">Status</option>
		</select>
	</div>

	<div class="select-control">
		<label for="question-direction">Direction</label>
		<select id="question-direction" bind:value={direction} onchange={update}>
			<option value="desc">Newest / descending</option>
			<option value="asc">Oldest / ascending</option>
		</select>
	</div>
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
	.status-control,
	.select-control {
		display: grid;
		gap: 0.25rem;
		position: relative;
	}
	.filters label {
		font-size: 0.8rem;
		font-weight: 650;
		color: #475569;
	}
	.status-button,
	.filters select {
		min-height: 2.35rem;
		padding: 0.45rem 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
		background: #fff;
		color: #0f172a;
		font: inherit;
		text-align: left;
	}
	.status-button {
		min-width: 8rem;
		cursor: pointer;
	}
	.status-options {
		position: absolute;
		z-index: 10;
		top: calc(100% + 0.35rem);
		left: 0;
		min-width: 13rem;
		margin: 0;
		padding: 0.65rem;
		background: #fff;
		border: 1px solid #cbd5e1;
		border-radius: 0.45rem;
		box-shadow: 0 8px 24px rgb(15 23 42 / 0.14);
	}
	.status-options legend {
		padding: 0 0.2rem;
		font-size: 0.8rem;
		font-weight: 700;
		color: #334155;
	}
	.status-options label {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.35rem 0.2rem;
		font-size: 0.9rem;
		color: #0f172a;
		cursor: pointer;
	}
	.status-options input {
		margin: 0;
	}
	@media (max-width: 640px) {
		.status-control,
		.select-control,
		.status-button,
		.filters select {
			width: 100%;
		}
		.status-options {
			position: static;
			width: 100%;
			box-shadow: none;
		}
	}
</style>
