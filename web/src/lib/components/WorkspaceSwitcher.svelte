<script lang="ts">
	import {
		currentWorkspaceId,
		selectWorkspace,
		workspaceHydrated,
		workspaces
	} from '$lib/stores/workspace';
</script>

<label class="switcher">
	<span class="sr-only">Current workspace</span>
	<select
		aria-label="Current workspace"
		value={$currentWorkspaceId ?? ''}
		disabled={!$workspaceHydrated || $workspaces.length === 0}
		onchange={(event) => selectWorkspace((event.currentTarget as HTMLSelectElement).value)}
	>
		<option value="" disabled>Select workspace</option>
		{#each $workspaces as workspace}
			<option value={workspace.id}>{workspace.name}</option>
		{/each}
	</select>
</label>

<style>
	.switcher select {
		min-width: 9rem;
		max-width: 14rem;
		padding: 0.38rem 0.5rem;
		border: 1px solid #93c5fd;
		border-radius: 0.4rem;
		background: #eff6ff;
		color: #172554;
	}
	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
</style>
