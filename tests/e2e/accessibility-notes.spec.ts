import { test, expect } from '@playwright/test';

test.describe('note accessibility', () => { test('note entry has labelled controls and keyboard focus', async ({ page }) => { await page.goto('/notes/new'); await expect(page.getByLabel('Title')).toBeVisible(); await expect(page.getByRole('textbox', { name: 'Note', exact: true })).toBeVisible(); await page.getByLabel('Title').focus(); await expect(page.getByLabel('Title')).toBeFocused(); }); });
