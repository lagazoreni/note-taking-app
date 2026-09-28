import { test, expect } from '@playwright/test';

test('external internet is blocked while localhost remains available', async ({ page }) => { await page.route('**/*', async (route) => { const url = new URL(route.request().url()); if (url.hostname !== '127.0.0.1' && url.hostname !== 'localhost') return route.abort(); await route.continue(); }); await page.goto('/'); await expect(page.getByRole('heading', { name: /keep questions in context/i })).toBeVisible(); });
