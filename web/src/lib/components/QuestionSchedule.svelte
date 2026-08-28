<script lang="ts">
	import { questionsApi } from '$lib/api/questions';
	import { pushToast } from '$lib/stores/toast';
	import type { Question } from '$lib/types/question';
	export let question: Question;
	export let onSave: ((question: Question) => void) | undefined;
	let dueDate = question.dueDate ?? '';
	let saving = false;
	async function save() {
		saving = true;
		try {
			const value = await questionsApi.update(question.id, {
				workspaceId: question.workspaceId,
				questionText: question.questionText,
				answerMarkdown: question.answerMarkdown,
				status: question.status,
				priority: question.priority,
				dueDate: dueDate || null,
				tagIds: question.tagIds,
				version: question.version
			});
			question = value;
			onSave?.(value);
			pushToast('Schedule updated', 'success');
		} catch (error) {
			pushToast(error instanceof Error ? error.message : 'Could not update schedule', 'error');
		} finally {
			saving = false;
		}
	}
</script>

<div class="schedule">
	<label for="due-date">Due date</label><input
		id="due-date"
		type="date"
		bind:value={dueDate}
	/><button type="button" onclick={save} disabled={saving}
		>{saving ? 'Saving…' : 'Save date'}</button
	>{#if question.dueDate && question.status !== 'answered' && question.dueDate < new Date()
				.toISOString()
				.slice(0, 10)}<strong class="overdue">Overdue</strong>{/if}
</div>

<style>
	.schedule {
		display: flex;
		align-items: end;
		gap: 0.5rem;
		flex-wrap: wrap;
		padding: 0.75rem 0;
	}
	.schedule label {
		font-size: 0.85rem;
		font-weight: 650;
	}
	.schedule input {
		padding: 0.5rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
	}
	.schedule button {
		border: 0;
		background: #dbeafe;
		color: #1e40af;
		padding: 0.5rem 0.75rem;
		border-radius: 0.35rem;
	}
	.overdue {
		color: #b91c1c;
		font-size: 0.85rem;
	}
</style>
