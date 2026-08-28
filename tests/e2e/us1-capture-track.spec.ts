import { test, expect } from '@playwright/test';

test('capture and track questions without external internet', async ({ page }) => {
	await page.route('**/*', async (route) => {
		const url = new URL(route.request().url());
		if (url.hostname !== '127.0.0.1' && url.hostname !== 'localhost') return route.abort();
		return route.continue();
	});
	await page.goto('/');
	await expect(page.getByRole('heading', { name: /keep questions in context/i })).toBeVisible();
});
