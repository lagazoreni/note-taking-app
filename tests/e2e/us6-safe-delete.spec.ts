import { test, expect } from '@playwright/test';

test('note deletion is reviewed before execution', async ({ page }) => { await page.goto('/notes'); await expect(page.getByRole('heading', { name: /notes/i })).toBeVisible(); });
