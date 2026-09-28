<script lang="ts">
	import { workspacesApi } from '$lib/api/workspaces';
	import {
		addWorkspace,
		currentWorkspaceId,
		selectWorkspace,
		workspaces
	} from '$lib/stores/workspace';
	import { pushToast } from '$lib/stores/toast';
	import type { Workspace } from '$lib/types/workspace';
	let name = '';
	let saving = false;
	async function createWorkspace() {
		if (!name.trim()) return;
		saving = true;
		try {
			const workspace = await workspacesApi.create(name.trim());
			addWorkspace(workspace);
			name = '';
			pushToast('Workspace created', 'success');
		} catch (error) {
			pushToast(error instanceof Error ? error.message : 'Could not create workspace', 'error');
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head><title>Workspaces · Noted</title></svelte:head>
<section class="page">
	<div class="heading">
		<div>
			<p class="eyebrow">Contexts</p>
			<h1>Workspaces</h1>
		</div>
	</div>
	<form
		onsubmit={(event) => {
			event.preventDefault();
			createWorkspace();
		}}
		class="create"
	>
		<label for="workspace-name">New workspace</label>
		<div>
			<input
				id="workspace-name"
				bind:value={name}
				maxlength="100"
				placeholder="Work, Personal, Studies"
			/><button class="button" disabled={saving || !name.trim()}>Create</button>
		</div>
	</form>
	{#if $workspaces.length === 0}<div class="empty">
			<h2>Start with a workspace</h2>
			<p>Separate areas keep notes and questions focused.</p>
		</div>{:else}<div class="list">
			{#each $workspaces as workspace}<button
					class:selected={$currentWorkspaceId === workspace.id}
					onclick={() => selectWorkspace(workspace.id)}
					><strong>{workspace.name}</strong><span>Version {workspace.version}</span></button
				>{/each}
		</div>{/if}
</section>

<style>
	.page {
		max-width: 52rem;
		margin: 0 auto;
	}
	.eyebrow {
		color: #2563eb;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		font-size: 0.8rem;
		font-weight: 700;
	}
	h1 {
		color: #172554;
	}
	.create {
		margin: 1.5rem 0;
		padding: 1rem;
		background: white;
		border: 1px solid #e2e8f0;
		border-radius: 0.6rem;
	}
	.create label {
		display: block;
		font-weight: 650;
		margin-bottom: 0.4rem;
	}
	.create div {
		display: flex;
		gap: 0.5rem;
	}
	input {
		width: 100%;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	.button {
		border: 0;
		background: #1d4ed8;
		color: white;
		padding: 0.6rem 1rem;
		border-radius: 0.4rem;
	}
	.button:disabled {
		opacity: 0.5;
	}
	.list {
		display: grid;
		gap: 0.6rem;
	}
	.list button {
		display: flex;
		justify-content: space-between;
		text-align: left;
		border: 1px solid #cbd5e1;
		border-radius: 0.55rem;
		background: white;
		padding: 1rem;
	}
	.list button.selected {
		border-color: #2563eb;
		box-shadow: 0 0 0 2px #bfdbfe;
	}
	.list span {
		color: #64748b;
		font-size: 0.85rem;
	}
	.empty {
		padding: 2rem;
		background: #eff6ff;
		border-radius: 0.6rem;
	}
</style>
