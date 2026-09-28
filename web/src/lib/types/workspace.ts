export interface Workspace {
	id: string;
	name: string;
	createdAt: string;
	updatedAt: string;
	version: number;
}

export interface Topic {
	id: string;
	workspaceId: string;
	name: string;
	createdAt: string;
	updatedAt: string;
	version: number;
}

export interface Tag {
	id: string;
	ownerWorkspaceId: string;
	name: string;
	availableWorkspaceIds: string[];
	createdAt: string;
	updatedAt: string;
	version: number;
}
