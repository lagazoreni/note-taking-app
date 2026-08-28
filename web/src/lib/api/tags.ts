import { api } from './client';
import type { Tag } from '$lib/types/workspace';

export interface TagPage {
	items: Tag[];
}

export const tagsApi = {
	list: async (workspaceId: string, signal?: AbortSignal): Promise<TagPage> => {
		const result = await api.get<TagPage>(
			`/api/v1/tags?workspaceId=${encodeURIComponent(workspaceId)}&includeShared=true`,
			signal
		);
		// Older server responses may encode an empty Go slice as null. Keep the
		// client contract stable so selectors can always render an array.
		return { ...result, items: result.items ?? [] };
	},
	create: (ownerWorkspaceId: string, name: string, signal?: AbortSignal) =>
		api.post<Tag>('/api/v1/tags', { ownerWorkspaceId, name }, signal),
	update: (id: string, name: string, version: number, signal?: AbortSignal) =>
		api.put<Tag>(`/api/v1/tags/${id}`, { name, version }, signal),
	setWorkspaceAccess: (
		id: string,
		workspaceIds: string[],
		version: number,
		removalDecisions: { workspaceId: string; removeAssignments: boolean }[],
		signal?: AbortSignal
	) =>
		api.put<Tag>(
			`/api/v1/tags/${id}/workspace-access`,
			{ workspaceIds, version, removalDecisions },
			signal
		)
};
