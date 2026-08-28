import { test, expect } from '@playwright/test';

test.describe('responsive secondary routes', () => { for (const path of ['/workspaces', '/search', '/settings/data']) test(path, async ({ page }) => { await page.setViewportSize({ width: 390, height: 800 }); await page.goto(path); await expect(page.locator('main')).toBeVisible(); }); });
