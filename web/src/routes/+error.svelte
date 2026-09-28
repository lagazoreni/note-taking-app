<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head><title>{page.status} · Noted</title></svelte:head>
<section class="error" aria-live="assertive">
	<p class="eyebrow">Something went wrong</p>
	<h1>{page.status}</h1>
	<p>{page.error?.message ?? 'The requested page could not be loaded.'}</p>
	<a class="button" href="/">Return home</a>
</section>

<style>
	.error {
		max-width: 42rem;
		margin: 5rem auto;
		padding: 2rem;
		text-align: center;
	}
	.eyebrow {
		color: #64748b;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		font-size: 0.8rem;
	}
	.button {
		display: inline-block;
		margin-top: 1rem;
		padding: 0.6rem 1rem;
		border-radius: 0.45rem;
		background: #1d4ed8;
		color: white;
		text-decoration: none;
	}
</style>
