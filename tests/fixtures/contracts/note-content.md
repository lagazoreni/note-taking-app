# Note Content Contract

**Version**: 1  
**Media type**: `text/markdown; charset=utf-8`

## Purpose

Notes are portable Markdown text with one application directive for placing a shared question in context. The visible question and answer remain canonical question records; they are never copied into note text.

## MVP Markdown subset

The editor must create and correctly display:

- paragraphs separated by a blank line;
- ordered lists using `1. ` syntax;
- unordered lists using `- ` syntax.

Other Markdown syntax may be retained as text but is not an MVP editing or display guarantee. Raw HTML is disabled and must never be rendered as trusted markup.

## Question directive

A question may be attached to a passage with a wrapped directive:

```text
{{question:550e8400-e29b-41d4-a716-446655440000}}selected passage{{/question}}
```

The bare opening directive remains valid for legacy or unanchored questions:

```text
{{question:550e8400-e29b-41d4-a716-446655440000}}
```

Grammar:

```abnf
question-directive = opening-directive [ passage closing-directive ]
opening-directive  = "{{question:" uuid "}}"
closing-directive  = "{{/question}}"
passage            = *OCTET
uuid               = 8HEXDIG "-" 4HEXDIG "-" 4HEXDIG "-" 4HEXDIG "-" 12HEXDIG
```

`{{/question}}` closes a wrapped highlight; it is not a second directive. IDs are serialized in lowercase. Parsers may accept uppercase hexadecimal but canonical save output normalizes it to lowercase.

## Integrity rules

1. Each directive ID must resolve to an existing question.
2. The question and note must belong to the same workspace.
3. Each directive must have one matching `NoteQuestion` relationship.
4. Each `NoteQuestion` relationship must have one matching directive.
5. A question may occur at most once in a given note.
6. Directive order determines `NoteQuestion.position`, beginning at zero.
7. `displayMode` is relationship metadata (`expanded`, `collapsed`, or `link`) and is not encoded in Markdown.
8. A malformed directive causes save validation to fail; it is not silently converted to text.
9. Relationship and Markdown changes are committed atomically.

## Rendering rules

- Text Markdown is parsed with raw HTML disabled and sanitized before browser insertion.
- A wrapped directive renders its passage as a clickable `<mark data-annotation-id="<uuid>">`; the directive tokens are never visible to the reader.
- A bare directive renders a compact marker chip associated with the canonical question rather than visible raw token text.
- Clicking a mark opens a card using the canonical question or annotation record.
- Legacy `displayMode` relationship values may remain stored for compatibility, but reading view uses highlight/card presentation.
- Exported Markdown retains directives so links can be restored from the export’s relationship records.

## Example

```markdown
# Optional heading retained as text

What causes this behavior?

{{question:550e8400-e29b-41d4-a716-446655440000}}What causes this behavior?{{/question}}

Possible areas to review:

1. Input assumptions
2. Boundary conditions
3. Existing evidence

- Follow the linked source
- Update the answer later
```

The corresponding relationship contains:

```json
{
  "noteId": "7f877d74-6a53-442a-9a16-c1c85f029fe8",
  "questionId": "550e8400-e29b-41d4-a716-446655440000",
  "displayMode": "collapsed",
  "position": 0
}
```
