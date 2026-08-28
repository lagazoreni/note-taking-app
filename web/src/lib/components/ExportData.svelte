<script lang="ts">
	import { exportData } from '$lib/api/dataPortability';
	import { pushToast } from '$lib/stores/toast';
	let working = false;
	async function run() {
		working = true;
		try {
			await exportData();
			pushToast('Export downloaded', 'success');
		} catch (error) {
			pushToast(error instanceof Error ? error.message : 'Export failed', 'error');
		} finally {
			working = false;
		}
	}
</script>

<button type="button" onclick={run} disabled={working}
	>{working ? 'Preparing export…' : 'Export all data'}</button
>

<style>
	button {
		border: 0;
		background: #1d4ed8;
		color: #fff;
		padding: 0.65rem 1rem;
		border-radius: 0.4rem;
	}
	button:disabled {
		opacity: 0.55;
	}
</style>
