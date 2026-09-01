<script lang="ts">
	import type { Question } from '$lib/types/question';
	export let questions: Question[] = [];
	export let emptyMessage = 'No questions match this view.';
	export let answered = false;
	const statusLabel: Record<string, string> = {
		unanswered: 'Unanswered',
		in_progress: 'In progress',
		deferred: 'Deferred',
		answered: 'Answered'
	};
	function overdue(question: Question): boolean {
		return (
			!!question.dueDate &&
			question.status !== 'answered' &&
			question.dueDate < new Date().toISOString().slice(0, 10)
		);
	}
</script>

{#if questions.length === 0}<div class="empty">
		<h2>{emptyMessage}</h2>
		<p>Questions stay connected to their source notes.</p>
	</div>{:else}<div
		class="questions"
		aria-label={answered ? 'Answered questions' : 'Active questions'}
	>
		{#each questions as question (question.id)}<article class="card">
				<div class="meta">
					<span class="status">{statusLabel[question.status] ?? question.status}</span><span
						>{question.priority}</span
					>{#if overdue(question)}<span class="overdue">Overdue · {question.dueDate}</span
						>{:else if question.dueDate}<span>Due {question.dueDate}</span>{/if}
				</div>
				<h2>
					<a
						href={question.linkedNotes[0]
							? `/questions/${question.id}?noteId=${question.linkedNotes[0].id}`
							: `/questions/${question.id}`}
					>{question.questionText}</a>
				</h2>
				{#if question.answerMarkdown}<p class="answer">{question.answerMarkdown}</p>{/if}
				<div class="links">
					{#each question.linkedNotes as note}<a href={`/notes/${note.id}`}>↳ {note.title}</a
						>{:else}<span>Unlinked question</span>{/each}
				</div>
			</article>{/each}
	</div>{/if}

<style>
	.empty {
		padding: 2rem;
		text-align: center;
		background: #fff;
		border: 1px dashed #94a3b8;
		border-radius: 0.6rem;
	}
	.empty p {
		color: #64748b;
	}
	.questions {
		display: grid;
		gap: 0.8rem;
	}
	.card {
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.6rem;
		padding: 1rem;
	}
	.card h2 {
		margin: 0.5rem 0;
		font-size: 1.1rem;
	}
	.card h2 a {
		color: #172554;
	}
	.meta,
	.links {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		color: #64748b;
		font-size: 0.85rem;
	}
	.status {
		color: #1d4ed8;
		font-weight: 700;
	}
	.overdue {
		color: #b91c1c;
		font-weight: 700;
	}
	.answer {
		color: #475569;
		white-space: pre-wrap;
	}
	.links a {
		color: #2563eb;
	}
</style>
