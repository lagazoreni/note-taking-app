<script lang="ts">
	import { onMount } from 'svelte';
	import { searchApi, type SearchResult } from '$lib/api/search';
	import { tagsApi } from '$lib/api/tags';
	import { currentWorkspaceId } from '$lib/stores/workspace';
	import type { Tag } from '$lib/types/workspace';

	let query = '';
	let contentScope: 'notes' | 'questions' | 'answers' | 'everything' = 'everything';
	let workspaceScope: 'current' | 'all' = 'current';
	let tagId = '';
	let results: SearchResult[] = [];
	let availableTags: Tag[] = [];
	let error = '';
	let tagError = '';
	let loading = false;
	let tagsLoading = false;
	let workspaceId: string | null = null;
	let tagsForWorkspace: string | null = null;
	let searchRequest = 0;
	let tagRequest = 0;

	function readUrlState() {
		const params = new URLSearchParams(window.location.search);
		query = params.get('q') ?? '';
		const content = params.get('contentScope');
		if (
			content === 'notes' ||
			content === 'questions' ||
			content === 'answers' ||
			content === 'everything'
		) {
			contentScope = content;
		}
		const scope = params.get('workspaceScope');
		if (scope === 'current' || scope === 'all') workspaceScope = scope;
		tagId = params.get('tagId') ?? '';
	}

	function syncUrlState() {
		if (typeof window === 'undefined') return;
		// This is a one-shot URL serialization for browser history, not reactive component state.
		// eslint-disable-next-line svelte/prefer-svelte-reactivity
		const params = new URLSearchParams(window.location.search);
		if (query.trim()) params.set('q', query.trim());
		else params.delete('q');
		params.set('contentScope', contentScope);
		params.set('workspaceScope', workspaceScope);
		if (tagId) params.set('tagId', tagId);
		else params.delete('tagId');
		const encoded = params.toString();
		window.history.replaceState(
			{},
			'',
			`${window.location.pathname}${encoded ? `?${encoded}` : ''}`
		);
	}

	async function loadTags(id: string) {
		const request = ++tagRequest;
		tagsLoading = true;
		tagError = '';
		availableTags = [];
		try {
			const response = await tagsApi.list(id);
			if (request === tagRequest && workspaceId === id) availableTags = response.items;
		} catch (cause) {
			if (request === tagRequest && workspaceId === id) {
				tagError = cause instanceof Error ? cause.message : 'Could not load workspace tags';
				availableTags = [];
			}
		} finally {
			if (request === tagRequest) tagsLoading = false;
		}
	}

	async function search() {
		syncUrlState();
		if (!query.trim()) {
			results = [];
			error = '';
			return;
		}
		const request = ++searchRequest;
		loading = true;
		error = '';
		try {
			const response = await searchApi({
				q: query.trim(),
				contentScope,
				workspaceScope,
				workspaceId: $currentWorkspaceId ?? undefined,
				tagId: tagId || undefined
			});
			if (request === searchRequest) results = response.items;
		} catch (cause) {
			if (request === searchRequest) {
				results = [];
				error = cause instanceof Error ? cause.message : 'Could not search';
			}
		} finally {
			if (request === searchRequest) loading = false;
		}
	}

	function filterChanged() {
		void search();
	}

	onMount(() => {
		readUrlState();
		const unsubscribe = currentWorkspaceId.subscribe((id) => {
			workspaceId = id;
			if (!id) {
				tagsForWorkspace = null;
				availableTags = [];
				return;
			}
			if (id !== tagsForWorkspace) {
				tagsForWorkspace = id;
				void loadTags(id);
			}
		});
		if (query.trim()) void search();
		return unsubscribe;
	});
</script>

<section>
	<p class="eyebrow">Discovery</p>
	<h1>Search</h1>
	<form
		onsubmit={(event) => {
			event.preventDefault();
			void search();
		}}
		class="search"
	>
		<label for="q">Search notes, questions, and answers</label>
		<div class="search-controls">
			<input id="q" bind:value={query} placeholder="Try a phrase" />
			<label class="sr-only" for="content-scope">Content scope</label>
			<select
				id="content-scope"
				aria-label="Content scope"
				bind:value={contentScope}
				onchange={filterChanged}
			>
				<option value="everything">Everything</option>
				<option value="notes">Notes</option>
				<option value="questions">Questions</option>
				<option value="answers">Answers</option>
			</select>
			<label class="sr-only" for="workspace-scope">Workspace scope</label>
			<select
				id="workspace-scope"
				aria-label="Workspace scope"
				bind:value={workspaceScope}
				onchange={filterChanged}
			>
				<option value="current">Current workspace</option>
				<option value="all">All workspaces</option>
			</select>
			<div class="tag-control">
				<label for="tag-filter">Tag</label>
				<select id="tag-filter" bind:value={tagId} onchange={filterChanged}>
					<option value="">All note tags</option>
					{#each availableTags as tag}
						<option value={tag.id}
							>{tag.name}{tag.availableWorkspaceIds.length > 1 ? ' · shared' : ''}</option
						>
					{/each}
				</select>
			</div>
			<button class="button" disabled={loading}>
				{loading ? 'Searching…' : 'Search'}
			</button>
		</div>
		{#if tagsLoading}
			<p class="tag-status" role="status">Loading note tags…</p>
		{:else if tagError}
			<p class="tag-status error-inline" role="status">{tagError}</p>
		{:else if availableTags.length === 0}
			<p class="tag-status">No workspace tags available. You can search all notes.</p>
		{/if}
	</form>
	{#if loading}
		<div role="status">Searching…</div>
	{:else if error}
		<div class="error" role="alert">{error}</div>
	{:else if results.length === 0 && query}
		<div class="empty">No matching content.</div>
	{:else if !query}
		<div class="empty">Enter a phrase to search your notes.</div>
	{:else}
		<div class="results">
			{#each results as result}
				<article>
					<div class="meta"><span>{result.type}</span><span>{result.workspaceName}</span></div>
					<h2>{result.title ?? result.type}</h2>
					<p>{result.snippet}</p>
					<a href={result.destination}>Open result</a>
				</article>
			{/each}
		</div>
	{/if}
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
	.search > label {
		display: block;
		font-weight: 650;
		margin-bottom: 0.4rem;
	}
	.search-controls {
		display: flex;
		gap: 0.5rem;
		align-items: end;
		flex-wrap: wrap;
	}
	input {
		flex: 1 1 12rem;
		min-width: 12rem;
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
	}
	select {
		padding: 0.65rem;
		border: 1px solid #cbd5e1;
		border-radius: 0.4rem;
		background: white;
	}
	.tag-control {
		display: grid;
		gap: 0.25rem;
		min-width: 10rem;
	}
	.tag-control label {
		font-size: 0.8rem;
		font-weight: 650;
		color: #475569;
	}
	.button {
		border: 0;
		background: #1d4ed8;
		color: #fff;
		padding: 0.6rem 1rem;
		border-radius: 0.4rem;
	}
	.button:disabled {
		opacity: 0.55;
	}
	.tag-status {
		margin: 0.65rem 0 0;
		color: #64748b;
		font-size: 0.85rem;
	}
	.error-inline {
		color: #991b1b;
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
	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		padding: 0;
		margin: -1px;
		overflow: hidden;
		clip: rect(0, 0, 0, 0);
		white-space: nowrap;
		border: 0;
	}
	@media (max-width: 640px) {
		.search-controls > input,
		.search-controls > select,
		.tag-control,
		.tag-control select,
		.button {
			width: 100%;
		}
		.search-controls > input {
			min-width: 0;
		}
	}
</style>
