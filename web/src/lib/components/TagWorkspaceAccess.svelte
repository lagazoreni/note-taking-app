<script lang="ts">
	import type { Tag, Workspace } from '$lib/types/workspace';
	export let tag: Tag;
	export let workspaces: Workspace[] = [];
	export let onSave: ((workspaceIds: string[]) => void) | undefined = undefined;
	let selected = tag.availableWorkspaceIds.filter((id) => id !== tag.ownerWorkspaceId);
</script>

<section class="access">
	<h2>Share “{tag.name}”</h2>
	<p>
		Sharing makes this tag selectable in another workspace. Existing assignments require
		confirmation before unsharing.
	</p>
	{#each workspaces.filter((workspace) => workspace.id !== tag.ownerWorkspaceId) as workspace}<label
			><input type="checkbox" value={workspace.id} bind:group={selected} />{workspace.name}</label
		>{/each}<button onclick={() => onSave?.(selected)}>Save access</button>
</section>

<style>
	.access {
		padding: 1rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.access h2 {
		font-size: 1rem;
	}
	.access p {
		color: #64748b;
	}
	.access label {
		display: block;
		margin: 0.4rem 0;
	}
	.access button {
		margin-top: 0.7rem;
		border: 0;
		background: #1d4ed8;
		color: #fff;
		padding: 0.55rem 0.8rem;
		border-radius: 0.35rem;
	}
</style>
