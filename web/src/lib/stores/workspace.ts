import { browser } from '$app/environment';
import { get, derived, writable } from 'svelte/store';
import type { Workspace } from '$lib/types/workspace';

const STORAGE_KEY = 'noted.currentWorkspace';
const initialWorkspaceId = browser ? localStorage.getItem(STORAGE_KEY) : null;

export const workspaces = writable<Workspace[]>([]);
export const currentWorkspaceId = writable<string | null>(initialWorkspaceId);
export const currentWorkspace = writable<Workspace | null>(null);

/** True after the first workspace list request has completed (successfully or not). */
export const workspaceHydrated = writable(false);
export const workspaceLoading = writable(false);
export const workspaceLoadError = writable<string | null>(null);
export const workspaceReady = derived(
	[workspaceHydrated, currentWorkspace],
	([$hydrated, $workspace]) => $hydrated && $workspace !== null
);

let workspaceList: Workspace[] = [];

function persistSelection(id: string | null): void {
	if (!browser) return;
	if (id) localStorage.setItem(STORAGE_KEY, id);
	else localStorage.removeItem(STORAGE_KEY);
}

function syncCurrentWorkspace(): void {
	const selected = get(currentWorkspaceId);
	currentWorkspace.set(workspaceList.find((workspace) => workspace.id === selected) ?? null);
}

function reconcileSelection(): void {
	const selected = get(currentWorkspaceId);
	const match = selected ? workspaceList.find((workspace) => workspace.id === selected) : undefined;
	if (match) {
		currentWorkspace.set(match);
		return;
	}

	const fallback = workspaceList[0]?.id ?? null;
	if (selected !== fallback) currentWorkspaceId.set(fallback);
	else currentWorkspace.set(fallback ? workspaceList[0] : null);
}

currentWorkspaceId.subscribe((id) => {
	persistSelection(id);
	syncCurrentWorkspace();
	if (
		get(workspaceHydrated) &&
		id &&
		workspaceList.length > 0 &&
		!workspaceList.some((workspace) => workspace.id === id)
	) {
		reconcileSelection();
	}
});

workspaces.subscribe((list) => {
	workspaceList = list;
	syncCurrentWorkspace();
	if (get(workspaceHydrated) || list.length > 0) reconcileSelection();
});

/** Mark the beginning of a fresh server-backed workspace list hydration. */
export function beginWorkspaceHydration(): void {
	workspaceLoading.set(true);
	workspaceLoadError.set(null);
	workspaceHydrated.set(false);
}

/** Replace the local list with the authoritative list and repair stale selection. */
export function setWorkspaceList(items: Workspace[]): void {
	workspaceHydrated.set(false);
	workspaces.set([...items]);
	workspaceHydrated.set(true);
	workspaceLoading.set(false);
	workspaceLoadError.set(null);
	reconcileSelection();
}

/** Keep the freshly-created workspace in the local list and select it immediately. */
export function addWorkspace(workspace: Workspace): void {
	const next = workspaceList.some((item) => item.id === workspace.id)
		? workspaceList.map((item) => (item.id === workspace.id ? workspace : item))
		: [...workspaceList, workspace];
	workspaces.set(next);
	workspaceHydrated.set(true);
	workspaceLoading.set(false);
	workspaceLoadError.set(null);
	selectWorkspace(workspace.id);
}

/** Leave the shell usable after a failed hydration without trusting a stale ID. */
export function failWorkspaceHydration(message: string): void {
	workspaceHydrated.set(false);
	workspaces.set([]);
	currentWorkspaceId.set(null);
	workspaceLoadError.set(message);
	workspaceLoading.set(false);
	workspaceHydrated.set(true);
}

export function selectWorkspace(id: string | null): void {
	const candidate = id?.trim() || null;
	if (
		candidate &&
		workspaceList.length > 0 &&
		!workspaceList.some((workspace) => workspace.id === candidate)
	) {
		reconcileSelection();
		return;
	}
	currentWorkspaceId.set(candidate);
}

/** Reset browser-local selection; useful when leaving a workspace-less setup state. */
export function clearWorkspaceSelection(): void {
	currentWorkspaceId.set(null);
}
