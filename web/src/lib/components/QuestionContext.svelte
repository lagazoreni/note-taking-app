<script lang="ts">
	import NoteExcerpt from '$lib/components/NoteExcerpt.svelte';
	import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
	import QuestionSchedule from '$lib/components/QuestionSchedule.svelte';
	import type { Note } from '$lib/types/note';
	import type { Question } from '$lib/types/question';

	export let question: Question;
	export let note: Note | null = null;
	export let selectedNoteId = '';
	export let onSelectNote: (id: string) => void;
	export let onSave: (q: Question) => void;

	function selectNote(event: Event) {
		const target = event.currentTarget as HTMLSelectElement;
		onSelectNote(target.value);
	}
</script>

<section class="context">
	<div class="note-pane">
		{#if question.linkedNotes.length > 1}
			<label for="context-note">Source note</label>
			<select id="context-note" value={selectedNoteId} onchange={selectNote}>
				{#each question.linkedNotes as linked}
					<option value={linked.id}>{linked.title}</option>
				{/each}
			</select>
		{/if}
		{#if note}
			<h2>{note.title}</h2>
			<NoteExcerpt markdown={note.bodyMarkdown} questionId={question.id} ariaLabel="Note excerpt" />
		{:else}
			<p>This question is currently unlinked.</p>
		{/if}
	</div>
	<div class="controls">
		<QuestionLifecycle {question} {onSave} />
		<QuestionSchedule {question} {onSave} />
	</div>
</section>

<style>
	.context {
		display: grid;
		grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
		gap: 1.25rem;
		align-items: start;
	}
	.note-pane,
	.controls {
		display: grid;
		gap: 0.65rem;
	}
	.note-pane h2 {
		margin: 0;
		color: #172554;
		font-size: 1.15rem;
	}
	.note-pane p {
		margin: 0;
		color: #334155;
	}
	.note-pane label {
		font-weight: 650;
	}
	.note-pane select {
		width: max-content;
		max-width: 100%;
		padding: 0.5rem 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	@media (max-width: 800px) {
		.context {
			grid-template-columns: 1fr;
		}
	}
</style>
