import { test, expect } from '@playwright/test';

test('organization navigation exposes accessible names', async ({ page }) => { await page.goto('/workspaces'); await expect(page.getByRole('combobox', { name: /current workspace/i })).toBeVisible(); await expect(page.getByRole('navigation', { name: /primary/i })).toBeVisible(); });
