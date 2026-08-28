<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { notesApi } from '$lib/api/notes';
	import NoteEditor from '$lib/editor/NoteEditor.svelte';
	import DeleteNoteReview from '$lib/components/DeleteNoteReview.svelte';
	import type { Note } from '$lib/types/note';
	import { pushToast } from '$lib/stores/toast';
	let note: Note | null = null;
	let preview: import('$lib/api/notes').DeletionPreview | null = null;
	let noteId = '';
	let loading = true;
	let error = '';
	$: noteId = page.params.noteId ?? '';
	onMount(async () => {
		try {
			note = await notesApi.get(noteId);
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load note';
		} finally {
			loading = false;
		}
	});
	function saved(value: Note) {
		note = value;
		pushToast('Note saved', 'success');
	}
	async function reviewDelete() {
		if (!note) return;
		try {
			preview = await notesApi.previewDeletion(note.id, note.version);
		} catch (cause) {
			pushToast(cause instanceof Error ? cause.message : 'Could not preview deletion', 'error');
		}
	}
	async function confirmDelete(
		questionDecisions: { questionId: string; action: 'keep_unlinked' | 'delete' }[],
		childDecisions: { childNoteId: string; action: 'make_root' | 'move_to_parent' }[]
	) {
		if (!note || !preview) return;
		try {
			await notesApi.deleteReviewed(note.id, {
				previewToken: preview.previewToken,
				version: preview.noteVersion,
				questionDecisions,
				childDecisions
			});
			pushToast('Note deleted', 'success');
			goto('/notes');
		} catch (cause) {
			pushToast(
				cause instanceof Error ? cause.message : 'Deletion failed; no changes were made',
				'error'
			);
		} finally {
			preview = null;
		}
	}
</script>

<svelte:head><title>{note?.title ?? 'Note'} · Noted</title></svelte:head>
{#if loading}<div role="status">Loading note…</div>{:else if error}<div class="error" role="alert">
		<h1>Could not load note</h1>
		<p>{error}</p>
		<a href="/notes">Back to notes</a>
	</div>{:else if note}<section class="page">
		<p class="eyebrow">Note</p>
		<div class="title-row">
			<h1>{note.title}</h1>
			<button class="delete" onclick={reviewDelete}>Delete note</button>
		</div>
		<NoteEditor workspaceId={note.workspaceId} existing={note} onSave={saved} />
	</section>
	{#if preview}<DeleteNoteReview
			{preview}
			onCancel={() => (preview = null)}
			onConfirm={confirmDelete}
		/>{/if}{/if}

<style>
	.page {
		max-width: 60rem;
		margin: auto;
	}
	.title-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
	}
	.delete {
		border: 0;
		background: #fee2e2;
		color: #991b1b;
		padding: 0.55rem 0.8rem;
		border-radius: 0.4rem;
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
	.error {
		padding: 2rem;
		background: #fef2f2;
		border-radius: 0.5rem;
	}
	.error a {
		color: #1d4ed8;
	}
</style>
