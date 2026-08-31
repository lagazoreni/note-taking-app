<script lang="ts">
	import { tick } from 'svelte';
	import { renderNoteHtml } from '$lib/editor/markdown';

	export let markdown = '';
	export let questionId = '';
	export let ariaLabel = 'Source note';

	$: html = renderNoteHtml(markdown);
	$: html, questionId, scrollToHighlight();

	async function scrollToHighlight() {
		await tick();
		const mark = document.querySelector('mark[data-annotation-id="' + questionId + '"]');
		mark?.scrollIntoView?.({ block: 'center' });
	}
</script>

<article aria-label={ariaLabel}>
	<!-- eslint-disable-next-line svelte/no-at-html-tags -- renderNoteHtml sanitizes the generated HTML with DOMPurify. -->
	{@html html}
</article>

<style>
	article {
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.65rem;
		padding: 1.1rem 1.2rem;
		line-height: 1.6;
		color: #0f172a;
		max-height: min(60vh, 28rem);
		overflow: auto;
	}
	article :global(mark[data-annotation-id]) {
		background: #fef3c7;
		border-bottom: 2px solid #d97706;
		padding: 0 0.1em;
	}
	article :global(mark.annotation-chip) {
		display: inline-block;
		width: 0.65rem;
		height: 0.65rem;
		border-radius: 999px;
		padding: 0;
		vertical-align: middle;
	}
</style>
