import { expect, test, type Page } from '@playwright/test';

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
	return page.url();
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

async function linkExistingQuestion(page: Page, passage: string, questionText: string) {
	await selectReaderText(page, passage);
	await page.getByRole('button', { name: 'Link existing question' }).click();
	const picker = page.getByRole('region', { name: 'Find an existing question' });
	await expect(picker).toBeVisible();
	await picker.getByText(questionText, { exact: true }).click();
}

test('links an existing question from reading to a second note', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const noteATitle = `Original note ${suffix}`;
	const noteBTitle = `Second note ${suffix}`;
	const passageA = `The original passage ${suffix}.`;
	const passageB = `The second passage ${suffix}.`;
	const questionText = `Why share this question? ${suffix}`;

	await createWorkspace(page, `Link existing ${suffix}`);
	const noteAUrl = await createNote(page, noteATitle, passageA);
	const noteBUrl = await createNote(page, noteBTitle, passageB);

	await page.goto(noteAUrl);
	await selectReaderText(page, passageA);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill(questionText);
	await page.getByRole('button', { name: 'Save question' }).click();
	await expect(page.locator('mark[data-annotation-id]')).toHaveText(passageA);

	await page.goto(noteBUrl);
	await linkExistingQuestion(page, passageB, questionText);
	await expect(page.locator('mark[data-annotation-id]')).toHaveText(passageB);

	await page.getByRole('link', { name: 'Active Questions' }).click();
	await page.getByRole('link', { name: questionText }).click();
	const sourceNote = page.getByLabel('Source note');
	await expect(sourceNote).toContainText(noteATitle);
	await expect(sourceNote).toContainText(noteBTitle);
});

test('creates an unlinked question and attaches it from reading', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const noteTitle = `Attach later ${suffix}`;
	const passage = `A passage waiting for a question ${suffix}.`;
	const questionText = `What should be attached later? ${suffix}`;

	await createWorkspace(page, `Attach question ${suffix}`);
	const noteUrl = await createNote(page, noteTitle, passage);

	await page.goto('/questions');
	await expect(page.getByRole('heading', { name: 'Active questions', exact: true })).toBeVisible();
	await page.getByRole('button', { name: 'New question' }).click();
	await page.getByRole('textbox', { name: /question text/i }).fill(questionText);
	await page.getByRole('button', { name: /create question|save question/i }).click();
	await expect(page.getByRole('link', { name: questionText })).toBeVisible();

	await page.goto(noteUrl);
	await linkExistingQuestion(page, passage, questionText);
	const highlight = page.locator('mark[data-annotation-id]').filter({ hasText: passage });
	await expect(highlight).toBeVisible();
	await highlight.click();
	await expect(page.getByRole('dialog', { name: /question/i })).toContainText(questionText);
});
