import { api } from './client';
import type {
	Question,
	QuestionPage,
	QuestionStatus,
	Priority,
	QuestionKind
} from '$lib/types/question';
export interface QuestionQuery {
	workspaceId: string;
	status?: QuestionStatus[];
	kind?: QuestionKind;
	cursor?: string;
	pageSize?: number;
	sort?: string;
	direction?: 'asc' | 'desc';
}
function params(query: QuestionQuery): string {
	const value = new URLSearchParams({ workspaceId: query.workspaceId });
	if (query.status?.length) value.set('status', query.status.join(','));
	if (query.kind) value.set('kind', query.kind);
	if (query.cursor) value.set('cursor', query.cursor);
	if (query.pageSize) value.set('pageSize', String(query.pageSize));
	if (query.sort) value.set('sort', query.sort);
	if (query.direction) value.set('direction', query.direction);
	return value.toString();
}
export const questionsApi = {
	list: (query: QuestionQuery, signal?: AbortSignal) =>
		api.get<QuestionPage>(`/api/v1/questions?${params(query)}`, signal),
	get: (id: string, signal?: AbortSignal) => api.get<Question>(`/api/v1/questions/${id}`, signal),
	create: (
		body: {
			workspaceId: string;
			questionText: string;
			status: QuestionStatus;
			priority: Priority;
			answerMarkdown?: string | null;
			tagIds: string[];
			kind?: QuestionKind;
		},
		signal?: AbortSignal
	) => api.post<Question>('/api/v1/questions', body, signal),
	update: (
		id: string,
		body: Partial<Question> & { workspaceId: string; version: number },
		signal?: AbortSignal
	) => api.put<Question>(`/api/v1/questions/${id}`, body, signal)
};
