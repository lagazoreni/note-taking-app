<script lang="ts">
	import { page } from '$app/state';
	import { currentWorkspaceId, workspaceReady } from '$lib/stores/workspace';
	import WorkspaceSwitcher from './WorkspaceSwitcher.svelte';

	$: workspace =
		$workspaceReady && $currentWorkspaceId
			? `?workspaceId=${encodeURIComponent($currentWorkspaceId)}`
			: '';
</script>

<header class="topbar">
	<a class="brand" href="/" aria-label="Application home">Noted</a>
	<WorkspaceSwitcher />
	<nav aria-label="Primary navigation">
		{#if $workspaceReady}
			<a class:active={page.url.pathname.startsWith('/notes')} href={`/notes${workspace}`}>Notes</a>
			<a class:active={page.url.pathname === '/questions'} href={`/questions${workspace}`}
				>Active Questions</a
			>
			<a class:active={page.url.pathname.startsWith('/answered')} href={`/answered${workspace}`}
				>Answered</a
			>
		{:else}
			<a class:active={page.url.pathname.startsWith('/workspaces')} href="/workspaces"
				>Set up workspace</a
			>
		{/if}
		<a class:active={page.url.pathname.startsWith('/search')} href="/search">Search</a>
		<a class:active={page.url.pathname.startsWith('/settings')} href="/settings/data">Data</a>
	</nav>
</header>

<style>
	.topbar {
		display: flex;
		align-items: center;
		gap: 1rem;
		flex-wrap: wrap;
		padding: 0.75rem 1rem;
		background: #172554;
		color: white;
	}
	.brand {
		color: white;
		font-weight: 800;
		font-size: 1.2rem;
		text-decoration: none;
		margin-right: 0.5rem;
	}
	nav {
		display: flex;
		gap: 0.35rem;
		flex-wrap: wrap;
	}
	nav a {
		color: #dbeafe;
		padding: 0.4rem 0.65rem;
		border-radius: 0.4rem;
		text-decoration: none;
	}
	nav a:hover,
	nav a.active {
		color: white;
		background: #1e40af;
	}
	@media (max-width: 620px) {
		.topbar {
			align-items: flex-start;
		}
		nav {
			width: 100%;
		}
		nav a {
			flex: 1 0 auto;
			text-align: center;
		}
	}
</style>
