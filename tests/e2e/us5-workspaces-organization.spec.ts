import { test, expect } from '@playwright/test';

test('workspace organization screen is available', async ({ page }) => { await page.goto('/workspaces'); await expect(page.getByRole('heading', { name: /workspaces/i })).toBeVisible(); });
