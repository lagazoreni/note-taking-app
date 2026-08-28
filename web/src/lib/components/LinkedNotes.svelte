<script lang="ts">
	import type { LinkedNoteSummary } from '$lib/types/question';
	export let notes: LinkedNoteSummary[] = [];
</script>

<section aria-label="Linked notes">
	<h3>Linked notes</h3>
	{#if notes.length === 0}<p class="muted">This question is currently unlinked.</p>{:else}<ul>
			{#each notes as note}<li>
					<a href={`/notes/${note.id}`}>{note.title}</a><span>{note.displayMode}</span>
				</li>{/each}
		</ul>{/if}
</section>

<style>
	section {
		margin-top: 1rem;
	}
	h3 {
		font-size: 1rem;
	}
	ul {
		padding: 0;
		list-style: none;
		display: grid;
		gap: 0.35rem;
	}
	li {
		display: flex;
		justify-content: space-between;
		gap: 0.5rem;
	}
	a {
		color: #1d4ed8;
	}
	span,
	.muted {
		color: #64748b;
		font-size: 0.85rem;
	}
</style>
