<script lang="ts">
	import NoteExcerpt from '$lib/components/NoteExcerpt.svelte';
	import QuestionLifecycle from '$lib/components/QuestionLifecycle.svelte';
	import QuestionSchedule from '$lib/components/QuestionSchedule.svelte';
	import { notesApi } from '$lib/api/notes';
	import { directiveIds, insertAnswerAfterDirective } from '$lib/editor/directives';
	import type { Note, NoteQuestionLinkWrite } from '$lib/types/note';
	import type { Question } from '$lib/types/question';

	export let question: Question;
	export let note: Note | null = null;
	export let selectedNoteId = '';
	export let onSelectNote: (id: string) => void;
	export let onSave: (q: Question) => void;
	let insertingAnswer = false;
	let insertError = '';
	$: canInsertAnswer =
		question.kind !== 'annotation' &&
		question.status === 'answered' &&
		Boolean(question.answerMarkdown) &&
		note !== null;

	function selectNote(event: Event) {
		const target = event.currentTarget as HTMLSelectElement;
		onSelectNote(target.value);
	}

	async function insertAnswer() {
		const currentNote = note;
		const answer = question.answerMarkdown;
		if (
			!currentNote ||
			!answer ||
			question.kind === 'annotation' ||
			question.status !== 'answered'
		) {
			return;
		}

		insertingAnswer = true;
		insertError = '';
		try {
			const bodyMarkdown = insertAnswerAfterDirective(currentNote.bodyMarkdown, question.id, answer);
			const questionLinks: NoteQuestionLinkWrite[] = directiveIds(bodyMarkdown).map(
				(id, position) => {
					const current = currentNote.questionLinks.find((link) => link.questionId === id);
					return current
						? { questionId: current.questionId, displayMode: current.displayMode, position }
						: { questionId: id, displayMode: 'collapsed', position };
				}
			);
			const updatedNote = await notesApi.update(currentNote.id, {
				workspaceId: currentNote.workspaceId,
				topicId: currentNote.topicId,
				parentNoteId: currentNote.parentNoteId,
				title: currentNote.title,
				bodyMarkdown,
				questionLinks,
				tagIds: [...(currentNote.tagIds ?? [])],
				version: currentNote.version
			});
			if (note?.id === currentNote.id) note = updatedNote;
		} catch (cause) {
			insertError = cause instanceof Error ? cause.message : 'Could not insert answer into note';
		} finally {
			insertingAnswer = false;
		}
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
			<NoteExcerpt
				markdown={note.bodyMarkdown}
				questionId={question.id}
				status={question.status}
				ariaLabel="Note excerpt"
			/>
			{#if canInsertAnswer}
				<button type="button" onclick={() => void insertAnswer()} disabled={insertingAnswer}>
					{insertingAnswer ? 'Inserting…' : 'Insert answer into note'}
				</button>
			{/if}
		{:else}
			<p>This question is currently unlinked.</p>
		{/if}
		{#if insertError}<p class="insert-error" role="alert">{insertError}</p>{/if}
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
	.note-pane button {
		width: max-content;
		border: 0;
		background: #2563eb;
		color: #fff;
		border-radius: 0.4rem;
		padding: 0.55rem 0.8rem;
	}
	.note-pane button:disabled {
		opacity: 0.55;
	}
	.insert-error {
		margin: 0;
		padding: 0.6rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.4rem;
	}
	@media (max-width: 800px) {
		.context {
			grid-template-columns: 1fr;
		}
	}
</style>
