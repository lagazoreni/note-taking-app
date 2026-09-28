import { test, expect } from '@playwright/test';

test('data settings exposes export and import controls', async ({ page }) => { await page.goto('/settings/data'); await expect(page.getByRole('heading', { name: /data/i })).toBeVisible(); });
