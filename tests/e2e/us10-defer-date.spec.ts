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
	await reader.evaluate((element, value) => {
		const walker = document.createTreeWalker(element, NodeFilter.SHOW_TEXT);
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
				element.dispatchEvent(new MouseEvent('mouseup', { bubbles: true }));
				return;
			}
		}
		throw new Error(`Could not select: ${value}`);
	}, text);
}

test('deferred questions require and save a resume date', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const workspaceName = `Deferred date ${suffix}`;
	const noteTitle = `Deferred note ${suffix}`;
	const passage = `Resume passage ${suffix}.`;
	const questionText = `Resume question ${suffix}`;
	const resumeDate = '2026-04-15';

	await createWorkspace(page, workspaceName);
	await createNote(page, noteTitle, passage);
	await selectReaderText(page, passage);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill(questionText);
	await page.getByRole('button', { name: 'Save question' }).click();

	await page.getByRole('link', { name: 'Active Questions' }).click();
	await expect(page.getByRole('link', { name: questionText, exact: true })).toBeVisible();
	await page.getByRole('link', { name: questionText, exact: true }).click();
	await expect(page).toHaveURL(/\/questions\/[0-9a-f-]{36}/i);

	await page.getByLabel('Status').selectOption('deferred');
	await page.getByRole('button', { name: 'Save question' }).click();
	await expect(page.getByText('A resume date is required to defer a question.')).toBeVisible();

	await page.getByLabel('Resume date').fill(resumeDate);
	const update = page.waitForResponse((response) => {
		return (
			response.request().method() === 'PUT' &&
			response.url().includes('/api/v1/questions/') &&
			response.ok()
		);
	});
	await page.getByRole('button', { name: 'Save question' }).click();
	await update;
	await expect(page.getByLabel('Resume date')).toHaveValue(resumeDate);
});
