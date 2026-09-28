import { writable } from 'svelte/store';

export type ToastKind = 'info' | 'success' | 'error';
export interface Toast {
	id: number;
	kind: ToastKind;
	message: string;
}

export const toasts = writable<Toast[]>([]);
let nextId = 1;
export function pushToast(message: string, kind: ToastKind = 'info', duration = 4500): number {
	const id = nextId++;
	toasts.update((items) => [...items, { id, kind, message }]);
	if (duration > 0) setTimeout(() => dismissToast(id), duration);
	return id;
}
export function dismissToast(id: number): void {
	toasts.update((items) => items.filter((item) => item.id !== id));
}
