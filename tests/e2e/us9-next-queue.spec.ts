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

async function captureQuestion(page: Page, passage: string, questionText: string) {
	await selectReaderText(page, passage);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	await page.getByLabel('Question text').fill(questionText);
	await page.getByRole('button', { name: 'Save question' }).click();
	const dialog = page.getByRole('dialog', { name: /question/i });
	await expect(dialog).toBeVisible();
	await dialog.getByRole('button', { name: 'Close' }).click();
}

async function setDueDate(page: Page, questionText: string, dueDate: string) {
	await page.goto('/questions');
	await page.getByRole('link', { name: questionText, exact: true }).click();
	await expect(page.getByLabel('Due date')).toBeVisible();
	await page.getByLabel('Due date').fill(dueDate);

	const update = page.waitForResponse((response) => {
		return (
			response.request().method() === 'PUT' &&
			response.url().includes('/api/v1/questions/') &&
			response.ok()
		);
	});
	await page.getByRole('button', { name: 'Save date' }).click();
	await update;
}

async function answerQuestion(page: Page, questionText: string, answer: string) {
	await page.goto('/questions');
	await page.getByRole('link', { name: questionText, exact: true }).click();
	await expect(page.getByLabel('Answer')).toBeVisible();
	await page.getByLabel('Answer').fill(answer);
	await page.getByLabel('Status').selectOption('answered');

	const update = page.waitForResponse((response) => {
		return (
			response.request().method() === 'PUT' &&
			response.url().includes('/api/v1/questions/') &&
			response.ok()
		);
	});
	await page.getByRole('button', { name: 'Save question' }).click();
	await update;
}

test('shows overdue questions in Next and omits answered questions', async ({ page }, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const workspaceName = `Next queue ${suffix}`;
	const noteTitle = `Next note ${suffix}`;
	const overduePassage = `Overdue passage ${suffix}.`;
	const answeredPassage = `Answered passage ${suffix}.`;
	const overdueQuestion = `Overdue question ${suffix}`;
	const answeredQuestion = `Answered question ${suffix}`;
	const yesterday = new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString().slice(0, 10);

	await createWorkspace(page, workspaceName);
	await createNote(page, noteTitle, `${overduePassage} ${answeredPassage}`);
	await captureQuestion(page, overduePassage, overdueQuestion);
	await captureQuestion(page, answeredPassage, answeredQuestion);

	await setDueDate(page, overdueQuestion, yesterday);
	await answerQuestion(page, answeredQuestion, `Answer ${suffix}`);

	await page.goto('/next');
	const overdueHeading = page.getByRole('heading', { name: 'Overdue', exact: true });
	await expect(overdueHeading).toBeVisible();
	await expect(overdueHeading.locator('..').getByRole('link', { name: overdueQuestion, exact: true })).toBeVisible();
	await expect(page.getByRole('link', { name: answeredQuestion, exact: true })).toHaveCount(0);

	await overdueHeading.locator('..').getByRole('link', { name: overdueQuestion, exact: true }).click();
	await expect(page).toHaveURL(/\/questions\/[0-9a-f-]{36}(?:\?[^#]*)?$/i);
});
