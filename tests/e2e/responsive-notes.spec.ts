import { test, expect } from '@playwright/test';

test.describe('responsive notes', () => { for (const width of [375, 1280]) test(`renders note entry at ${width}px`, async ({ page }) => { await page.setViewportSize({ width, height: 800 }); await page.goto('/notes/new'); await expect(page.getByRole('heading', { name: /create a note/i })).toBeVisible(); }); });
