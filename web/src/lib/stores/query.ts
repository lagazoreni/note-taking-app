import { writable } from 'svelte/store';

export const queryVersion = writable(0);
export const invalidatedKeys = writable<string[]>([]);

export function invalidate(...keys: string[]): void {
	queryVersion.update((version) => version + 1);
	invalidatedKeys.update((current) => [...new Set([...current, ...keys])]);
}
export function clearInvalidations(): void {
	invalidatedKeys.set([]);
}
