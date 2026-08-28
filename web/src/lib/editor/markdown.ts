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

export function renderNoteHtml(markdown: string): string {
	rejectRawHtml(markdown);
	const tokens = tokenizeDirectives(markdown);
	const marks: { id: string }[] = [];
	let rebuilt = '';
	for (const token of tokens) {
		if (token.kind === 'question' && token.directive) {
			const index = marks.length;
			marks.push({ id: token.directive.id });
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
			const id = marks[Number(index)]?.id ?? '';
			return `<mark data-annotation-id="${id}">${inner}</mark>`;
		})
		.replace(/notedchip(\d+)Z/g, (_match, index: string) => {
			const id = marks[Number(index)]?.id ?? '';
			return `<mark data-annotation-id="${id}" class="annotation-chip"></mark>`;
		});
	if (typeof window === 'undefined') return withMarks;
	return DOMPurify.sanitize(withMarks, {
		USE_PROFILES: { html: true },
		ADD_ATTR: ['data-annotation-id']
	});
}
