import { api, request } from './client';
export interface ImportPreview {
	importId: string;
	expiresAt: string;
	formatVersion: number;
	archiveSha256: string;
	counts: Record<string, number>;
	conflicts: {
		conflictId: string;
		entityType: string;
		entityId: string;
		reason: string;
		allowedResolutions: string[];
	}[];
}
export async function exportData(): Promise<void> {
	const response = await fetch('/api/v1/export', { method: 'POST' });
	if (!response.ok) throw new Error('Could not export data');
	const blob = await response.blob();
	const url = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = url;
	link.download = 'noted-export.zip';
	link.click();
	URL.revokeObjectURL(url);
}
export async function validateImport(file: File): Promise<ImportPreview> {
	const body = new FormData();
	body.set('archive', file);
	return request<ImportPreview>('/api/v1/imports/validate', { method: 'POST', body });
}
export async function applyImport(
	importId: string,
	conflictResolutions: unknown[] = []
): Promise<unknown> {
	return api.post(`/api/v1/imports/${importId}/apply`, { conflictResolutions });
}
export async function cancelImport(importId: string): Promise<void> {
	await api.delete(`/api/v1/imports/${importId}`);
}
