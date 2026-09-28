<script lang="ts">
	import { onMount } from 'svelte';
	import AppNavigation from '$lib/components/AppNavigation.svelte';
	import { workspacesApi } from '$lib/api/workspaces';
	import { markServiceAvailable, markServiceUnavailable } from '$lib/stores/connectivity';
	import {
		beginWorkspaceHydration,
		failWorkspaceHydration,
		setWorkspaceList,
		workspaceLoadError,
		workspaceLoading
	} from '$lib/stores/workspace';
	import { pushToast, toasts, dismissToast } from '$lib/stores/toast';
	import { evaluateReminders } from '$lib/reminders/reminders';

	let loaded = false;
	async function loadWorkspaces(): Promise<void> {
		beginWorkspaceHydration();
		try {
			const result = await workspacesApi.list();
			setWorkspaceList(result.items ?? []);
			markServiceAvailable();
		} catch (error) {
			failWorkspaceHydration(error instanceof Error ? error.message : 'Could not load workspaces');
			markServiceUnavailable();
		}
	}

	onMount(async () => {
		await loadWorkspaces();
		await evaluateReminders();
		loaded = true;
	});
</script>

<svelte:head>
	<title>Noted · Questions in context</title>
	<meta name="description" content="A focused local-first note and question workspace" />
</svelte:head>

<AppNavigation />
{#if !loaded || $workspaceLoading}
	<div class="loading" role="status">Loading your local workspace…</div>
{/if}
{#if $workspaceLoadError}
	<div class="load-error" role="alert">
		<span>Could not load your workspaces: {$workspaceLoadError}</span>
		<button onclick={loadWorkspaces}>Retry</button>
	</div>
{/if}
{#if $toasts.length > 0}
	<div class="toasts" aria-live="polite">
		{#each $toasts as toast (toast.id)}
			<div class="toast {toast.kind}" role={toast.kind === 'error' ? 'alert' : 'status'}>
				<span>{toast.message}</span><button
					aria-label="Dismiss notification"
					onclick={() => dismissToast(toast.id)}>×</button
				>
			</div>
		{/each}
	</div>
{/if}
<main><slot /></main>

<style>
	:global(*) {
		box-sizing: border-box;
	}
	:global(html) {
		font-family: Inter, ui-sans-serif, system-ui, sans-serif;
		color: #172033;
		background: #f8fafc;
	}
	:global(body) {
		margin: 0;
		min-width: 280px;
	}
	:global(button),
	:global(input),
	:global(textarea),
	:global(select) {
		font: inherit;
	}
	:global(button),
	:global(.button) {
		cursor: pointer;
	}
	main {
		max-width: 1180px;
		margin: 0 auto;
		padding: 1.5rem;
	}
	.loading {
		padding: 0.45rem 1rem;
		color: #475569;
		background: #e0f2fe;
		font-size: 0.9rem;
	}
	.load-error {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 0.65rem 1rem;
		color: #991b1b;
		background: #fef2f2;
	}
	.load-error button {
		border: 1px solid #fca5a5;
		border-radius: 0.35rem;
		background: white;
		color: #991b1b;
		padding: 0.35rem 0.65rem;
	}
	.toasts {
		position: fixed;
		right: 1rem;
		bottom: 1rem;
		z-index: 20;
		display: grid;
		gap: 0.5rem;
		max-width: min(28rem, calc(100vw - 2rem));
	}
	.toast {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		background: white;
		box-shadow: 0 8px 30px #0f172a2b;
		border-left: 4px solid #2563eb;
	}
	.toast.success {
		border-color: #16a34a;
	}
	.toast.error {
		border-color: #dc2626;
	}
	.toast button {
		border: 0;
		background: transparent;
		font-size: 1.2rem;
	}
	@media (max-width: 620px) {
		main {
			padding: 1rem;
		}
	}
</style>
