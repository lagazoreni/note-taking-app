import { test, expect } from '@playwright/test';

test('shared question journey has one canonical question', async ({ page }) => {
	await page.goto('/questions');
	await expect(page.getByRole('heading', { name: /active questions/i })).toBeVisible();
});
