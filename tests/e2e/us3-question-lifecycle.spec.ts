import { test, expect } from '@playwright/test';

test('answered view is a distinct filtered view', async ({ page }) => { await page.goto('/answered'); await expect(page.getByRole('heading', { name: /answered questions/i })).toBeVisible(); });
