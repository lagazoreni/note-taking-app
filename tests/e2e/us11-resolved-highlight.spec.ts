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

test('answered highlights look resolved and can insert their answer into the note', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const workspaceName = `Resolved highlight ${suffix}`;
	const noteTitle = `Resolved note ${suffix}`;
	const passage = 'The mitochondria is the powerhouse of the cell.';
	const questionText = `What does mitochondria do? ${suffix}`;
	const answer = 'It produces energy for the cell.';

	await createWorkspace(page, workspaceName);
	await createNote(page, noteTitle, passage);

	await selectReaderText(page, passage);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill(questionText);
	await page.getByRole('button', { name: 'Save question' }).click();
	await expect(page.getByRole('dialog', { name: /question/i })).toBeVisible();

	await page.getByLabel('Answer').fill(answer);
	await page.getByLabel('Status').selectOption('answered');
	const questionUpdate = page.waitForResponse(
		(response) =>
			response.request().method() === 'PUT' && response.url().includes('/api/v1/questions/')
	);
	await page.getByRole('button', { name: 'Save question' }).click();
	await questionUpdate;
	await page.getByRole('button', { name: 'Close' }).click();

	// Answering alone must not mutate the note body.
	await expect(page.locator('article blockquote')).toHaveCount(0);

	await page.reload();
	await expect(page.locator('article blockquote')).toHaveCount(0);
	const mark = page.locator('article mark[data-annotation-id]').first();
	await expect(mark).toHaveAttribute('data-status', 'answered');
	await mark.click();
	await expect(page.getByRole('dialog', { name: /question/i })).toBeVisible();

	const insertButton = page.getByRole('button', { name: 'Insert answer into note' });
	await expect(insertButton).toBeVisible();
	const noteUpdate = page.waitForResponse(
		(response) =>
			response.request().method() === 'PUT' && response.url().includes('/api/v1/notes/')
	);
	await insertButton.click();
	await noteUpdate;
	await page.getByRole('button', { name: 'Close' }).click();

	await page.reload();
	await expect(page.locator('article blockquote')).toContainText(answer);
});
