import { api } from './client';
export interface SearchResult {
	type: 'note' | 'question' | 'answer';
	id: string;
	workspaceId: string;
	workspaceName: string;
	title?: string;
	snippet: string;
	destination: string;
}
export interface SearchPage {
	items: SearchResult[];
	nextCursor?: string | null;
}
export function searchApi(
	query: {
		q: string;
		contentScope: 'notes' | 'questions' | 'answers' | 'everything';
		workspaceScope: 'current' | 'all';
		workspaceId?: string;
	},
	signal?: AbortSignal
): Promise<SearchPage> {
	const params = new URLSearchParams({
		q: query.q,
		contentScope: query.contentScope,
		workspaceScope: query.workspaceScope
	});
	if (query.workspaceId) params.set('workspaceId', query.workspaceId);
	return api.get<SearchPage>(`/api/v1/search?${params}`, signal);
}
