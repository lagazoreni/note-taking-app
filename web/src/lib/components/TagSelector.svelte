<script lang="ts">
	import type { Tag } from '$lib/types/workspace';

	export let tags: Tag[] = [];
	export let selected: string[] = [];
	export let onChange: ((value: string[]) => void) | undefined = undefined;

	$: available = tags ?? [];
	$: selectedIDs = selected ?? [];

	function toggle(tagID: string) {
		selected = selectedIDs.includes(tagID)
			? selectedIDs.filter((id) => id !== tagID)
			: [...selectedIDs, tagID];
		onChange?.(selected);
	}
</script>

<fieldset>
	<legend>Tags</legend>
	{#if available.length === 0}
		<p class="empty">No workspace tags available yet.</p>
	{:else}
		<div class="tag-list">
			{#each available as tag}
				<label for={`note-tag-${tag.id}`}>
					<input
						id={`note-tag-${tag.id}`}
						type="checkbox"
						value={tag.id}
						aria-label={tag.name}
						checked={selectedIDs.includes(tag.id)}
						onchange={() => toggle(tag.id)}
					/>
					<span>{tag.name}</span>
					{#if tag.availableWorkspaceIds.length > 1}<small>shared</small>{/if}
				</label>
			{/each}
		</div>
	{/if}
</fieldset>

<style>
	fieldset {
		border: 1px solid #e2e8f0;
		border-radius: 0.4rem;
		padding: 0.6rem;
	}
	legend {
		font-weight: 650;
	}
	.tag-list {
		display: flex;
		flex-wrap: wrap;
		gap: 0.45rem;
	}
	label {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		margin: 0.1rem 0;
		padding: 0.35rem 0.5rem;
		border: 1px solid #cbd5e1;
		border-radius: 999px;
		background: #fff;
		cursor: pointer;
	}
	input {
		margin: 0;
	}
	small {
		color: #64748b;
		font-size: 0.75rem;
	}
	.empty {
		margin: 0;
		color: #64748b;
		font-size: 0.9rem;
	}
</style>
