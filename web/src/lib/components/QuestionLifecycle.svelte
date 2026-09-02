<script lang="ts">
	import { questionsApi } from '$lib/api/questions';
	import { pushToast } from '$lib/stores/toast';
	import { invalidate } from '$lib/stores/query';
	import type { Question, QuestionStatus } from '$lib/types/question';
	export let question: Question;
	export let onSave: ((question: Question) => void) | undefined = undefined;
	let answer = question.answerMarkdown ?? '';
	let status: QuestionStatus = question.status;
	let dueDate = question.dueDate ?? '';
	let saving = false;
	let error = '';
	async function save() {
		saving = true;
		error = '';
		if (status === 'answered' && !answer.trim()) {
			error = 'An answer is required before marking a question answered.';
			saving = false;
			return;
		}
		if (question.kind !== 'annotation' && status === 'deferred' && !dueDate.trim()) {
			error = 'A resume date is required to defer a question.';
			saving = false;
			return;
		}
		try {
			const value = await questionsApi.update(question.id, {
				workspaceId: question.workspaceId,
				questionText: question.questionText,
				answerMarkdown: answer.trim() || null,
				status,
				priority: question.priority,
				dueDate: dueDate.trim() || null,
				tagIds: question.tagIds,
				version: question.version
			});
			question = value;
			answer = value.answerMarkdown ?? '';
			dueDate = value.dueDate ?? '';
			onSave?.(value);
			invalidate('questions', `question:${value.id}`);
			pushToast('Question updated', 'success');
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not save question';
		} finally {
			saving = false;
		}
	}
	function reopen() {
		status = 'in_progress';
	}
</script>

<form
	class="lifecycle"
	onsubmit={(event) => {
		event.preventDefault();
		save();
	}}
>
	<label for="answer">Answer</label><textarea
		id="answer"
		bind:value={answer}
		rows="8"
		placeholder="Record what you found…"
	></textarea><label for="status">Status</label><select id="status" bind:value={status}
		><option value="unanswered">Unanswered</option><option value="in_progress">In progress</option
		><option value="deferred">Deferred</option><option value="answered">Answered</option></select
	>{#if status === 'deferred'}<label for="resume-date">Resume date</label><input
		id="resume-date"
		type="date"
		bind:value={dueDate}
	/>{/if}{#if question.status === 'answered'}<button type="button" class="secondary" onclick={reopen}
			>Reopen in progress</button
		>{/if}{#if error}<div class="error" role="alert">{error}</div>{/if}<button
		class="save"
		disabled={saving}>{saving ? 'Saving…' : 'Save question'}</button
	>
</form>

<style>
	.lifecycle {
		display: grid;
		gap: 0.55rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		padding: 1rem;
		border-radius: 0.6rem;
	}
	label {
		font-weight: 650;
	}
	textarea,
	select,
	input {
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	.save,
	.secondary {
		width: max-content;
		padding: 0.6rem 0.9rem;
		border: 0;
		border-radius: 0.4rem;
	}
	.save {
		background: #1d4ed8;
		color: #fff;
	}
	.secondary {
		background: #dbeafe;
		color: #1e40af;
	}
	.error {
		padding: 0.6rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.4rem;
	}
</style>
