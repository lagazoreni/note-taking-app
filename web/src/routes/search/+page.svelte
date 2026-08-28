<script lang="ts">
	import { searchApi, type SearchResult } from '$lib/api/search';
	import { currentWorkspaceId } from '$lib/stores/workspace';
	let query = '';
	let contentScope: 'notes' | 'questions' | 'answers' | 'everything' = 'everything';
	let workspaceScope: 'current' | 'all' = 'current';
	let results: SearchResult[] = [];
	let error = '';
	let loading = false;
	async function search() {
		if (!query.trim()) return;
		loading = true;
		error = '';
		try {
			results = (
				await searchApi({
					q: query,
					contentScope,
					workspaceScope,
					workspaceId: $currentWorkspaceId ?? undefined
				})
			).items;
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not search';
		} finally {
			loading = false;
		}
	}
</script>

<section>
	<p class="eyebrow">Discovery</p>
	<h1>Search</h1>
	<form
		onsubmit={(event) => {
			event.preventDefault();
			search();
		}}
		class="search"
	>
		<label for="q">Search notes, questions, and answers</label>
		<div>
			<input id="q" bind:value={query} placeholder="Try a phrase" /><select
				aria-label="Content scope"
				bind:value={contentScope}
				><option value="everything">Everything</option><option value="notes">Notes</option><option
					value="questions">Questions</option
				><option value="answers">Answers</option></select
			><select aria-label="Workspace scope" bind:value={workspaceScope}
				><option value="current">Current workspace</option><option value="all"
					>All workspaces</option
				></select
			><button class="button" disabled={loading}>{loading ? 'Searching…' : 'Search'}</button>
		</div>
	</form>
	{#if error}<div class="error" role="alert">
			{error}
		</div>{:else if results.length === 0 && query}<div class="empty">
			No matching content.
		</div>{:else}<div class="results">
			{#each results as result}<article>
					<div class="meta"><span>{result.type}</span><span>{result.workspaceName}</span></div>
					<h2>{result.title ?? result.type}</h2>
					<p>{result.snippet}</p>
					<a href={result.destination}>Open result</a>
				</article>{/each}
		</div>{/if}
</section>

<style>
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
	.search {
		background: #fff;
		padding: 1rem;
		border: 1px solid #e2e8f0;
		border-radius: 0.6rem;
	}
	.search label {
		display: block;
		font-weight: 650;
		margin-bottom: 0.4rem;
	}
	.search div {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}
	input {
		flex: 1;
		min-width: 12rem;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	select {
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	.button {
		border: 0;
		background: #1d4ed8;
		color: #fff;
		padding: 0.6rem 1rem;
		border-radius: 0.4rem;
	}
	.empty {
		margin-top: 1rem;
		color: #64748b;
	}
	.error {
		margin-top: 1rem;
		padding: 1rem;
		background: #fef2f2;
		color: #991b1b;
		border-radius: 0.5rem;
	}
	.results {
		display: grid;
		gap: 0.75rem;
		margin-top: 1rem;
	}
	.results article {
		padding: 1rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.results h2 {
		font-size: 1.05rem;
		margin: 0.4rem 0;
	}
	.results p {
		color: #475569;
	}
	.results a {
		color: #1d4ed8;
	}
	.meta {
		display: flex;
		gap: 0.6rem;
		color: #64748b;
		font-size: 0.8rem;
		text-transform: capitalize;
	}
</style>
