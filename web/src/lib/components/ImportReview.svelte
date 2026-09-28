<script lang="ts">
	import {
		applyImport,
		cancelImport,
		validateImport,
		type ImportPreview
	} from '$lib/api/dataPortability';
	import { pushToast } from '$lib/stores/toast';
	let preview: ImportPreview | null = null;
	let working = false;
	let input: HTMLInputElement;
	async function choose() {
		const file = input.files?.[0];
		if (!file) return;
		working = true;
		try {
			preview = await validateImport(file);
		} catch (error) {
			pushToast(error instanceof Error ? error.message : 'Invalid archive', 'error');
		} finally {
			working = false;
		}
	}
	async function apply() {
		if (!preview) return;
		working = true;
		try {
			await applyImport(preview.importId, []);
			pushToast('Import applied', 'success');
			preview = null;
		} catch (error) {
			pushToast(
				error instanceof Error ? error.message : 'Import failed; existing data was preserved',
				'error'
			);
		} finally {
			working = false;
		}
	}
	async function cancel() {
		if (preview) await cancelImport(preview.importId);
		preview = null;
	}
</script>

<section class="import">
	<label for="archive">Validate an export archive</label><input
		id="archive"
		bind:this={input}
		type="file"
		accept=".zip,application/zip"
		onchange={choose}
		disabled={working}
	/>{#if preview}<div class="review">
			<h2>Import review</h2>
			<p>{preview.conflicts.length} conflicts · expires {preview.expiresAt}</p>
			<button type="button" onclick={apply} disabled={working}>Apply import</button><button
				type="button"
				class="cancel"
				onclick={cancel}>Cancel</button
			>
		</div>{/if}
</section>

<style>
	.import {
		padding: 1rem;
		background: #fff;
		border: 1px solid #e2e8f0;
		border-radius: 0.5rem;
	}
	.import label {
		display: block;
		font-weight: 650;
		margin-bottom: 0.4rem;
	}
	.review {
		margin-top: 1rem;
		padding: 0.8rem;
		background: #eff6ff;
		border-radius: 0.4rem;
	}
	.review button {
		border: 0;
		background: #1d4ed8;
		color: #fff;
		padding: 0.55rem 0.8rem;
		border-radius: 0.35rem;
		margin-right: 0.4rem;
	}
	.review .cancel {
		background: #e2e8f0;
		color: #172033;
	}
</style>
