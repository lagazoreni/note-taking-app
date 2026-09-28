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

export function findSelectionInMarkdown(
	markdown: string,
	selectedText: string
): { index: number; length: number } {
	if (!selectedText.trim()) throw new Error('selection is empty');

	const markerCharacters = new Set(['*', '_', '`']);
	const rangeIncludingMarkers = (start: number, end: number) => {
		let rangeStart = start;
		let rangeEnd = end;
		if (
			markerCharacters.has(markdown[rangeStart - 1] ?? '') &&
			markerCharacters.has(markdown[rangeEnd] ?? '')
		) {
			while (rangeStart > 0 && markerCharacters.has(markdown[rangeStart - 1] ?? '')) rangeStart -= 1;
			while (rangeEnd < markdown.length && markerCharacters.has(markdown[rangeEnd] ?? '')) rangeEnd += 1;
		}
		return { index: rangeStart, length: rangeEnd - rangeStart };
	};

	const index = markdown.indexOf(selectedText);
	if (index >= 0) return rangeIncludingMarkers(index, index + selectedText.length);

	const collapseWhitespace = (
		value: string,
		skip?: (character: string) => boolean
	): { text: string; indexMap: number[]; endMap: number[] } => {
		const characters: string[] = [];
		const indexMap: number[] = [];
		const endMap: number[] = [];
		let sourceIndex = 0;

		while (sourceIndex < value.length) {
			const character = value[sourceIndex];
			if (skip?.(character)) {
				sourceIndex += 1;
				continue;
			}

			const whitespace = value.slice(sourceIndex).match(/^\s+/);
			if (whitespace) {
				characters.push(' ');
				indexMap.push(sourceIndex);
				endMap.push(sourceIndex + whitespace[0].length);
				sourceIndex += whitespace[0].length;
				continue;
			}

			characters.push(character);
			indexMap.push(sourceIndex);
			endMap.push(sourceIndex + 1);
			sourceIndex += 1;
		}

		return { text: characters.join(''), indexMap, endMap };
	};

	const collapsedSelection = selectedText.replace(/\s+/g, ' ');
	const collapsedMarkdown = collapseWhitespace(markdown);
	const collapsedIndex = collapsedMarkdown.text.indexOf(collapsedSelection);
	if (collapsedIndex >= 0) {
		const end = collapsedIndex + collapsedSelection.length;
		return rangeIncludingMarkers(
			collapsedMarkdown.indexMap[collapsedIndex],
			collapsedMarkdown.endMap[end - 1]
		);
	}

	const readableMarkdown = collapseWhitespace(markdown, (character) => markerCharacters.has(character));
	const readableIndex = readableMarkdown.text.indexOf(collapsedSelection);
	if (readableIndex >= 0) {
		const end = readableIndex + collapsedSelection.length;
		return rangeIncludingMarkers(
			readableMarkdown.indexMap[readableIndex],
			readableMarkdown.endMap[end - 1]
		);
	}

	throw new Error('Could not find that passage in the note. Try selecting plain text.');
}

export function wrapSelection(markdown: string, selectedText: string, id: string): string {
	const { index, length } = findSelectionInMarkdown(markdown, selectedText);
	const selectedSpan = markdown.slice(index, index + length);
	const before = markdown.slice(0, index);
	const after = markdown.slice(index + length);
	return `${before}{{question:${id.toLowerCase()}}}${selectedSpan}{{/question}}${after}`;
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
