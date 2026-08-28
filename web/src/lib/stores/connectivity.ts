import { browser } from '$app/environment';
import { writable } from 'svelte/store';

export const localServiceAvailable = writable(true);
export const internetBlocked = writable(false);

export function markServiceUnavailable(): void {
	localServiceAvailable.set(false);
}
export function markServiceAvailable(): void {
	localServiceAvailable.set(true);
}

if (browser) {
	window.addEventListener('online', () => internetBlocked.set(false));
	window.addEventListener('offline', () => internetBlocked.set(true));
}
