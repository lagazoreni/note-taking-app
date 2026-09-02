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
	await page.getByLabel('Note').fill(body);
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
	const dialog = page.getByRole('dialog');
	await expect(dialog).toBeVisible();
	await dialog.getByRole('button', { name: 'Close' }).click();
}

async function captureAnnotation(page: Page, passage: string, annotationText: string) {
	await selectReaderText(page, passage);
	await page.getByRole('button', { name: 'Add annotation' }).click();
	await page.getByLabel('Annotation').fill(annotationText);
	await page.getByRole('button', { name: 'Save annotation' }).click();
	const dialog = page.getByRole('dialog');
	await expect(dialog).toBeVisible();
	await dialog.getByRole('button', { name: 'Close' }).click();
}

async function answerOpenQuestion(page: Page, questionText: string, answer: string) {
	const dialog = page.getByRole('dialog', { name: /question/i });
	await expect(dialog).toBeVisible();
	await expect(dialog.getByText(questionText, { exact: true })).toBeVisible();
	await dialog.getByLabel('Answer').fill(answer);
	await dialog.getByLabel('Status').selectOption('answered');

	const update = page.waitForResponse((response) => {
		return (
			response.request().method() === 'PUT' &&
			response.url().includes('/api/v1/questions/') &&
			response.ok()
		);
	});
	await dialog.getByRole('button', { name: 'Save question' }).click();
	await update;
	await dialog.getByRole('button', { name: 'Close' }).click();
}

test('walks open questions in note order and excludes annotations and answered questions', async ({
	page
}, info) => {
	await externalNetworkBlock(page);
	const suffix = `${info.project.name}-${Date.now()}`;
	const workspaceName = `Note rail ${suffix}`;
	const noteTitle = `Rail note ${suffix}`;
	const firstPassage = `First rail passage ${suffix}.`;
	const annotationPassage = `Annotation rail passage ${suffix}.`;
	const secondPassage = `Second rail passage ${suffix}.`;
	const answeredPassage = `Answered rail passage ${suffix}.`;
	const firstQuestion = `First open question ${suffix}`;
	const secondQuestion = `Second open question ${suffix}`;
	const answeredQuestion = `Already answered question ${suffix}`;
	const annotationText = `Annotation only ${suffix}`;

	await createWorkspace(page, workspaceName);
	await createNote(
		page,
		noteTitle,
		[ firstPassage, annotationPassage, secondPassage, answeredPassage ].join(' ')
	);

	await captureQuestion(page, firstPassage, firstQuestion);
	await captureAnnotation(page, annotationPassage, annotationText);
	await captureQuestion(page, secondPassage, secondQuestion);
	await captureQuestion(page, answeredPassage, answeredQuestion);

	await page.getByText(answeredPassage, { exact: true }).click();
	await answerOpenQuestion(page, answeredQuestion, `Answer ${suffix}`);
	await page.reload();
	await expect(page.getByRole('article', { name: 'Reading note' })).toBeVisible();

	let rail = page.getByRole('complementary', { name: 'Open questions in this note' });
	await expect(rail).toBeVisible();
	await expect(rail).toContainText(firstQuestion);
	await expect(rail).toContainText(secondQuestion);
	await expect(rail).not.toContainText(annotationText);
	await expect(rail).not.toContainText(answeredQuestion);

	await rail.getByRole('button', { name: 'Next unanswered' }).click();
	await expect(page.getByRole('dialog', { name: /question/i })).toContainText(firstQuestion);
	await answerOpenQuestion(page, firstQuestion, `First answer ${suffix}`);
	await page.reload();
	await expect(page.getByRole('article', { name: 'Reading note' })).toBeVisible();

	rail = page.getByRole('complementary', { name: 'Open questions in this note' });
	await expect(rail).toContainText(secondQuestion);
	await expect(rail).not.toContainText(firstQuestion);
	await rail.getByRole('button', { name: 'Next unanswered' }).click();
	await expect(page.getByRole('dialog', { name: /question/i })).toContainText(secondQuestion);
	await answerOpenQuestion(page, secondQuestion, `Second answer ${suffix}`);
	await page.reload();
	await expect(page.getByRole('article', { name: 'Reading note' })).toBeVisible();

	rail = page.getByRole('complementary', { name: 'Open questions in this note' });
	await expect(rail).toContainText('No open questions in this note');
	const nextButton = rail.getByRole('button', { name: 'Next unanswered' });
	if (await nextButton.count()) await expect(nextButton).toBeDisabled();
	await expect(rail).not.toContainText(annotationText);
	await expect(rail).not.toContainText(answeredQuestion);
});
