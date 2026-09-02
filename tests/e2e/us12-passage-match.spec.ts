import { test, expect, type Page } from '@playwright/test';

const externalNetworkBlock = async (page: Page) => {
	await page.route('**/*', async (route) => {
		const url = new URL(route.request().url());
		if (url.hostname !== '127.0.0.1' && url.hostname !== 'localhost') return route.abort();
		return route.continue();
	});
};

async function createWorkspace(page: Page, name: string) {
	await page.goto('/workspaces');
	await page.getByLabel('New workspace').fill(name);
	await page.getByRole('button', { name: 'Create' }).click();
	await expect(page.getByText(/workspace created/i)).toBeVisible();
}

async function createNote(page: Page, title: string, body: string) {
	await page.goto('/notes/new');
	await page.getByLabel('Title').fill(title);
	await page.getByRole('textbox', { name: 'Note', exact: true }).fill(body);
	await page.getByRole('button', { name: 'Save note' }).click();
	await expect(page).toHaveURL(/\/notes\/[0-9a-f-]{36}/i);
}

async function selectReaderText(page: Page, text: string) {
	const reader = page.getByRole('article', { name: 'Reading note' });
	await expect(reader).toBeVisible();
	await reader.evaluate((el, value) => {
		const walker = document.createTreeWalker(el, NodeFilter.SHOW_TEXT);
		let node: Node | null;
		while ((node = walker.nextNode())) {
			const index = node.textContent?.indexOf(value) ?? -1;
			if (index >= 0) {
				const range = document.createRange();
				range.setStart(node, index);
				range.setEnd(node, index + value.length);
				const selection = window.getSelection();
				selection?.removeAllRanges();
				selection?.addRange(range);
				el.dispatchEvent(new MouseEvent('mouseup', { bubbles: true }));
				return;
			}
		}
		throw new Error(`Could not select: ${value}`);
	}, text);
}

test('captures a selection rendered from bold Markdown', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;

	await createWorkspace(page, `Passage match ${suffix}`);
	await createNote(page, `Bold note ${suffix}`, 'This is **bold** text.');

	await selectReaderText(page, 'bold');
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill('Why is this bold?');
	await page.getByRole('button', { name: 'Save question' }).click();

	const highlight = page.locator('mark[data-annotation-id]');
	await expect(highlight).toHaveText('bold');
	await expect(highlight.locator('strong')).toHaveText('bold');
	await expect(page.getByRole('alert')).toHaveCount(0);
	await expect(page.locator('body')).not.toContainText('{{question:');
});
