<script lang="ts">
	import { onDestroy } from 'svelte';
	import { currentWorkspaceId, workspaceHydrated } from '$lib/stores/workspace';
	import { notesApi } from '$lib/api/notes';
	import NoteTree from '$lib/components/NoteTree.svelte';
	import type { Note } from '$lib/types/note';

	let notes: Note[] = [];
	let loading = false;
	let error = '';
	let controller: AbortController | null = null;

	$: workspaceId = $currentWorkspaceId;
	$: hydrated = $workspaceHydrated;
	$: if (hydrated && workspaceId) loadNotes(workspaceId);

	async function loadNotes(id: string) {
		controller?.abort();
		controller = new AbortController();
		const signal = controller.signal;
		loading = true;
		error = '';
		try {
			const result = await notesApi.list(id, signal);
			if (signal.aborted) return;
			notes = result.items ?? [];
		} catch (cause) {
			if (signal.aborted) return;
			error = cause instanceof Error ? cause.message : 'Could not load notes';
		} finally {
			if (!signal.aborted) loading = false;
		}
	}

	onDestroy(() => controller?.abort());
</script>

<section class="page">
	<div class="heading">
		<div>
			<p class="eyebrow">{workspaceId ? 'Current workspace' : 'Choose a workspace'}</p>
			<h1>Notes</h1>
		</div>
		<a class="button" href="/notes/new">New note</a>
	</div>
	{#if !hydrated}
		<div role="status">Loading notes…</div>
	{:else if !workspaceId}
		<div class="empty">
			<h2>Create a workspace</h2>
			<p>Choose a workspace to start taking notes.</p>
			<a href="/workspaces">Go to workspaces</a>
		</div>
	{:else if loading}
		<div role="status">Loading notes…</div>
	{:else if error}
		<div class="error" role="alert">
			{error}
			<button type="button" onclick={() => workspaceId && loadNotes(workspaceId)}>Retry</button>
		</div>
	{:else if notes.length === 0}
		<div class="empty">
			<h2>Your notes will appear here</h2>
			<p>Create a note and add questions inline to start tracking follow-up work.</p>
		</div>
	{:else}
		<NoteTree {notes} />
	{/if}
</section>

<style>
	.page {
		max-width: 60rem;
		margin: auto;
	}
	.heading {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
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
	.button {
		padding: 0.6rem 1rem;
		background: #1d4ed8;
		color: #fff;
		border-radius: 0.4rem;
		text-decoration: none;
	}
	.empty {
		margin-top: 1.5rem;
		background: #fff;
		border: 1px dashed #94a3b8;
		border-radius: 0.6rem;
		padding: 2rem;
		text-align: center;
	}
	.empty p {
		color: #64748b;
	}
	.empty a {
		color: #1d4ed8;
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
