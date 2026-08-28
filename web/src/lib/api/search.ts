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

export interface SearchQuery {
	q: string;
	contentScope: 'notes' | 'questions' | 'answers' | 'everything';
	workspaceScope: 'current' | 'all';
	workspaceId?: string;
	tagId?: string;
}

export function searchApi(query: SearchQuery, signal?: AbortSignal): Promise<SearchPage> {
	const params = new URLSearchParams({
		q: query.q,
		contentScope: query.contentScope,
		workspaceScope: query.workspaceScope
	});
	if (query.workspaceId) params.set('workspaceId', query.workspaceId);
	if (query.tagId) params.set('tagId', query.tagId);
	return api.get<SearchPage>(`/api/v1/search?${params}`, signal);
}
