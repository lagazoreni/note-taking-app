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
			if (index < 0) continue;

			const range = document.createRange();
			range.setStart(node, index);
			range.setEnd(node, index + value.length);
			const selection = window.getSelection();
			selection?.removeAllRanges();
			selection?.addRange(range);
			element.dispatchEvent(new MouseEvent('mouseup', { bubbles: true }));
			return;
		}
		throw new Error(`Could not select: ${value}`);
	}, text);
}

async function clickOutsideCaptureLayer(page: Page) {
	// Dispatch on body itself so the target is outside the reader and every
	// floating capture layer without navigating away from the note.
	await page.locator('body').dispatchEvent('mousedown');
	await page.locator('body').dispatchEvent('mouseup');
	await page.locator('body').dispatchEvent('click');
}

async function runCaptureFocusFlow(page: Page, suffix: string) {
	await externalNetworkBlock(page);
	const questionPassage = 'This passage is used to test question focus.';
	const annotationPassage = 'This passage is used to test annotation focus.';
	const questionText = `Why does this passage matter? ${suffix}`;
	const annotationText = `Remember this annotation. ${suffix}`;

	await createWorkspace(page, `Capture focus ${suffix}`);
	await createNote(
		page,
		`Capture focus note ${suffix}`,
		`${questionPassage}\n\n${annotationPassage}`
	);

	// Escape must dismiss a focused question composer without creating a mark.
	await selectReaderText(page, questionPassage);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	const discardedQuestion = page.getByLabel('Question text');
	await discardedQuestion.fill(`discarded question ${suffix}`);
	await expect(discardedQuestion).toBeFocused();
	await expect(page.locator('.composer')).toBeVisible();
	await page.keyboard.press('Escape');
	await expect(page.locator('.composer')).toHaveCount(0);
	await expect(page.locator('article mark[data-annotation-id]')).toHaveCount(0);

	// The saved question must keep its textbox focused and use the entered value.
	await selectReaderText(page, questionPassage);
	await page.getByRole('button', { name: 'Ask a question' }).click();
	const questionInput = page.getByLabel('Question text');
	await questionInput.fill(questionText);
	await expect(questionInput).toBeFocused();
	await expect(questionInput).toHaveValue(questionText);
	await expect(page.locator('.composer')).toBeVisible();
	await page.getByRole('button', { name: 'Save question' }).click();

	const questionDialog = page.getByRole('dialog', { name: /question/i });
	await expect(questionDialog).toBeVisible();
	await expect(questionDialog.getByText(questionText, { exact: true })).toBeVisible();
	await questionDialog.getByRole('button', { name: 'Close' }).click();

	// Cancel must dismiss an annotation composer without creating a mark.
	await selectReaderText(page, annotationPassage);
	await page.getByRole('button', { name: 'Add annotation' }).click();
	const discardedAnnotation = page.getByLabel('Annotation');
	await discardedAnnotation.fill(`discarded annotation ${suffix}`);
	await expect(discardedAnnotation).toBeFocused();
	await expect(page.locator('.composer')).toBeVisible();
	await page.getByRole('button', { name: 'Cancel' }).click();
	await expect(page.locator('.composer')).toHaveCount(0);
	await expect(page.locator('article mark[data-annotation-id]')).toHaveCount(1);

	// An intentional outside click must also dismiss the focused annotation
	// composer, after which the same passage can be captured successfully.
	await selectReaderText(page, annotationPassage);
	await page.getByRole('button', { name: 'Add annotation' }).click();
	const annotationInput = page.getByLabel('Annotation');
	await annotationInput.fill(annotationText);
	await expect(annotationInput).toBeFocused();
	await expect(annotationInput).toHaveValue(annotationText);
	await expect(page.locator('.composer')).toBeVisible();
	await clickOutsideCaptureLayer(page);
	await expect(page.locator('.composer')).toHaveCount(0);

	await selectReaderText(page, annotationPassage);
	await page.getByRole('button', { name: 'Add annotation' }).click();
	const savedAnnotationInput = page.getByLabel('Annotation');
	await savedAnnotationInput.fill(annotationText);
	await expect(savedAnnotationInput).toBeFocused();
	await expect(page.locator('.composer')).toBeVisible();
	await page.getByRole('button', { name: 'Save annotation' }).click();

	const annotationDialog = page.getByRole('dialog', { name: /annotation/i });
	await expect(annotationDialog).toBeVisible();
	await expect(annotationDialog.getByText(annotationText, { exact: true })).toBeVisible();
}

test.describe('capture composer focus', () => {
	test('keeps question and annotation capture focused on desktop', async ({ page }, info) => {
		await runCaptureFocusFlow(page, `desktop-${info.project.name}-${Date.now()}`);
	});

	test.describe('mobile viewport', () => {
		test.use({ viewport: { width: 375, height: 667 } });

		test('keeps question and annotation capture focused on mobile', async ({ page }, info) => {
			await runCaptureFocusFlow(page, `mobile-${info.project.name}-${Date.now()}`);
		});
	});
});
