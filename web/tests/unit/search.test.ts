import { beforeEach, describe, expect, it, vi } from 'vitest';
import { searchApi } from '$lib/api/search';

describe('searchApi', () => {
	beforeEach(() => {
		vi.restoreAllMocks();
	});

	it('serializes the optional tag filter', async () => {
		const fetchMock = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
			new Response(JSON.stringify({ items: [], nextCursor: null }), {
				status: 200,
				headers: { 'content-type': 'application/json' }
			})
		);
		await searchApi({
			q: 'lighthouse',
			contentScope: 'notes',
			workspaceScope: 'current',
			workspaceId: '11111111-1111-4111-8111-111111111111',
			tagId: '22222222-2222-4222-8222-222222222222'
		});
		const requestURL = String(fetchMock.mock.calls[0][0]);
		expect(requestURL).toContain('tagId=22222222-2222-4222-8222-222222222222');
	});
});
