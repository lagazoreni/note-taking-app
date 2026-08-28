import { api } from './client';
import type { Workspace } from '$lib/types/workspace';
export const workspacesApi = {
	list: (signal?: AbortSignal) => api.get<{ items: Workspace[] }>('/api/v1/workspaces', signal),
	create: (name: string, signal?: AbortSignal) =>
		api.post<Workspace>('/api/v1/workspaces', { name }, signal),
	get: (id: string, signal?: AbortSignal) => api.get<Workspace>(`/api/v1/workspaces/${id}`, signal),
	update: (id: string, name: string, version: number, signal?: AbortSignal) =>
		api.put<Workspace>(`/api/v1/workspaces/${id}`, { name, version }, signal)
};
