<script lang="ts">
	import { workspaceHydrated, workspaceReady } from '$lib/stores/workspace';
</script>

<svelte:head><title>Noted · Home</title></svelte:head>
{#if !$workspaceHydrated}
	<section class="hero" aria-busy="true">
		<p class="eyebrow">Local-first notes</p>
		<h1>Keep questions in context.</h1>
		<p class="lead" role="status">Loading your workspace setup…</p>
	</section>
{:else if !$workspaceReady}
	<section class="hero setup-state">
		<p class="eyebrow">Get started</p>
		<h1>Create your first workspace.</h1>
		<p class="lead">
			A workspace keeps notes and questions in the right context. Create one before starting your
			first note.
		</p>
		<div class="actions">
			<a class="button primary" href="/workspaces">Create workspace</a>
		</div>
	</section>
{:else}
	<section class="hero">
		<p class="eyebrow">Local-first notes</p>
		<h1>Keep questions in context.</h1>
		<p class="lead">
			Capture a thought in a note, follow it centrally, and keep your work available without the
			internet.
		</p>
		<div class="actions">
			<a class="button primary" href="/workspaces">Choose a workspace</a><a
				class="button secondary"
				href="/notes/new">New note</a
			>
		</div>
	</section>
{/if}
<section class="cards" aria-label="Getting started">
	<article>
		<strong>Write in context</strong>
		<p>Use focused Markdown notes with inline shared questions.</p>
	</article>
	<article>
		<strong>Track what is open</strong>
		<p>Active and answered views are filtered views over the same questions.</p>
	</article>
	<article>
		<strong>Stay local</strong>
		<p>The Go service and SQLite keep your data on this device.</p>
	</article>
</section>

<style>
	.hero {
		max-width: 54rem;
		padding: 4rem 0 2rem;
	}
	.eyebrow {
		color: #2563eb;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		font-size: 0.8rem;
	}
	h1 {
		margin: 0.5rem 0;
		font-size: clamp(2.3rem, 8vw, 4.8rem);
		line-height: 1;
		color: #172554;
	}
	.lead {
		max-width: 42rem;
		font-size: 1.2rem;
		color: #475569;
	}
	.actions {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
		margin-top: 1.5rem;
	}
	.button {
		padding: 0.7rem 1rem;
		border-radius: 0.45rem;
		text-decoration: none;
		font-weight: 650;
	}
	.primary {
		background: #1d4ed8;
		color: white;
	}
	.secondary {
		background: white;
		color: #1d4ed8;
		border: 1px solid #93c5fd;
	}
	.cards {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 1rem;
	}
	.cards article {
		background: white;
		border: 1px solid #e2e8f0;
		border-radius: 0.65rem;
		padding: 1rem;
	}
	.cards p {
		color: #64748b;
	}
	@media (max-width: 680px) {
		.hero {
			padding-top: 2rem;
		}
		.cards {
			grid-template-columns: 1fr;
		}
	}
</style>
