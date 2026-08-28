<script lang="ts">
	import type { Note } from '$lib/types/note';
	export let notes: Note[] = [];
	$: children = (parent: string | null) =>
		notes.filter((note) => (note.parentNoteId ?? null) === parent);
</script>

<nav aria-label="Note hierarchy">
	<ul>
		{#each children(null) as note}<li>
				<a href={`/notes/${note.id}`}>{note.title}</a>{#if children(note.id).length}<ul>
						{#each children(note.id) as child}<li>
								<a href={`/notes/${child.id}`}>{child.title}</a>
							</li>{/each}
					</ul>{/if}
			</li>{:else}<li class="muted">No notes yet.</li>{/each}
	</ul>
</nav>

<style>
	ul {
		list-style: none;
		padding-left: 0;
	}
	ul ul {
		padding-left: 1rem;
		border-left: 1px solid #cbd5e1;
	}
	li {
		margin: 0.35rem 0;
	}
	a {
		color: #1d4ed8;
	}
	.muted {
		color: #64748b;
	}
</style>
