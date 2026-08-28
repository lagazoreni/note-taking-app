<script lang="ts">
	import { goto } from '$app/navigation';
	import { currentWorkspaceId, workspaceHydrated, workspaceReady } from '$lib/stores/workspace';
	import { pushToast } from '$lib/stores/toast';
	import NoteEditor from '$lib/editor/NoteEditor.svelte';
	import type { Note } from '$lib/types/note';
	function saved(note: Note) {
		pushToast('Note saved', 'success');
		goto(`/notes/${note.id}`);
	}
</script>

<section class="page">
	<p class="eyebrow">New note</p>
	<h1>Create a note</h1>
	{#if !$workspaceHydrated}
		<div class="empty" role="status">Loading workspace setup…</div>
	{:else if $workspaceReady && $currentWorkspaceId}<NoteEditor
			workspaceId={$currentWorkspaceId}
			onSave={saved}
		/>{:else}<div class="empty" role="status">
			<p>Select or create a workspace first.</p>
			<a href="/workspaces">Go to workspaces</a>
		</div>{/if}
</section>

<style>
	.page {
		max-width: 60rem;
		margin: auto;
	}
	.eyebrow {
		color: #2563eb;
		text-transform: uppercase;
		font-size: 0.8rem;
		font-weight: 700;
		letter-spacing: 0.08em;
	}
	h1 {
		color: #172554;
	}
	.empty {
		padding: 1.5rem;
		background: #fff;
		border-radius: 0.5rem;
	}
	.empty a {
		color: #1d4ed8;
	}
</style>
