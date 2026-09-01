<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { questionsApi } from '$lib/api/questions';
	import { notesApi } from '$lib/api/notes';
	import type { Note } from '$lib/types/note';
	import type { Question } from '$lib/types/question';
	import QuestionContext from '$lib/components/QuestionContext.svelte';
	let question: Question | null = null;
	let note: Note | null = null;
	let selectedNoteId = '';
	let questionId = '';
	let loading = true;
	let error = '';
	$: questionId = page.params.questionId ?? '';
	async function load() {
		loading = true;
		error = '';
		try {
			const loadedQuestion = question ?? (await questionsApi.get(questionId));
			question = loadedQuestion;
			const requestedNoteId = page.url.searchParams.get('noteId');
			const requestedLinkedNoteId = requestedNoteId && loadedQuestion.linkedNotes.some((linked) => linked.id === requestedNoteId)
				? requestedNoteId
				: undefined;
			const selectedLinkedNoteId = selectedNoteId && loadedQuestion.linkedNotes.some((linked) => linked.id === selectedNoteId)
				? selectedNoteId
				: undefined;
			const noteId = selectedLinkedNoteId ?? requestedLinkedNoteId ?? loadedQuestion.linkedNotes[0]?.id;
			selectedNoteId = noteId ?? '';
			note = noteId ? await notesApi.get(noteId) : null;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load question';
		} finally {
			loading = false;
		}
	}
	onMount(() => {
		void load();
	});
	function saved(value: Question) {
		question = value;
	}
	async function selectNote(id: string) {
		selectedNoteId = id;
		note = null;
		error = '';
		try {
			note = await notesApi.get(id);
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load note';
		}
	}
	function retry() {
		void load();
	}
</script>

{#if loading}<div role="status">Loading question…</div>{:else if error}<div
		role="alert"
		class="error"
	>
		{error}
		<button type="button" onclick={retry}>Retry</button>
	</div>{:else if question}<section class="page">
		<p class="eyebrow">Question</p>
		<h1>{question.questionText}</h1>
		<QuestionContext
			{question}
			{note}
			selectedNoteId={selectedNoteId}
			onSelectNote={selectNote}
			onSave={saved}
		/>
	</section>{/if}

<style>
	.page {
		max-width: 58rem;
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
	.linked {
		margin-top: 1.5rem;
		padding: 1rem;
		background: #fff;
		border-radius: 0.5rem;
		border: 1px solid #e2e8f0;
	}
	.linked a {
		display: block;
		color: #1d4ed8;
		margin: 0.3rem 0;
	}
	.error {
		padding: 1rem;
		background: #fef2f2;
		color: #991b1b;
	}
</style>
