import { ApiError, type ErrorPayload } from './errors';

export interface RequestOptions extends RequestInit {
	signal?: AbortSignal;
}

async function decodeError(response: Response): Promise<ErrorPayload> {
	try {
		const body = (await response.json()) as { error?: ErrorPayload };
		if (body.error) return body.error;
	} catch {
		// Fall through to a safe generic message.
	}
	return {
		code: response.status >= 500 ? 'INTERNAL_ERROR' : 'VALIDATION_FAILED',
		message: response.statusText || 'Request failed',
		requestId: response.headers.get('x-request-id') ?? ''
	};
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
	const headers = new Headers(options.headers);
	if (options.body && !(options.body instanceof FormData) && !headers.has('content-type'))
		headers.set('content-type', 'application/json');
	let response: Response;
	try {
		response = await fetch(path, { ...options, headers });
	} catch (error) {
		throw error instanceof Error ? error : new Error('Unable to reach the local service');
	}
	if (!response.ok) throw new ApiError(response.status, await decodeError(response));
	if (response.status === 204) return undefined as T;
	return (await response.json()) as T;
}

export const api = {
	get: <T>(path: string, signal?: AbortSignal) => request<T>(path, { method: 'GET', signal }),
	post: <T>(path: string, body?: unknown, signal?: AbortSignal) =>
		request<T>(path, {
			method: 'POST',
			body: body === undefined ? undefined : JSON.stringify(body),
			signal
		}),
	put: <T>(path: string, body: unknown, signal?: AbortSignal) =>
		request<T>(path, { method: 'PUT', body: JSON.stringify(body), signal }),
	delete: <T>(path: string, signal?: AbortSignal) => request<T>(path, { method: 'DELETE', signal })
};
