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

test('lists notes and captures a highlight question and annotation', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const workspaceName = `Highlight ${suffix}`;
	const noteA = `Mitochondria ${suffix}`;
	const noteB = `ATP ${suffix}`;

	await createWorkspace(page, workspaceName);
	await createNote(
		page,
		noteA,
		'The mitochondria is the powerhouse of the cell. ATP stores energy.'
	);
	await createNote(page, noteB, 'A second note for the list.');

	await page.goto('/notes');
	await expect(page.getByRole('link', { name: noteA })).toBeVisible();
	await expect(page.getByRole('link', { name: noteB })).toBeVisible();

	await page.getByRole('link', { name: noteA }).click();
	await selectReaderText(page, 'The mitochondria is the powerhouse of the cell.');
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill('What does mitochondria do?');
	await page.getByRole('button', { name: 'Save question' }).click();
	const questionHighlight = page
		.locator('article mark[data-annotation-id]')
		.filter({ hasText: 'The mitochondria is the powerhouse of the cell.' })
		.first();
	await expect(questionHighlight).toBeVisible();
	await expect(page.locator('body')).not.toContainText('{{question:');

	await questionHighlight.click();
	const questionDialog = page.getByRole('dialog', { name: /question/i });
	await expect(questionDialog).toBeVisible();
	await expect(questionDialog.getByText('What does mitochondria do?', { exact: true })).toBeVisible();
	await page.getByRole('button', { name: 'Close' }).click();

	await selectReaderText(page, 'ATP stores energy.');
	await page.getByRole('button', { name: 'Add annotation' }).click();
	await page.getByLabel('Annotation').fill('Keep this for later.');
	await page.getByRole('button', { name: 'Save annotation' }).click();
	const annotationHighlight = page
		.locator('article mark[data-annotation-id]')
		.filter({ hasText: 'ATP stores energy.' })
		.first();
	await expect(annotationHighlight).toBeVisible();

	await page.getByRole('link', { name: 'Active Questions' }).click();
	await expect(page.getByText('What does mitochondria do?')).toBeVisible();
	await expect(page.getByText('Keep this for later.')).toHaveCount(0);
});
