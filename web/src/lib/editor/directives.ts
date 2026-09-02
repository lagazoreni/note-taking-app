import type { DisplayMode } from '$lib/types/question';

const uuid = '[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}';
const directive = new RegExp(`\\{\\{question:(${uuid})\\}\\}`, 'g');
const nextOpening = new RegExp(`\\{\\{question:${uuid}\\}\\}`);
const closeTag = '{{/question}}';

export interface QuestionDirective {
	id: string;
	displayMode: DisplayMode;
	position: number;
	snippet?: string;
	wrapped?: boolean;
}
export interface TextSegment {
	kind: 'text' | 'question';
	value: string;
	directive?: QuestionDirective;
}

export function tokenizeDirectives(markdown: string): TextSegment[] {
	const segments: TextSegment[] = [];
	let cursor = 0;
	let position = 0;
	for (const match of markdown.matchAll(directive)) {
		const index = match.index ?? 0;
		if (index < cursor) continue;
		if (index > cursor) segments.push({ kind: 'text', value: markdown.slice(cursor, index) });
		const after = markdown.slice(index + match[0].length);
		const closeIdx = after.indexOf(closeTag);
		const nextOpen = after.search(nextOpening);
		let snippet = '';
		let wrapped = false;
		let consumed = match[0].length;
		if (closeIdx >= 0 && (nextOpen === -1 || closeIdx < nextOpen)) {
			snippet = after.slice(0, closeIdx);
			wrapped = true;
			consumed = match[0].length + closeIdx + closeTag.length;
		}
		segments.push({
			kind: 'question',
			value: match[0],
			directive: {
				id: match[1].toLowerCase(),
				displayMode: 'collapsed',
				position: position++,
				snippet,
				wrapped
			}
		});
		cursor = index + consumed;
	}
	if (cursor < markdown.length) segments.push({ kind: 'text', value: markdown.slice(cursor) });
	return segments;
}

export function directiveIds(markdown: string): string[] {
	return tokenizeDirectives(markdown)
		.filter((item) => item.kind === 'question')
		.map((item) => item.directive!.id);
}
export function insertDirective(markdown: string, id: string): string {
	const suffix = markdown && !markdown.endsWith('\n') ? '\n\n' : '';
	return `${markdown}${suffix}{{question:${id.toLowerCase()}}}`;
}
export function wrapSelection(markdown: string, selectedText: string, id: string): string {
	if (!selectedText) throw new Error('selection is empty');
	const index = markdown.indexOf(selectedText);
	if (index < 0) throw new Error('selection not found in note');
	const before = markdown.slice(0, index);
	const after = markdown.slice(index + selectedText.length);
	return `${before}{{question:${id.toLowerCase()}}}${selectedText}{{/question}}${after}`;
}

function startsWithAnswerBlockquote(text: string, answer: string): boolean {
	const afterWhitespace = text.replace(/^\s+/, '');
	if (!afterWhitespace.startsWith('>')) return false;

	const quotedLines: string[] = [];
	for (const line of afterWhitespace.split(/\r?\n/)) {
		if (!/^\s*>/.test(line)) break;
		quotedLines.push(line.replace(/^\s*>\s?/, ''));
	}
	return quotedLines.join('\n').includes(answer);
}

export function insertAnswerAfterDirective(markdown: string, id: string, answer: string): string {
	const trimmedAnswer = answer.trim();
	if (!trimmedAnswer) throw new Error('answer is empty');

	const needle = id.toLowerCase();
	for (const match of markdown.matchAll(directive)) {
		if (match[1].toLowerCase() !== needle) continue;

		const openingEnd = (match.index ?? 0) + match[0].length;
		const afterOpening = markdown.slice(openingEnd);
		const closeIndex = afterOpening.indexOf(closeTag);
		const nextOpeningIndex = afterOpening.search(nextOpening);
		const wrapped = closeIndex >= 0 && (nextOpeningIndex === -1 || closeIndex < nextOpeningIndex);
		const directiveEnd = wrapped ? openingEnd + closeIndex + closeTag.length : openingEnd;
		const afterDirective = markdown.slice(directiveEnd);

		if (startsWithAnswerBlockquote(afterDirective, trimmedAnswer)) return markdown;

		const blockquote = trimmedAnswer
			.split(/\r?\n/)
			.map((line) => `> ${line}`)
			.join('\n');
		return `${markdown.slice(0, directiveEnd)}\n\n${blockquote}\n${afterDirective}`;
	}

	return markdown;
}

export function stripDirectives(markdown: string): string {
	return tokenizeDirectives(markdown)
		.map((token) => {
			if (token.kind === 'question') {
				return token.directive?.wrapped ? (token.directive.snippet ?? '') : '';
			}
			return token.value;
		})
		.join('');
}
export function removeDirective(markdown: string, id: string): string {
	const needle = id.toLowerCase();
	return tokenizeDirectives(markdown)
		.map((token) => {
			if (token.kind === 'question' && token.directive?.id === needle) {
				return token.directive.wrapped ? (token.directive.snippet ?? '') : '';
			}
			return token.value;
		})
		.join('')
		.replace(/\n{3,}/g, '\n\n')
		.trim();
}
