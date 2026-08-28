import { writable } from 'svelte/store';
import type { QuestionStatus } from '$lib/types/question';
export interface QuestionQueryState {
	status: QuestionStatus[];
	topicId: string;
	tagId: string;
	hasAnswer: '' | 'true' | 'false';
	hasLinkedNotes: '' | 'true' | 'false';
	sort: string;
	direction: 'asc' | 'desc';
}
export const questionQuery = writable<QuestionQueryState>({
	status: [],
	topicId: '',
	tagId: '',
	hasAnswer: '',
	hasLinkedNotes: '',
	sort: 'updatedAt',
	direction: 'desc'
});
export function resetQuestionQuery(): void {
	questionQuery.set({
		status: [],
		topicId: '',
		tagId: '',
		hasAnswer: '',
		hasLinkedNotes: '',
		sort: 'updatedAt',
		direction: 'desc'
	});
}
