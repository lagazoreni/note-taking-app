import { api } from './client';
import type { Note, NotePage, NoteWrite } from '$lib/types/note';
export interface DeletionPreview {
	previewToken: string;
	expiresAt: string;
	noteId: string;
	noteVersion: number;
	singlyLinkedQuestions: import('$lib/types/question').QuestionSummary[];
	multiplyLinkedQuestions: import('$lib/types/question').QuestionSummary[];
	childNotes: { id: string; title: string; version: number }[];
}
export const notesApi = {
	list: (workspaceId: string, signal?: AbortSignal) =>
		api.get<NotePage>(`/api/v1/notes?workspaceId=${encodeURIComponent(workspaceId)}`, signal),
	create: (write: NoteWrite, signal?: AbortSignal) =>
		api.post<Note>('/api/v1/notes', write, signal),
	get: (id: string, signal?: AbortSignal) => api.get<Note>(`/api/v1/notes/${id}`, signal),
	update: (id: string, write: NoteWrite & { version: number }, signal?: AbortSignal) =>
		api.put<Note>(`/api/v1/notes/${id}`, write, signal),
	previewDeletion: (id: string, version: number, signal?: AbortSignal) =>
		api.post<DeletionPreview>(`/api/v1/notes/${id}/deletion-preview`, { version }, signal),
	deleteReviewed: (id: string, command: unknown, signal?: AbortSignal) =>
		api.post<void>(`/api/v1/notes/${id}/delete`, command, signal)
};
