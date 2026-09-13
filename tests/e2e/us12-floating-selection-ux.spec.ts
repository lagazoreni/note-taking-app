import { test, expect, type Locator, type Page } from '@playwright/test';

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

function longNoteBody(passage: string): string {
	return [
		passage,
		...Array.from(
			{ length: 36 },
			(_, index) =>
				`Filler paragraph ${index + 1}: this deliberately long reading note keeps the selection controls away from the end of the document while the passage remains near the top.`
		)
	].join('\n\n');
}

type Box = { x: number; y: number; width: number; height: number };

async function assertInViewport(page: Page, locator: Locator, name: string): Promise<Box> {
	const box = await locator.boundingBox();
	if (!box) throw new Error(`${name} has no rendered bounding box`);
	const viewport = page.viewportSize();
	if (!viewport) throw new Error('The Playwright page has no viewport');

	expect(box.x, `${name} is clipped on the left`).toBeGreaterThanOrEqual(-1);
	expect(box.y, `${name} is clipped above the viewport`).toBeGreaterThanOrEqual(-1);
	expect(box.x + box.width, `${name} is clipped on the right`).toBeLessThanOrEqual(viewport.width + 1);
	expect(box.y + box.height, `${name} is clipped below the viewport`).toBeLessThanOrEqual(
		viewport.height + 1
	);
	return box;
}

async function assertTouchTargets(container: Locator, name: string) {
	const buttons = await container.getByRole('button').all();
	for (const [index, button] of buttons.entries()) {
		const box = await button.boundingBox();
		if (!box) throw new Error(`${name} button ${index + 1} has no rendered bounding box`);
		expect(box.width, `${name} button ${index + 1} is too narrow for touch`).toBeGreaterThanOrEqual(44);
		expect(box.height, `${name} button ${index + 1} is too short for touch`).toBeGreaterThanOrEqual(44);
	}
}

async function selectReaderText(page: Page, text: string): Promise<Box> {
	const reader = page.getByRole('article', { name: 'Reading note' });
	await expect(reader).toBeVisible();
	return reader.evaluate((element, value) => {
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
			const rect = range.getBoundingClientRect();
			element.dispatchEvent(new MouseEvent('mouseup', { bubbles: true }));
			return { x: rect.x, y: rect.y, width: rect.width, height: rect.height };
		}
		throw new Error(`Could not select: ${value}`);
	}, text);
}

async function runFloatingSelectionFlow(page: Page, suffix: string) {
	await externalNetworkBlock(page);
	const passage = 'This passage near the top should open the floating actions.';
	const questionText = `Why is this passage important? ${suffix}`;

	await createWorkspace(page, `Floating selection ${suffix}`);
	await createNote(page, `Long floating note ${suffix}`, longNoteBody(passage));
	await page.evaluate(() => window.scrollTo(0, 0));

	const selectionBox = await selectReaderText(page, passage);
	const viewport = page.viewportSize();
	if (viewport) {
		expect(selectionBox.y).toBeGreaterThanOrEqual(0);
		expect(selectionBox.y + selectionBox.height).toBeLessThanOrEqual(viewport.height);
	}
	const toolbar = page.getByRole('toolbar', { name: 'Selection actions' });
	await expect(toolbar).toBeVisible();
	const toolbarBox = await assertInViewport(page, toolbar, 'selection toolbar');
	if (viewport && viewport.width <= 640) await assertTouchTargets(toolbar, 'selection toolbar');
	await expect.poll(() => page.evaluate(() => window.scrollY)).toBe(0);

	// The toolbar must stay close to the selected passage instead of being appended
	// after the full long article.
	expect(Math.abs(toolbarBox.y - selectionBox.y)).toBeLessThan(320);

	await toolbar.getByRole('button', { name: 'Ask a question' }).click();
	const composer = page.locator('.composer');
	await expect(composer).toBeVisible();
	const composerBox = await assertInViewport(page, composer, 'question composer');
	if (viewport && viewport.width <= 640) await assertTouchTargets(composer, 'question composer');
	expect(Math.abs(composerBox.y - selectionBox.y)).toBeLessThan(420);

	await composer.getByLabel('Question text').fill(questionText);
	await composer.getByRole('button', { name: 'Save question' }).click();

	const highlight = page
		.locator('article mark[data-annotation-id]')
		.filter({ hasText: passage })
		.first();
	await expect(highlight).toBeVisible();

	// Saving opens the new card. Close it before exercising the click-and-scroll
	// path independently.
	const initialCard = page.getByRole('dialog', { name: /question/i });
	await expect(initialCard).toBeVisible();
	await initialCard.getByRole('button', { name: 'Close' }).click();

	await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight));
	await expect.poll(() => page.evaluate(() => window.scrollY)).toBeGreaterThan(0);
	await highlight.evaluate((element) => element.scrollIntoView({ block: 'start', inline: 'nearest' }));
	const highlightBox = await assertInViewport(page, highlight, 'highlight');

	await highlight.click();
	const card = page.getByRole('dialog', { name: /question/i });
	await expect(card).toBeVisible();
	const cardBox = await assertInViewport(page, card, 'question card');

	// Desktop cards should remain anchored to the highlighted passage rather than
	// falling back to a fixed bottom-corner dialog.
	if (viewport && viewport.width > 640) {
		expect(Math.abs(cardBox.y - highlightBox.y)).toBeLessThan(260);
	}
}

test.describe('floating selection UX', () => {
	test('keeps selection actions and the question card in context on desktop', async ({ page }, info) => {
		await runFloatingSelectionFlow(page, `desktop-${info.project.name}-${Date.now()}`);
	});

	test.describe('mobile viewport', () => {
		test.use({ viewport: { width: 375, height: 667 } });

		test('keeps selection actions and the question card accessible without clipping', async ({ page }, info) => {
			await runFloatingSelectionFlow(page, `mobile-${info.project.name}-${Date.now()}`);
		});
	});
});
