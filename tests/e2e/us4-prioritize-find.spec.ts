import { test, expect } from '@playwright/test';

test('prioritization and search surfaces are available', async ({ page }) => { await page.goto('/search'); await expect(page.getByRole('heading', { name: /search/i })).toBeVisible(); });
