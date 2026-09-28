import { marked } from 'marked';
import DOMPurify from 'dompurify';
import { tokenizeDirectives } from './directives';

function escapeHtml(value: string): string {
	return value.replace(
		/[&<>"']/g,
		(character) =>
			({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[character] ??
			character
	);
}

const renderer = new marked.Renderer();
renderer.html = ({ text }: { text: string }) => escapeHtml(text);

export function rejectRawHtml(markdown: string): void {
	if (/<\/?[a-z][^>]*>/i.test(markdown)) throw new Error('Raw HTML is not allowed in notes');
}

export function renderMarkdown(markdown: string): string {
	rejectRawHtml(markdown);
	const html = marked.parse(markdown, {
		renderer,
		gfm: true,
		breaks: false,
		async: false
	}) as string;
	if (typeof window === 'undefined') return html;
	return DOMPurify.sanitize(html, { USE_PROFILES: { html: true } });
}

export function renderNoteHtml(
	markdown: string,
	questions: { id: string; status?: string; kind?: string }[] = []
): string {
	rejectRawHtml(markdown);
	const tokens = tokenizeDirectives(markdown);
	const questionById = new Map(questions.map((question) => [question.id.toLowerCase(), question]));
	const marks: { id: string; status: string; kind: string }[] = [];
	let rebuilt = '';
	for (const token of tokens) {
		if (token.kind === 'question' && token.directive) {
			const index = marks.length;
			const question = questionById.get(token.directive.id);
			marks.push({
				id: token.directive.id,
				status: question?.status || 'unanswered',
				kind: question?.kind || 'question'
			});
			if (token.directive.wrapped) {
				rebuilt += `notedwrap${index}Z${token.directive.snippet ?? ''}notedend${index}Z`;
			} else {
				rebuilt += `notedchip${index}Z`;
			}
		} else {
			rebuilt += token.value;
		}
	}
	const html = marked.parse(rebuilt, {
		renderer,
		gfm: true,
		breaks: false,
		async: false
	}) as string;
	const withMarks = html
		.replace(/notedwrap(\d+)Z([\s\S]*?)notedend\1Z/g, (_match, index: string, inner: string) => {
			const mark = marks[Number(index)];
			const id = mark?.id ?? '';
			const status = mark?.status ?? 'unanswered';
			const kind = mark?.kind ?? 'question';
			return `<mark data-annotation-id="${id}" data-status="${escapeHtml(status)}" data-kind="${escapeHtml(kind)}">${inner}</mark>`;
		})
		.replace(/notedchip(\d+)Z/g, (_match, index: string) => {
			const mark = marks[Number(index)];
			const id = mark?.id ?? '';
			const status = mark?.status ?? 'unanswered';
			const kind = mark?.kind ?? 'question';
			return `<mark data-annotation-id="${id}" data-status="${escapeHtml(status)}" data-kind="${escapeHtml(kind)}" class="annotation-chip"></mark>`;
		});
	if (typeof window === 'undefined') return withMarks;
	return DOMPurify.sanitize(withMarks, {
		USE_PROFILES: { html: true },
		ADD_ATTR: ['data-annotation-id', 'data-status', 'data-kind']
	});
}
