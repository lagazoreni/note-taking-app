<script lang="ts">
	import type { Topic } from '$lib/types/workspace';
	export let topics: Topic[] = [];
	export let value: string | null = null;
</script>

<label for="topic">Topic</label><select id="topic" bind:value
	><option value="">No topic</option>{#each topics as topic}<option value={topic.id}
			>{topic.name}</option
		>{/each}</select
>

<style>
	label {
		font-weight: 650;
		display: block;
		margin-bottom: 0.25rem;
	}
	select {
		padding: 0.55rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
	}
</style>
