<script lang="ts">
	import type { Question } from '$lib/types/question';
	export let questions: Question[] = [];
	export let onSelect: ((question: Question) => void) | undefined = undefined;
	let query = '';
	$: matches = (questions ?? []).filter((question) =>
		question.questionText.toLowerCase().includes(query.toLowerCase())
	);
</script>

<section class="picker" aria-label="Find an existing question">
	<label for="question-search">Link an existing question</label><input
		id="question-search"
		bind:value={query}
		placeholder="Search question text"
	/>{#if query && matches.length === 0}<p class="muted">No matching questions.</p>{:else}<ul>
			{#each matches as question}<li>
					<button type="button" onclick={() => onSelect?.(question)}
						><strong>{question.questionText}</strong><span
							>{(question.linkedNotes ?? []).map((note) => note.title).join(', ') || 'Unlinked'} · {question.status}</span
						></button
					>
				</li>{/each}
		</ul>{/if}
</section>

<style>
	.picker {
		padding: 0.8rem;
		background: #f8fafc;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.picker label {
		display: block;
		font-weight: 650;
		margin-bottom: 0.35rem;
	}
	input {
		width: 100%;
		padding: 0.6rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
	}
	ul {
		list-style: none;
		padding: 0;
		margin: 0.6rem 0 0;
		display: grid;
		gap: 0.35rem;
		max-height: 16rem;
		overflow: auto;
	}
	button {
		width: 100%;
		text-align: left;
		padding: 0.55rem;
		background: white;
		border: 1px solid #cbd5e1;
		border-radius: 0.35rem;
	}
	button span {
		display: block;
		color: #64748b;
		font-size: 0.8rem;
		margin-top: 0.2rem;
	}
	.muted {
		color: #64748b;
	}
</style>
