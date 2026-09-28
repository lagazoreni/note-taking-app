<script lang="ts">
	import ExportData from '$lib/components/ExportData.svelte';
	import ImportReview from '$lib/components/ImportReview.svelte';
</script>

<section>
	<p class="eyebrow">Portability</p>
	<h1>Data</h1>
	<div class="empty">
		<h2>Export and restore</h2>
		<p>Archives are validated before any data changes.</p>
		<div class="actions"><ExportData /><ImportReview /></div>
	</div>
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
	.empty {
		max-width: 42rem;
		margin-top: 1.5rem;
		padding: 1.5rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.6rem;
	}
	.actions {
		display: grid;
		gap: 1rem;
		margin-top: 1rem;
	}
</style>
