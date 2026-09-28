export interface FieldError {
	field: string;
	message: string;
}

export interface ErrorPayload {
	code: string;
	message: string;
	requestId: string;
	fieldErrors?: FieldError[];
	details?: Record<string, unknown>;
}

export class ApiError extends Error {
	readonly status: number;
	readonly code: string;
	readonly requestId: string;
	readonly fieldErrors: FieldError[];
	readonly details: Record<string, unknown> | undefined;

	constructor(status: number, payload: ErrorPayload) {
		super(payload.message);
		this.name = 'ApiError';
		this.status = status;
		this.code = payload.code;
		this.requestId = payload.requestId;
		this.fieldErrors = payload.fieldErrors ?? [];
		this.details = payload.details;
	}

	get isConflict(): boolean {
		return this.status === 409 || this.code === 'VERSION_CONFLICT';
	}
}

export function isApiError(error: unknown): error is ApiError {
	return error instanceof ApiError;
}
