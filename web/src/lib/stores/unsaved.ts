import { writable } from 'svelte/store';

export interface UnsavedState {
	key: string;
	dirty: boolean;
	message?: string;
}
export const unsaved = writable<UnsavedState | null>(null);
export function markUnsaved(key: string, message?: string): void {
	unsaved.set({ key, dirty: true, message });
}
export function markSaved(key: string): void {
	unsaved.update((current) => (current?.key === key ? null : current));
}
