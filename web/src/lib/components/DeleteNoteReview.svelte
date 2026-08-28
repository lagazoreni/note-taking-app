<script lang="ts">
	import type { QuestionSummary } from '$lib/types/question';
	interface Preview {
		previewToken: string;
		expiresAt: string;
		noteId: string;
		noteVersion: number;
		singlyLinkedQuestions: QuestionSummary[];
		multiplyLinkedQuestions: QuestionSummary[];
		childNotes: { id: string; title: string; version: number }[];
	}
	export let preview: Preview;
	export let onCancel: (() => void) | undefined = undefined;
	export let onConfirm:
		| ((
				questionDecisions: { questionId: string; action: 'keep_unlinked' | 'delete' }[],
				childDecisions: { childNoteId: string; action: 'make_root' | 'move_to_parent' }[]
		  ) => void)
		| undefined = undefined;
	let questions: Record<string, 'keep_unlinked' | 'delete'> = {};
	let children: Record<string, 'make_root' | 'move_to_parent'> = {};
	$: for (const question of preview.singlyLinkedQuestions)
		questions[question.id] ??= 'keep_unlinked';
	$: for (const child of preview.childNotes) children[child.id] ??= 'make_root';
</script>

<div class="backdrop" role="presentation">
	<div class="dialog" role="dialog" aria-modal="true" aria-labelledby="delete-title">
		<h2 id="delete-title">Review note deletion</h2>
		<p>
			This review is required before anything is changed. Questions linked elsewhere are kept and
			only this note link is removed.
		</p>
		{#if preview.singlyLinkedQuestions.length}<h3>Questions that would become unlinked</h3>
			{#each preview.singlyLinkedQuestions as question}<label class="decision"
					><span>{question.questionText}</span><select bind:value={questions[question.id]}
						><option value="keep_unlinked">Keep question unlinked</option><option value="delete"
							>Delete question</option
						></select
					></label
				>{/each}{/if}{#if preview.multiplyLinkedQuestions.length}<h3>
				Questions kept through another note
			</h3>
			<ul>
				{#each preview.multiplyLinkedQuestions as question}<li>{question.questionText}</li>{/each}
			</ul>{/if}{#if preview.childNotes.length}<h3>Child notes</h3>
			{#each preview.childNotes as child}<label class="decision"
					><span>{child.title}</span><select bind:value={children[child.id]}
						><option value="make_root">Make a root note</option><option value="move_to_parent"
							>Move to deleted note’s parent</option
						></select
					></label
				>{/each}{/if}
		<footer>
			<button type="button" class="secondary" onclick={() => onCancel?.()}>Cancel</button><button
				type="button"
				class="danger"
				onclick={() =>
					onConfirm?.(
						Object.entries(questions).map(([questionId, action]) => ({ questionId, action })),
						Object.entries(children).map(([childNoteId, action]) => ({ childNoteId, action }))
					)}>Delete after review</button
			>
		</footer>
	</div>
</div>

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		background: #0f172a99;
		display: grid;
		place-items: center;
		padding: 1rem;
		z-index: 30;
	}
	.dialog {
		max-width: 38rem;
		width: 100%;
		max-height: 90vh;
		overflow: auto;
		background: #fff;
		border-radius: 0.65rem;
		padding: 1.2rem;
		box-shadow: 0 20px 50px #0f172a55;
	}
	.dialog h2 {
		color: #172554;
	}
	.dialog p {
		color: #475569;
	}
	.decision {
		display: grid;
		grid-template-columns: 1fr 14rem;
		gap: 0.75rem;
		align-items: center;
		padding: 0.5rem 0;
		border-top: 1px solid #e2e8f0;
	}
	.decision select {
		padding: 0.45rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
	}
	h3 {
		font-size: 1rem;
		margin-top: 1.2rem;
	}
	footer {
		display: flex;
		justify-content: flex-end;
		gap: 0.6rem;
		margin-top: 1.2rem;
	}
	.secondary,
	.danger {
		border: 0;
		padding: 0.6rem 0.85rem;
		border-radius: 0.4rem;
	}
	.secondary {
		background: #e2e8f0;
	}
	.danger {
		background: #b91c1c;
		color: #fff;
	}
	@media (max-width: 600px) {
		.decision {
			grid-template-columns: 1fr;
		}
		.dialog {
			padding: 0.9rem;
		}
	}
</style>
