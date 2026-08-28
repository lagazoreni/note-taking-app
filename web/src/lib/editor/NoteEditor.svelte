<script lang="ts">
	import { api } from '$lib/api/client';
	import { directiveIds, insertDirective, wrapSelection } from './directives';
	import { markSaved, markUnsaved } from '$lib/stores/unsaved';
	import { pushToast } from '$lib/stores/toast';
	import { questionsApi } from '$lib/api/questions';
	import QuestionPicker from '$lib/components/QuestionPicker.svelte';
	import NoteReader from './NoteReader.svelte';
	import type { Note, NoteQuestionLinkWrite } from '$lib/types/note';
	import type { DisplayMode, Question, QuestionKind } from '$lib/types/question';

	export let workspaceId: string;
	export let existing: Note | null = null;
	export let onSave: ((note: Note) => void) | undefined = undefined;

	type EditorMode = 'read' | 'edit';

	let mode: EditorMode = existing ? 'read' : 'edit';
	let title = existing?.title ?? '';
	let bodyMarkdown = existing?.bodyMarkdown ?? '';
	let links: NoteQuestionLinkWrite[] = (existing?.questionLinks ?? []).map((link) => ({
		questionId: link.questionId,
		displayMode: link.displayMode,
		position: link.position
	}));
	let questions: Question[] = (existing?.questionLinks ?? []).map(questionFromLink);
	let questionText = '';
	let questionMode: DisplayMode = 'collapsed';
	let availableQuestions: Question[] = [];
	let addingQuestion = false;
	let saving = false;
	let error = '';
	let lastSavedTitle = existing?.title ?? '';
	let lastSavedBodyMarkdown = existing?.bodyMarkdown ?? '';
	let hydratedNoteId: string | null = null;
	let hydrationRequest = 0;

	$: if (
		(title || bodyMarkdown) &&
		(title !== lastSavedTitle || bodyMarkdown !== lastSavedBodyMarkdown)
	)
		markUnsaved(existing?.id ?? 'new-note');
	$: if (workspaceId)
		questionsApi
			.list({
				workspaceId,
				status: ['unanswered', 'in_progress', 'deferred', 'answered'],
				kind: 'question'
			})
			.then((result) => (availableQuestions = result.items))
			.catch(() => undefined);

	// Summaries returned with a note intentionally do not contain lifecycle fields. Hydrate them
	// before passing questions to the reading card and lifecycle editor.
	$: if (existing && existing.id !== hydratedNoteId) {
		const note = existing;
		hydratedNoteId = note.id;
		mode = 'read';
		title = note.title;
		bodyMarkdown = note.bodyMarkdown;
		lastSavedTitle = note.title;
		lastSavedBodyMarkdown = note.bodyMarkdown;
		const noteLinks = note.questionLinks ?? [];
		links = noteLinks.map((link) => ({
			questionId: link.questionId,
			displayMode: link.displayMode,
			position: link.position
		}));
		questions = noteLinks.map(questionFromLink);
		void hydrateQuestions(note, ++hydrationRequest);
	}

	function questionFromLink(link: Note['questionLinks'][number]): Question {
		return {
			...link.question,
			kind: link.question.kind ?? 'question',
			answerMarkdown: null,
			tagIds: [],
			reminder: null,
			linkedNotes: [],
			createdAt: '',
			updatedAt: '',
			dueDate: link.question.dueDate
		};
	}

	async function hydrateQuestions(note: Note, request: number) {
		try {
			const fullQuestions = await Promise.all(
				(note.questionLinks ?? []).map((link) => questionsApi.get(link.questionId))
			);
			if (request !== hydrationRequest || existing?.id !== note.id) return;
			questions = fullQuestions;
		} catch (cause) {
			if (request !== hydrationRequest || existing?.id !== note.id) return;
			error = cause instanceof Error ? cause.message : 'Could not load linked questions';
		}
	}

	function selectExisting(question: Question) {
		if (links.some((link) => link.questionId === question.id)) {
			error = 'That question is already in this note.';
			return;
		}
		const position = links.length;
		bodyMarkdown = insertDirective(bodyMarkdown, question.id);
		links = [...links, { questionId: question.id, displayMode: questionMode, position }];
		questions = [...questions, question];
		pushToast('Existing question linked', 'success');
	}

	async function addQuestion() {
		if (!questionText.trim()) return;
		addingQuestion = true;
		error = '';
		try {
			const question = await questionsApi.create({
				workspaceId,
				questionText: questionText.trim(),
				kind: 'question',
				status: 'unanswered',
				priority: 'none',
				answerMarkdown: null,
				tagIds: []
			});
			const position = links.length;
			bodyMarkdown = insertDirective(bodyMarkdown, question.id);
			links = [...links, { questionId: question.id, displayMode: questionMode, position }];
			questions = [...questions, question];
			questionText = '';
			pushToast('Question added to this note', 'success');
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not add the question.';
		} finally {
			addingQuestion = false;
		}
	}

	async function save(): Promise<Note | null> {
		saving = true;
		error = '';
		try {
			const ids = directiveIds(bodyMarkdown);
			const questionLinks = ids.map((id, position) => {
				const current = links.find((link) => link.questionId === id);
				return current
					? { ...current, position }
					: { questionId: id, displayMode: 'collapsed' as DisplayMode, position };
			});
			const payload = {
				workspaceId,
				title: title.trim(),
				bodyMarkdown,
				topicId: existing?.topicId ?? null,
				parentNoteId: existing?.parentNoteId ?? null,
				questionLinks,
				tagIds: existing?.tagIds ?? []
			};
			const note = existing
				? await api.put<Note>(`/api/v1/notes/${existing.id}`, {
						...payload,
						version: existing.version
					})
				: await api.post<Note>('/api/v1/notes', payload);
			const unsavedKey = existing?.id ?? 'new-note';
			links = questionLinks;
			title = note.title;
			bodyMarkdown = note.bodyMarkdown;
			lastSavedTitle = title;
			lastSavedBodyMarkdown = bodyMarkdown;
			markSaved(unsavedKey);
			// Keep the authoritative version locally so a subsequent edit does not reuse a stale token.
			existing = note;
			onSave?.(note);
			return note;
		} catch (cause) {
			error =
				cause instanceof Error ? cause.message : 'Could not save note. Your draft is still here.';
			return null;
		} finally {
			saving = false;
		}
	}

	async function captureFromReader(
		kind: QuestionKind,
		passage: string,
		text: string
	): Promise<Question> {
		if (!existing) throw new Error('Save the note before adding a highlight');
		const question = await questionsApi.create({
			workspaceId,
			questionText: text,
			kind,
			status: 'unanswered',
			priority: 'none',
			answerMarkdown: null,
			tagIds: []
		});
		const nextBody = wrapSelection(bodyMarkdown, passage, question.id);
		const nextLinks: NoteQuestionLinkWrite[] = directiveIds(nextBody).map((id, position) => ({
			questionId: id,
			displayMode: 'collapsed',
			position
		}));
		bodyMarkdown = nextBody;
		links = nextLinks;
		questions = [...questions, question];
		const savedNote = await save();
		if (!savedNote) throw new Error(error || 'Could not save highlight');
		return question;
	}

	function toggleMode() {
		if (!existing) return;
		mode = mode === 'read' ? 'edit' : 'read';
	}
</script>

{#if existing}
	<div class="mode-switch">
		<button type="button" onclick={toggleMode}
			>{mode === 'read' ? 'Edit note' : 'Reading view'}</button
		>
	</div>
{/if}

{#if mode === 'read' && existing}
	<NoteReader markdown={bodyMarkdown} {questions} onCapture={captureFromReader} />
{:else}
	<form
		class="editor"
		onsubmit={(event) => {
			event.preventDefault();
			void save();
		}}
	>
		{#if error}<div class="error" role="alert">
				<span>{error}</span><button type="button" onclick={() => void save()}>Retry save</button>
			</div>{/if}
		<label for="note-title">Title</label><input
			id="note-title"
			bind:value={title}
			maxlength="300"
			required
			placeholder="A focused note"
		/>
		<label for="note-body">Note</label><textarea
			id="note-body"
			bind:value={bodyMarkdown}
			rows="14"
			placeholder="Write paragraphs, lists, and questions…"
		></textarea>
		<section class="question-tools" aria-labelledby="question-tools-title">
			<h2 id="question-tools-title">Inline questions</h2>
			<div class="question-create">
				<label for="question-text">Question text</label><input
					id="question-text"
					bind:value={questionText}
					placeholder="What do I need to find out?"
				/><label for="question-mode">Presentation</label><select
					id="question-mode"
					bind:value={questionMode}
					><option value="collapsed">Collapsed</option><option value="expanded">Expanded</option
					><option value="link">Link</option></select
				><button
					type="button"
					onclick={() => void addQuestion()}
					disabled={addingQuestion || !questionText.trim()}
					>{addingQuestion ? 'Adding…' : 'Add question'}</button
				>
			</div>
			<QuestionPicker
				questions={availableQuestions.filter(
					(question) => !links.some((link) => link.questionId === question.id)
				)}
				onSelect={selectExisting}
			/>
			{#if questions.length > 0}<div class="question-preview">
					{#each questions as question, index}<div>
							<span>{index + 1}. {question.questionText}</span><span class="mode"
								>{links.find((link) => link.questionId === question.id)?.displayMode ??
									'collapsed'}</span
							>
						</div>{/each}
				</div>{/if}
		</section>
		<p class="hint">
			Paragraphs and ordered/unordered lists are supported. Question directives are saved atomically
			with the note.
		</p>
		<button class="save" disabled={saving || !title.trim()}
			>{saving ? 'Saving…' : 'Save note'}</button
		>
	</form>
{/if}
{#if error && mode === 'read'}<div class="error" role="alert">
		<span>{error}</span><button type="button" onclick={() => void save()}>Retry save</button>
	</div>{/if}

<style>
	.mode-switch {
		display: flex;
		justify-content: flex-end;
		margin-bottom: 0.55rem;
	}
	.mode-switch button {
		border: 1px solid #bfdbfe;
		background: #eff6ff;
		color: #1d4ed8;
		border-radius: 0.4rem;
		padding: 0.55rem 0.8rem;
	}
	.editor {
		display: grid;
		gap: 0.55rem;
		background: white;
		border: 1px solid #e2e8f0;
		border-radius: 0.65rem;
		padding: 1rem;
	}
	label {
		font-weight: 650;
		margin-top: 0.4rem;
	}
	input,
	textarea,
	select {
		width: 100%;
		padding: 0.7rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
		resize: vertical;
	}
	textarea {
		line-height: 1.5;
	}
	.question-tools {
		margin-top: 0.75rem;
		padding: 0.85rem;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.question-tools h2 {
		margin: 0 0 0.5rem;
		font-size: 1rem;
	}
	.question-create {
		display: grid;
		grid-template-columns: 1fr 9rem auto;
		gap: 0.5rem;
		align-items: end;
	}
	.question-create label {
		grid-row: 1;
		font-size: 0.85rem;
	}
	.question-create input {
		grid-column: 1;
	}
	.question-create select {
		grid-column: 2;
	}
	.question-create button {
		grid-column: 3;
		padding: 0.68rem 0.8rem;
		border: 0;
		background: #2563eb;
		color: white;
		border-radius: 0.4rem;
		white-space: nowrap;
	}
	.question-preview {
		display: grid;
		gap: 0.35rem;
		margin-top: 0.75rem;
	}
	.question-preview div {
		display: flex;
		justify-content: space-between;
		gap: 0.5rem;
		padding: 0.45rem 0.6rem;
		background: white;
		border-radius: 0.35rem;
	}
	.mode {
		color: #2563eb;
		font-size: 0.8rem;
	}
	.hint {
		color: #64748b;
		font-size: 0.9rem;
		margin: 0;
	}
	.save {
		width: fit-content;
		border: 0;
		background: #1d4ed8;
		color: white;
		border-radius: 0.4rem;
		padding: 0.65rem 1rem;
	}
	.save:disabled {
		opacity: 0.55;
	}
	.error {
		display: flex;
		justify-content: space-between;
		gap: 0.5rem;
		padding: 0.7rem;
		color: #991b1b;
		background: #fef2f2;
		border-radius: 0.4rem;
	}
	.error button {
		border: 1px solid #fecaca;
		background: white;
		border-radius: 0.3rem;
	}
	@media (max-width: 640px) {
		.question-create {
			grid-template-columns: 1fr;
		}
		.question-create label,
		.question-create input,
		.question-create select,
		.question-create button {
			grid-column: 1;
		}
	}
</style>
