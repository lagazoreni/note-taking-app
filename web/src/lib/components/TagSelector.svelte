<script lang="ts">
	import type { Tag } from '$lib/types/workspace';
	export let tags: Tag[] = [];
	export let selected: string[] = [];
</script>

<fieldset>
	<legend>Tags</legend>{#each tags as tag}<label
			><input
				type="checkbox"
				value={tag.id}
				bind:group={selected}
			/>{tag.name}{#if tag.availableWorkspaceIds.length > 1}<small> shared</small>{/if}</label
		>{/each}
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
	label {
		display: block;
		margin: 0.3rem 0;
	}
	small {
		color: #64748b;
	}
</style>
